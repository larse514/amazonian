package commandlineargs

import (
	"errors"
	"flag"
	"fmt"
)

//CommandLineArgs is a struct representing items pulled from the command line
type CommandLineArgs struct {
	VPC            string
	Priority       string
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
	maxSizePrt := flag.String("MaxSize", "1", "Max number of host machines cluster can scale to (Required only if clusterExists is false)")
	instanceTypePrt := flag.String("InstanceType", "t2.medium", "Type of machine. (Required only if clusterExists is false, defaults to t2.medium)")
	//parse the values
	flag.Parse()
	//validate arguments
	err := validateArguments(*vpcPtr, *priorityPtr, *imagePtr, *imagePtr, *serviceNamePtr, *containerNamePtr, *clusterNamePtr)
	//if a required parameter is not specified, log error and exit
	if err != nil {
		flag.PrintDefaults()
		return CommandLineArgs{}, err
	}

	args := CommandLineArgs{}
	args.VPC = *vpcPtr
	args.Priority = *priorityPtr
	args.HostedZoneName = *hostedZonePtr
	args.Image = *imagePtr
	args.ServiceName = *serviceNamePtr
	args.ContainerName = *containerNamePtr
	args.ClusterName = *clusterNamePtr
	args.ClusterExists = *clusterExistsPtr
	args.SubnetIDs = *subnetPrt
	args.KeyName = *keyNamePrt
	args.ClusterSize = *cluserSizePrt
	args.MaxSize = *maxSizePrt
	args.InstanceType = *instanceTypePrt

	return args, nil

}
