package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/service"
)

const (
	containerTemplatePath = "ias/cloudformation/containertemplate.yml"
	ecsTemplatePath       = "ias/cloudformation/ecs.yml"
)

func main() {
	//get command line args
	vpcPtr := flag.String("VPC", "", "VPC to deploy target group. (Required)")
	priorityPtr := flag.String("Priority", "", "Priority use in Target Group Rules. (Required)")
	hostedZonePtr := flag.String("HostedZoneName", "", "HostedZoneName used to create dns entry for services. (Required)")
	imagePtr := flag.String("Image", "", "Docker Repository Image (Required)")
	serviceNamePtr := flag.String("ServiceName", "", "Name ECS Service Name (Required)")
	containerNamePtr := flag.String("ContainerName", "", "Name ECS Container Name (Required)")
	clusterNamePtr := flag.String("ClusterName", "", "Name ECS Cluster to use (Required)")
	clusterExistsPtr := flag.Bool("ClusterExists", false, "If cluster exists, defaults to false if not provided (Required)")
	subnetPrt := flag.String("Subnets", "", "List of VPC Subnets to deploy cluster to (Required only if clusterExists is false)")
	keyNamePrt := flag.String("KeyName", "", "Key name to use for cluster (Required only if clusterExists is false)")
	cluserSizePrt := flag.String("ClusterSize", "1", "Number of host machines for cluster (Required only if clusterExists is false)")
	mazSizePrt := flag.String("MaxSize", "1", "Max number of host machines cluster can scale to (Required only if clusterExists is false)")
	instanceTypePrt := flag.String("InstanceType", "t2.medium", "Type of machine. (Required only if clusterExists is false, defaults to t2.medium)")

	//parse the values
	flag.Parse()
	//validate arguments
	err := commandlineargs.ValidateArguments(*vpcPtr, *priorityPtr, *imagePtr, *imagePtr, *serviceNamePtr, *containerNamePtr, *clusterNamePtr)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	// or environment variables
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create CloudFormation client in region
	svc := cloudformation.New(sess)

	//Initialize the dependencies
	containerExecutor := cf.CFExecutor{Client: svc}
	serv := service.EcsService{Executor: containerExecutor}
	ecs := cluster.Ecs{Resource: cf.Stack{Client: svc}, Executor: containerExecutor}

	//check if the cluster exists, if not create it
	if !*clusterExistsPtr {
		//for now lets hardcode the ECSCluster params
		//todo-refactor to not one giant line
		clusterStruct := cluster.EcsCluster{}
		clusterStruct.DomainName = *hostedZonePtr
		clusterStruct.KeyName = *keyNamePrt
		clusterStruct.VpcID = *vpcPtr
		clusterStruct.SubnetIDs = *subnetPrt
		clusterStruct.DesiredCapacity = *cluserSizePrt
		clusterStruct.MaxSize = *mazSizePrt
		clusterStruct.InstanceType = *instanceTypePrt
		//create cluster
		err = ecs.CreateCluster(*clusterNamePtr, clusterStruct)
		if err != nil {
			println("error creating cluster ", err.Error())
			os.Exit(1)
		}
	}

	//now get the cluster based on the stack name provided
	ecs, err = ecs.GetCluster(*clusterNamePtr)

	if err != nil {
		fmt.Printf("error retrieving stack %s", *clusterNamePtr)
		os.Exit(1)
	}

	//let's get the priority for the next service
	// priority, err := cf.LoadBalancer.GetHighestPriority()
	//create the service struct, this is the struct that defines everything we need to create a container service
	//(note that for the time being only ECS is supported)
	serviceStruct := service.EcsService{}
	serviceStruct.Vpc = *vpcPtr
	serviceStruct.Priority = *priorityPtr
	serviceStruct.Image = *imagePtr
	serviceStruct.ServiceName = *serviceNamePtr
	serviceStruct.ContainerName = *containerNamePtr
	serviceStruct.HostedZoneName = *hostedZonePtr

	//attempt to create the service
	err = serv.CreateService(&ecs, serviceStruct, *serviceNamePtr)

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	serviceName := strings.ToLower(*serviceNamePtr)
	dnsName := "https://" + serviceName + "." + *hostedZonePtr
	fmt.Printf("Successfully created Container Service: %s, with url: %s \n", *serviceNamePtr, dnsName)
}
