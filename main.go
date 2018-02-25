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
	vpcParam             = "VpcId"
	priorityParam        = "Priority"
	hostedZoneNameParam  = "HostedZoneName"
	eLBHostedZoneIDParam = "ecslbhostedzoneid"
	eLBDNSNameParam      = "ecslbdnsname"
	eLBARNParam          = "ecslbarn"
	clusterARNParam      = "ecscluster"
	aLBListenerARNParam  = "alblistener"
	imageParam           = "image"
	containerTemplateURL = "https://s3.amazonaws.com/ecs.bucket.template/ecstenant/containertemplate.yml"
	ecsTemplateURL       = "https://s3.amazonaws.com/ecs.bucket.template/ecs/ecs.yml"
	serviceNameParam     = "ServiceName"
	containerNameParam   = "ContainerName"

	//export param names
	clusterArn      = "ecscluster"
	ecsHostedZoneID = "ecslbhostedzoneid"
	albListener     = "alblistener"
	ecsDNSName      = "ecslbdnsname"
	ecsLbArn        = "ecslbarn"
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
	//create the service struct, this is the struct that defines everything we need to create a container service
	//(note that for the time being only ECS is supported)
	serviceStruct := cf.EcsService{Vpc: *vpcPtr, Priority: *priorityPtr, Image: *imagePtr, ServiceName: *serviceNamePtr, ContainerName: *containerNamePtr, HostedZoneName: *hostedZonePtr}

	ecs := cluster.Ecs{}
	//check if the cluster exists, if not create it
	if !*clusterExistsPtr {
		//for now lets hardcode the ECSCluster params
		//todo-refactor to not one giant line
		clusterStruct := cf.EcsCluster{DomainName: *hostedZonePtr, KeyName: *keyNamePrt, VpcID: *vpcPtr, SubnetIDs: *subnetPrt, DesiredCapacity: *cluserSizePrt, MaxSize: *mazSizePrt, InstanceType: *instanceTypePrt}
		//create the parameters
		clusterParams := cf.CreateClusterParameters(clusterStruct)
		//initialize executor to create the cluster
		ecs = cluster.Ecs{Executor: cf.CFExecutor{Client: svc, StackName: *clusterNamePtr, TemplateURL: ecsTemplateURL, Parameters: clusterParams}}
		//create cluster
		err := ecs.CreateCluster(*clusterNamePtr)
		if err != nil {
			println("error creating cluster ", err.Error())
			os.Exit(1)
		}
	}
	//initialize ecs to retrieve cluster
	ecs = cluster.Ecs{Resource: cf.Stack{Client: svc}}

	//now get the cluster based on the stack name provided
	ecs, err = ecs.GetCluster(*clusterNamePtr)

	if err != nil {
		fmt.Printf("error retrieving stack %s", *clusterNamePtr)
		os.Exit(1)
	}
	//Grab the output parameters form the ECS Cluster that was just fetched
	ecsParameters := ecs.GetOutputParameters()

	//generate the parameters to create an ECS Service
	parameters := cf.CreateServiceParameters(ecsParameters, serviceStruct, *clusterNamePtr)

	//initialize the thing that will actually create the stack
	containerExecutor := cf.CFExecutor{Client: svc, StackName: *serviceNamePtr, TemplateURL: containerTemplateURL, Parameters: parameters}
	//take the executor and initialize the ECS Servie
	//todo-this very badly needs to be renamed
	serv := service.ECSService{Executor: containerExecutor}
	//attempt to create the service
	err = serv.CreateService()

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	serviceName := strings.ToLower(*serviceNamePtr)
	dnsName := "https://" + serviceName + "." + *hostedZonePtr
	fmt.Printf("Successfully created Container Service: %s, with url: %s \n", *serviceNamePtr, dnsName)
}
