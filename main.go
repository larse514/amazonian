package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/larse514/amazonian/cloudformation"
	"github.com/larse514/amazonian/commandlineargs"
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

	fmt.Printf("textPtr: %s", *vpcPtr)
	cloudformation.ListStacks()
	fmt.Printf("finished listing?")
}
