package main

import (
	"flag"
	"fmt"
	"os"

	// "github.com/larse514/amazonian/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/service"
)

const (
	vpcParam             = "VPC"
	priorityParam        = "Priority"
	hostedZoneNameParam  = "HostedZoneName"
	eLBHostedZoneIDParam = "ELBHostedZoneId"
	eLBDNSNameParam      = "ELBDNSName"
	eLBARNParam          = "ELBARN"
	clusterARNParam      = "ClusterARN"
	aLBListenerARNParam  = "ALBListenerARN"
	imageParam           = "Image"
	stackName            = "default-stack"
	templateURL          = "default-url"
)

func main() {
	//get command line args
	vpcPtr := flag.String("VPC", "", "VPC to deploy target group. (Required)")
	//todo-remove this and add dynamic lookup
	priorityPtr := flag.String("Priority", "", "Priority use in Target Group Rules. (Required)")
	hostedZonePtr := flag.String("HostedZoneName", "", "HostedZoneName used to create dns entry for services. (Required)")
	elbHostedZoneIDPtr := flag.String("ELBHostedZoneId", "", "ELBHostedZoneId used to lookup dns entry of loadbalancer for DNS entries. (Required)")
	elbDNSNamePtr := flag.String("ELBDNSName", "", "ELBDNSName used to lookup dns entry of loadbalancer for DNS entries. (Required)")
	elbARNPtr := flag.String("ELBARN", "", "ELBARN used to reference load balancer. (Required)")
	clusterArnPtr := flag.String("ClusterARN", "", "ARN of Cluster to be used to run containers. (Required)")
	albListernArnPtr := flag.String("ALBListenerARN", "", "ALB Listener Arn. (Required)")
	image := flag.String("Image", "", "Docker Repository Image (Required)")
	//parse the values
	flag.Parse()
	//validate arguments
	err := commandlineargs.ValidateArguments(*vpcPtr, *priorityPtr, *hostedZonePtr, *elbHostedZoneIDPtr, *elbDNSNamePtr, *elbARNPtr, *clusterArnPtr, *albListernArnPtr, *image)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
	//just brute force create the map we need, todo- probably refactor to a file we read in?
	parameterMap := make(map[string]string, 0)

	parameterMap[vpcParam] = *vpcPtr
	parameterMap[priorityParam] = *priorityPtr
	parameterMap[hostedZoneNameParam] = *hostedZonePtr
	parameterMap[eLBHostedZoneIDParam] = *elbHostedZoneIDPtr
	parameterMap[eLBDNSNameParam] = *elbDNSNamePtr
	parameterMap[eLBARNParam] = *elbARNPtr
	parameterMap[clusterARNParam] = *clusterArnPtr
	parameterMap[aLBListenerARNParam] = *albListernArnPtr
	parameterMap[imageParam] = *image

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

	parameters := cf.CreateCloudformationParameters(parameterMap)

	executor := cf.CFExecutor{Client: svc, StackName: stackName, TemplateURL: templateURL, Parameters: parameters}
	serv := service.ECSService{Executor: executor}
	err = serv.CreateService()

	if err != nil {
		fmt.Printf("error creating service")
		os.Exit(1)
	}
	fmt.Printf("textPt r: %s", executor)
	// cloudformation.List Stacks()
	fmt.Printf("finished listing?")
}
