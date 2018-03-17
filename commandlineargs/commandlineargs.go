package commandlineargs

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const (
	service    = "amazonian-service"
	vpc        = "amazonian-vpc"
	container  = "amazonian-container"
	cluster    = "amazonian-cluster"
	runeString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tenant     = "amazonian"
)

//CommandLineArgs is a struct representing items pulled from the command line
type CommandLineArgs struct {
	VPC              string
	VPCName          string
	HostedZoneName   string
	Image            string
	ServiceName      string
	ContainerName    string
	ClusterName      string
	ClusterSubnetIDs string
	WSSubnetIDs      string
	KeyName          string
	ClusterSize      string
	MaxSize          string
	InstanceType     string
	PortMapping      string
	Tenant           string
	//ECS cluster params
	ECSClusterARN    string
	ECSHostedZoneID  string
	ECSDNSName       string
	ECSALBArn        string
	ECSALBListener   string
	ECSALBFullName   string
	ECSClusterExists bool
}

//validateArguments method to validate all required command line args are specified
func validateArguments(args ...string) error {
	if args == nil {
		return errors.New("No command line args were specified")
	}
	for _, arg := range args {
		if arg == "" {
			return errors.New("Unspecified required command line args")
		}
	}
	return nil
}

//GenerateArgs is a method to parse command line arguments
func GenerateArgs() (CommandLineArgs, error) {
	args := createArgs()
	//validate general arguments
	err := validateArguments(args.VPCName, args.Image, args.HostedZoneName, args.ServiceName, args.ContainerName, args.ClusterName, args.PortMapping)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		return CommandLineArgs{}, err
	}
	//validate cluster arguments
	ecsExist, err := doesECSExist(args)
	if err != nil {
		fmt.Printf("Invalid input: If at least one existing ECS parameter is provided, all must be provided")
		flag.PrintDefaults()
		return CommandLineArgs{}, err
	}
	args.ECSClusterExists = ecsExist
	return args, nil

}

//Method to grab arguments from command line and convert them to CommandLineArgs struct
func createArgs() CommandLineArgs {
	//todo-refactor flags more for unit testing
	vpcPtr := flag.String("VPCId", "", "VPC to deploy target group. (Required)")
	vpcNamePrt := flag.String("VPCName", createRandomString(vpc), "VPC Name to deploy target group. (Required if VPCId is not passed)")
	portMappingPtr := flag.String("PortMapping", "", "Port used by container (Required)")

	hostedZonePtr := flag.String("HostedZoneName", "", "HostedZoneName used to create dns entry for services. (Required)")
	imagePtr := flag.String("Image", "", "Docker Repository Image (Required)")

	serviceNamePtr := flag.String("ServiceName", createRandomString(service), "Name ECS Service Name (Required)")
	containerNamePtr := flag.String("ContainerName", "", "Name ECS Container Name (Required)")
	clusterNamePtr := flag.String("ClusterName", createRandomString(cluster), "Name ECS Cluster to use (Required)")
	elbSubnetPtr := flag.String("ELBSubnets", "", "List of VPC Subnets to deploy Elastic Load Balancers to (Required only if clusterExists is false)")
	clusterSubnetsPtr := flag.String("ClusterSubnets", "", "List of VPC Subnets to deploy cluster to (Required only if clusterExists is false)")

	keyNamePrt := flag.String("KeyName", "", "Key name to use for cluster (Required only if clusterExists is false)")
	cluserSizePrt := flag.String("ClusterSize", "1", "Number of host machines for cluster (Required only if clusterExists is false)")
	maxSizePrt := flag.String("MaxSize", "1", "Max number of host machines cluster can scale to (Required only if clusterExists is false)")
	instanceTypePrt := flag.String("InstanceType", "t2.medium", "Type of machine. (Required only if clusterExists is false, defaults to t2.medium)")

	//Existing ECS Cluster Params
	ecsClusterARNPtr := flag.String("ECSClusterARN", "", "AWS ECS Cluster Amazon Resource Name (ARN))")
	ecsHostedZoneIDPtr := flag.String("ECSALBHostedZoneID", "", "AWS ECS Cluster Application Load Balancer Hosted Zone")
	ecsALNDNSNamePtr := flag.String("ECSALNDNSName", "", "AWS ECS Cluster Application Load Balancer DNS Name")
	ecsALBArnPtr := flag.String("ECSALBArn", "", "AWS ECS Cluster Application Load Balancer Amazon Resource Name (ARN)")
	ecsALBListenerPtr := flag.String("ECSALBListener", "", "AWS ECS Cluster Application Load Balancer Listener")
	ecsALBFullNamePtr := flag.String("ECSALBFullName", "", "AWS ECS Cluster Application Load Balancer Full Name")

	//parse the values
	flag.Parse()

	if *containerNamePtr == "" {
		*containerNamePtr = *serviceNamePtr
	}
	args := CommandLineArgs{
		VPC:              *vpcPtr,
		VPCName:          *vpcNamePrt,
		HostedZoneName:   *hostedZonePtr,
		Image:            *imagePtr,
		ServiceName:      *serviceNamePtr,
		ContainerName:    *containerNamePtr,
		ClusterName:      *clusterNamePtr,
		ClusterSubnetIDs: *clusterSubnetsPtr,
		WSSubnetIDs:      *elbSubnetPtr,
		KeyName:          *keyNamePrt,
		ClusterSize:      *cluserSizePrt,
		MaxSize:          *maxSizePrt,
		InstanceType:     *instanceTypePrt,
		PortMapping:      *portMappingPtr,
		Tenant:           *vpcNamePrt,
		ECSClusterARN:    *ecsClusterARNPtr,
		ECSHostedZoneID:  *ecsHostedZoneIDPtr,
		ECSDNSName:       *ecsALNDNSNamePtr,
		ECSALBArn:        *ecsALBArnPtr,
		ECSALBListener:   *ecsALBListenerPtr,
		ECSALBFullName:   *ecsALBFullNamePtr,
	}
	return args
}

//Helper Methods

//Method to create random string to be used to generate auto generated names based on initial name
func createRandomString(starterString string) string {
	result := starterString + randomString(8)
	return result
}

//Method to create random string
func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letter = []rune(runeString)

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

//Method to check whether one of the inputted parameters exist
func doAnyParametersExist(args ...string) bool {
	if args == nil {
		return false
	}
	for _, arg := range args {
		if arg != "" {
			return true
		}
	}
	return false
}

//method to validate ECS Parameters
func doesECSExist(args CommandLineArgs) (bool, error) {
	//first, check if any of the existing cluster parameters exist
	if doAnyParametersExist(args.ECSALBArn, args.ECSALBFullName, args.ECSALBListener,
		args.ECSClusterARN, args.ECSDNSName, args.ECSHostedZoneID) {
		//if any of them do, then we require all of them so validate
		err := validateArguments(args.ECSALBArn, args.ECSALBFullName, args.ECSALBListener,
			args.ECSClusterARN, args.ECSDNSName, args.ECSHostedZoneID)
		if err != nil {
			return false, err
		}

		return true, nil

	}
	return false, nil
}
