package main

import (
	"flag"
	"fmt"
	"os"

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
	templateURL          = "https://s3.amazonaws.com/ecs.bucket.template/ecstenant/containertemplate.yml"
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
	stackNamePtr := flag.String("StackName", "", "Name of aws cloudformation stack (Required)")
	serviceNamePtr := flag.String("ServiceName", "", "Name ECS Service Name (Required)")
	containerNamePtr := flag.String("ContainerName", "", "Name ECS Container Name (Required)")
	clusterNamePtr := flag.String("ClusterName", "", "Name ECS Cluster to use (Required)")

	//parse the values
	flag.Parse()
	//validate arguments
	err := commandlineargs.ValidateArguments(*vpcPtr, *priorityPtr, *imagePtr, *imagePtr, *stackNamePtr, *serviceNamePtr, *containerNamePtr, *clusterNamePtr)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set stack name, template url
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	//todo-maybe meove this out even further
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create CloudFormation client in region
	svc := cloudformation.New(sess)
	//create the service struct
	serviceStruct := ecsService{vpc: *vpcPtr, priority: *priorityPtr, image: *imagePtr, serviceName: *serviceNamePtr, containerName: *containerNamePtr, hostedZoneName: *hostedZonePtr}
	//initialize ecs to retrieve clsuter
	ecs := cluster.Ecs{Resource: cf.Stack{Client: svc}}

	//now get the cluster based on the stack name provided
	ecs, err = ecs.GetCluster(*clusterNamePtr)

	if err != nil {
		fmt.Printf("error retrieving stack %s", *clusterNamePtr)
		os.Exit(1)
	}
	//Grab the output parameters form the ECS Cluster that was just fetched
	ecsParameters := ecs.GetOutputParameters()
	//now we need to convert this (albiet awkwardly for the time being) to Cloudformation Parameters
	//we do as such first by converting everything to a key value map
	//key being the CF Param name, value is the value to provide

	parameterMap := CreateServiceParameters(ecsParameters, serviceStruct, *clusterNamePtr)
	//now convert the key value map to a list of cloudformation.Parameter 's
	parameters := cf.CreateCloudformationParameters(parameterMap)

	//initialize the thing that will actually create the stack
	executor := cf.CFExecutor{Client: svc, StackName: *stackNamePtr, TemplateURL: templateURL, Parameters: parameters}
	//take the executor and initialize the ECS Servie
	//todo-this very badly needs to be renamed
	serv := service.ECSService{Executor: executor}
	//attempt to create the service
	err = serv.CreateService()

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	fmt.Printf("textPt r: %s", executor)
	// cloudformation.List Stacks()
	fmt.Printf("finished listing?")
}

//CreateServiceParameters will create the Parameter list to generate a cluster service
func CreateServiceParameters(outputs map[string]string, service ecsService, clusterStackName string) map[string]string {
	parameterMap := make(map[string]string, 0)
	//todo-refactor this bloody hardcoded mess
	parameterMap[vpcParam] = service.vpc
	parameterMap[priorityParam] = service.priority
	parameterMap[imageParam] = service.image
	parameterMap[hostedZoneNameParam] = service.hostedZoneName
	parameterMap[serviceNameParam] = service.serviceName
	parameterMap[containerNameParam] = service.containerName
	parameterMap[clusterARNParam] = outputs[clusterStackName]
	parameterMap[eLBHostedZoneIDParam] = outputs[ecsHostedZoneID+"-"+clusterStackName]
	parameterMap[eLBDNSNameParam] = outputs[ecsDNSName+"-"+clusterStackName]
	parameterMap[eLBARNParam] = outputs[ecsLbArn+"-"+clusterStackName]
	parameterMap[aLBListenerARNParam] = outputs[albListener+"-"+clusterStackName]

	return parameterMap
}

type ecsService struct {
	vpc            string
	priority       string
	hostedZoneName string
	image          string
	serviceName    string
	containerName  string
}
