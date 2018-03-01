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
	"github.com/larse514/amazonian/service"
)

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
	containerExecutor := cf.CFExecutor{Client: svc}
	serv := service.EcsService{Executor: containerExecutor, LoadBalancer: cf.AWSElb{Client: elb}}
	ecs := cluster.Ecs{Resource: cf.Stack{Client: svc}, Executor: containerExecutor}

	//check if the cluster exists, if not create it
	if !args.ClusterExists {
		//create cluster
		clusterStruct := cluster.EcsCluster{}
		clusterStruct.DomainName = args.HostedZoneName
		clusterStruct.KeyName = args.KeyName
		clusterStruct.VpcID = args.VPC
		clusterStruct.SubnetIDs = args.SubnetIDs
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
	fmt.Println("service struct ", serviceStruct)
	fmt.Println("serviceName ", args.ServiceName)
	//attempt to create the service
	err = serv.CreateService(&ecs, serviceStruct, args.ServiceName)

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	serviceName := strings.ToLower(args.ServiceName)
	dnsName := "https://" + serviceName + "." + args.HostedZoneName
	fmt.Printf("Successfully created Container Service: %s, with url: %s \n", args.ServiceName, dnsName)
}
