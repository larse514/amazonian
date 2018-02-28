package commandlineargs

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
)

const (
	service   = "service"
	container = "container"
	cluster   = "cluster"
)

//CommandLineArgs is a struct representing items pulled from the command line
type CommandLineArgs struct {
	VPC            string
	HostedZoneName string
	Image          string
	ServiceName    string
	ContainerName  string
	ClusterName    string
	ClusterExists  bool
	SubnetIDs      string
	KeyName        string
	ClusterSize    string
	MaxSize        string
	InstanceType   string
}

//validateArguments method to validate all required command line args are specified
func validateArguments(args ...string) error {
	if args == nil {
		return errors.New("No command line args were specified")
	}
	for _, arg := range args {
		println("arg: ", arg)
		if arg == "" {
			fmt.Printf("argument %s missing", arg)
			return errors.New("Unspecified required command line args")
		}
	}
	return nil
}

//GenerateArgs is a method to parse command line arguments
func GenerateArgs() (CommandLineArgs, error) {
	//todo-refactor flags more for unit testing
	vpcPtr := flag.String("VPC", "", "VPC to deploy target group. (Required)")
	hostedZonePtr := flag.String("HostedZoneName", "", "HostedZoneName used to create dns entry for services. (Required)")
	imagePtr := flag.String("Image", "", "Docker Repository Image (Required)")
	serviceNamePtr := flag.String("ServiceName", createRandomString(service), "Name ECS Service Name (Required)")
	containerNamePtr := flag.String("ContainerName", createRandomString(container), "Name ECS Container Name (Required)")
	clusterNamePtr := flag.String("ClusterName", createRandomString(cluster), "Name ECS Cluster to use (Required)")
	clusterExistsPtr := flag.Bool("ClusterExists", false, "If cluster exists, defaults to false if not provided")
	subnetPrt := flag.String("Subnets", "", "List of VPC Subnets to deploy cluster to (Required only if clusterExists is false)")
	keyNamePrt := flag.String("KeyName", "", "Key name to use for cluster (Required only if clusterExists is false)")
	cluserSizePrt := flag.String("ClusterSize", "1", "Number of host machines for cluster (Required only if clusterExists is false)")
	maxSizePrt := flag.String("MaxSize", "1", "Max number of host machines cluster can scale to (Required only if clusterExists is false)")
	instanceTypePrt := flag.String("InstanceType", "t2.medium", "Type of machine. (Required only if clusterExists is false, defaults to t2.medium)")
	//parse the values
	println(*serviceNamePtr)
	flag.Parse()
	args := CommandLineArgs{
		VPC:            *vpcPtr,
		HostedZoneName: *hostedZonePtr,
		Image:          *imagePtr,
		ServiceName:    *serviceNamePtr,
		ContainerName:  *containerNamePtr,
		ClusterName:    *clusterNamePtr,
		ClusterExists:  *clusterExistsPtr,
		SubnetIDs:      *subnetPrt,
		KeyName:        *keyNamePrt,
		ClusterSize:    *cluserSizePrt,
		MaxSize:        *maxSizePrt,
		InstanceType:   *instanceTypePrt,
	}
	fmt.Println(args)
	//validate arguments
	err := validateArguments(*vpcPtr, *imagePtr, *hostedZonePtr)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		return CommandLineArgs{}, err
	}

	return args, nil

}

func createRandomString(starterString string) string {
	result := starterString + strconv.Itoa(rand.Intn(900000))
	return result
}
