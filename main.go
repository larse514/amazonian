package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/cloud"
	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/network"
	"github.com/larse514/amazonian/output"
	"github.com/larse514/amazonian/service"
)

const (
	fileName = "amazonian-output"
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
	ecs := cluster.Ecs{Resource: cf.Stack{Client: svc}, Executor: cfExecutor}
	stack := cf.Stack{Client: svc}
	vpc := network.VPC{Executor: cfExecutor}
	serv := service.EcsService{Executor: cfExecutor, LoadBalancer: cf.AWSElb{Client: elb}, Resource: stack}

	aws := cloud.AWS{Vpc: &vpc, Stack: &stack, Ecs: ecs, Serv: &serv}

	err = aws.CreateDeployment(&args)

	if err != nil {
		fmt.Printf("Error encountered when creating deployment %f", err)
		os.Exit(1)
	}
	//format name to lower case since it's an http url
	serviceName := strings.ToLower(args.ServiceName)
	url := "https://" + serviceName + "." + args.HostedZoneName
	err = output.WriteOutputFile(output.Output{FileName: fileName, ServiceName: args.ServiceName, ClusterName: args.ClusterName, ServiceURL: url, VPCId: args.VPC, VPCName: args.VPCName})
	if err != nil {
		fmt.Println("Error writing output file ", err.Error())
	}
	fmt.Printf("Successfully created Container Service: %s, with url: %s \n", args.ServiceName, url)
}
