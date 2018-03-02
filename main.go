package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/network"
	"github.com/larse514/amazonian/output"
	"github.com/larse514/amazonian/service"
)

const (
	fileName = "amazonian-output"
	tenant   = "amazonian"
)

//Todo- refactor main to be testable
func main() {
	args, err := commandlineargs.GenerateArgs()
	//if a required parameter is not specified, log error and exit
	if err != nil {
		os.Exit(1)
	}
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	// or environment variables
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create CloudFormation and LoadBalancer client in region
	svc := cloudformation.New(sess)
	elb := elbv2.New(sess)

	//Initialize the dependencies
	cfExecutor := cf.CFExecutor{Client: svc}
	serv := service.EcsService{Executor: cfExecutor, LoadBalancer: cf.AWSElb{Client: elb}}
	ecs := cluster.Ecs{Resource: cf.Stack{Client: svc}, Executor: cfExecutor}
	stack := cf.Stack{Client: svc}

	if !args.VPCExists {
		vpc := network.CreateDefaultVPC(args.VPCName, tenant)
		vpc.Executor = cfExecutor
		err := vpc.CreateNetwork()
		if err != nil {
			println("error creating vpc ", err.Error())
			os.Exit(1)
		}
	}
	if !args.VPCExists && args.VPC == "" {
		fmt.Println("VPC doesn't exist and VPCId was not provided, looking up values by name")
		//let's grab the vpc to get out needed output values
		vpcStack, err := stack.GetStack(&args.VPCName)
		if err != nil {
			println("error creating vpc ", err.Error())
			os.Exit(1)
		}
		fmt.Println("retrieved stack ", vpcStack.GoString())
		//i'm sorry, need to really refactor this whole block
		args.VPC = cf.GetOutputValue(vpcStack, "VPC-"+tenant)
		args.WSSubnetIDs = cf.GetOutputValue(vpcStack, "WSSubnet1-"+tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+tenant)
		//todo-get VPC private subnets to work
		args.ClusterSubnetIDs = cf.GetOutputValue(vpcStack, "WSSubnet1-"+tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+tenant)
	}
	//check if the cluster exists, if not create it
	if !args.ClusterExists {
		//create cluster
		clusterStruct := cluster.EcsCluster{}
		clusterStruct.DomainName = args.HostedZoneName
		clusterStruct.KeyName = args.KeyName
		clusterStruct.VpcID = args.VPC
		clusterStruct.ClusterSubnetIds = args.ClusterSubnetIDs
		clusterStruct.WSSubnetIds = args.WSSubnetIDs
		clusterStruct.DesiredCapacity = args.ClusterSize
		clusterStruct.MaxSize = args.MaxSize
		clusterStruct.InstanceType = args.InstanceType
		//create cluster
		err = ecs.CreateCluster(args.ClusterName, clusterStruct)
		if err != nil {
			println("error creating cluster ", err.Error())
			os.Exit(1)
		}
	}

	//now get the cluster based on the stack name provided
	ecs, err = ecs.GetCluster(args.ClusterName)

	if err != nil {
		fmt.Printf("error retrieving stack %s", args.ClusterName)
		os.Exit(1)
	}

	//let's get the priority for the next service
	// priority, err := cf.LoadBalancer.GetHighestPriority()
	//create the service struct, this is the struct that defines everything we need to create a container service
	//(note that for the time being only ECS is supported)
	serviceStruct := service.EcsService{}
	serviceStruct.Vpc = args.VPC
	serviceStruct.Image = args.Image
	serviceStruct.ServiceName = args.ServiceName
	serviceStruct.ContainerName = args.ContainerName
	serviceStruct.HostedZoneName = args.HostedZoneName
	//attempt to create the service
	err = serv.CreateService(&ecs, serviceStruct, args.ServiceName)

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	serviceName := strings.ToLower(args.ServiceName)
	url := "https://" + serviceName + "." + args.HostedZoneName
	err = output.WriteOutputFile(output.Output{fileName, args.ServiceName, args.ClusterName, url, args.VPC, args.VPCName})
	if err != nil {
		fmt.Println("Error writing output file ", err.Error())
	}
	fmt.Printf("Successfully created Container Service: %s, with url: %s \n", args.ServiceName, url)
}
