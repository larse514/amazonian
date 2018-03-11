package commandlineargs

import (
	"errors"
	"flag"
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
	VPC               string
	VPCName           string
	HostedZoneName    string
	Image             string
	ServiceName       string
	ContainerName     string
	ClusterName       string
	ClusterSubnetIDs  string
	WSSubnetIDs       string
	KeyName           string
	ClusterSize       string
	MaxSize           string
	InstanceType      string
	PortMapping       string
	Tenant            string
	EcsLbHostedZoneID string
	EcsLBDNSName      string
	EcsLBARN          string
	EcsLBFullName     string
	EcsARN            string
	EcsLBListenerARN  string
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
	//validate arguments
	err := validateArguments(args.VPCName, args.Image, args.HostedZoneName, args.ServiceName, args.ContainerName, args.ClusterName, args.PortMapping)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		return CommandLineArgs{}, err
	}

	return args, nil

}

func createRandomString(starterString string) string {
	result := starterString + randomString(8)
	return result
}
func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letter = []rune(runeString)

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func createArgs() CommandLineArgs {
	//todo-refactor flags more for unit testing
	//VPC parameters
	vpcIDPtr := flag.String("VPCId", "", "VPC to deploy target group. (Required)")
	vpcNamePrt := flag.String("VPCName", createRandomString(vpc), "VPC Name to deploy target group. (Required if VPCId is not passed)")
	clusterNamePtr := flag.String("ClusterName", createRandomString(cluster), "Name ECS Cluster to use (Required)")
	elbSubnetPtr := flag.String("ELBSubnets", "", "List of VPC Subnets to deploy Elastic Load Balancers to (Required only if clusterExists is false)")
	hostedZonePtr := flag.String("HostedZoneName", "", "HostedZoneName used to create dns entry for services. (Required)")

	imagePtr := flag.String("Image", "", "Docker Repository Image (Required)")

	//service definitions
	serviceNamePtr := flag.String("ServiceName", createRandomString(service), "Name ECS Service Name (Required)")
	containerNamePtr := flag.String("ContainerName", createRandomString(container), "Name ECS Container Name (Required)")
	portMappingPtr := flag.String("PortMapping", "", "Port used by container (Required)")

	clusterSubnetsPtr := flag.String("ClusterSubnets", "", "List of VPC Subnets to deploy cluster to (Required only if clusterExists is false)")

	keyNamePrt := flag.String("KeyName", "", "Key name to use for cluster (Required only if clusterExists is false)")
	cluserSizePrt := flag.String("ClusterSize", "1", "Number of host machines for cluster (Required only if clusterExists is false)")
	maxSizePrt := flag.String("MaxSize", "1", "Max number of host machines cluster can scale to (Required only if clusterExists is false)")
	instanceTypePrt := flag.String("InstanceType", "t2.medium", "Type of machine. (Required only if clusterExists is false, defaults to t2.medium)")
	//custom cluster parameters
	ecsLBHostedZoneIDPtr := flag.String("ALBHostedZoneID", "", "Elastic Load Balancer Canonincal Hosted Zone Id, provide if using existing cluster")
	ecsLBDNSNamePtr := flag.String("ALBDNSName", "", "Elastic Load Balancer DNS Name, provide if using existing cluster")
	ecsLBDARNPtr := flag.String("ALBARN", "", "Elastic Load Balancer ARN, provide if using existing cluster")
	ecsLBFullNamePtr := flag.String("ALBFullName", "", "Elastic Load Balancer Full Name, provide if using existing cluster")
	ecsARNPtr := flag.String("ECSARN", "", "ECS Cluster ARN, provide if using existing cluster")
	ecsLBListenerARNPtr := flag.String("ECSALBListenerARN", "", "ECS Cluster ARN, provide if using existing cluster")

	//parse the values
	flag.Parse()
	args := CommandLineArgs{
		VPC:               *vpcIDPtr,
		VPCName:           *vpcNamePrt,
		HostedZoneName:    *hostedZonePtr,
		Image:             *imagePtr,
		ServiceName:       *serviceNamePtr,
		ContainerName:     *containerNamePtr,
		ClusterName:       *clusterNamePtr,
		ClusterSubnetIDs:  *clusterSubnetsPtr,
		WSSubnetIDs:       *elbSubnetPtr,
		KeyName:           *keyNamePrt,
		ClusterSize:       *cluserSizePrt,
		MaxSize:           *maxSizePrt,
		InstanceType:      *instanceTypePrt,
		PortMapping:       *portMappingPtr,
		Tenant:            *vpcNamePrt,
		EcsLbHostedZoneID: *ecsLBHostedZoneIDPtr,
		EcsLBDNSName:      *ecsLBDNSNamePtr,
		EcsLBARN:          *ecsLBDARNPtr,
		EcsLBFullName:     *ecsLBFullNamePtr,
		EcsARN:            *ecsARNPtr,
		EcsLBListenerARN:  *ecsLBListenerARNPtr,
	}
	return args
}
