package cloud

import (
	"fmt"
	"os"

	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/service"

	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/network"
)

//Cloud is an interface representating actions to take on a Cloud provider
type Cloud interface {
	CreateDeployment()
}

//AWS is an implementation of Cloud interface representing AWS cloud
type AWS struct {
	Vpc   *network.VPC
	Stack *cf.Stack
	Ecs   *cluster.Ecs
	Serv  *service.EcsService
}

//CreateDeployment is method to create networking (VPC, Subnets, etc) for AWS
func (aws AWS) CreateDeployment(args commandlineargs.CommandLineArgs) error {
	if !args.VPCExists {
		fmt.Printf("VPC doesn't exist, creating %s...", args.VPCName)
		vpcInput := network.CreateDefaultVPC(args.VPCName, args.Tenant)
		err := aws.Vpc.CreateNetwork(vpcInput)
		if err != nil {
			println("error creating vpc ", err.Error())
			os.Exit(1)
		}
	}
	if !args.VPCExists && args.VPC == "" {
		//let's grab the vpc to get needed output values
		vpcStack, err := aws.Stack.GetStack(&args.VPCName)
		if err != nil {
			fmt.Println("error creating vpc ", err.Error())
			os.Exit(1)
		}
		//i'm sorry, need to really refactor this whole block
		args.VPC = cf.GetOutputValue(vpcStack, "VPC-"+args.Tenant)
		args.WSSubnetIDs = cf.GetOutputValue(vpcStack, "WSSubnet1-"+args.Tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+args.Tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+args.Tenant)
		//todo-get VPC private subnets to work
		args.ClusterSubnetIDs = cf.GetOutputValue(vpcStack, "WSSubnet1-"+args.Tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+args.Tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+args.Tenant)
	}
	//check if the cluster exists, if not create it
	if !args.ClusterExists {
		fmt.Printf("Cluster doesn't exist, creating %s...", args.ClusterName)

		//create cluster
		clusterStruct := cluster.EcsCluster{}
		clusterStruct.DomainName = args.HostedZoneName
		clusterStruct.KeyName = args.KeyName
		clusterStruct.VpcID = args.VPC
		clusterStruct.ClusterSubnetIds = args.ClusterSubnetIDs
		clusterStruct.WSSubnetIds = args.WSSubnetIDs
		clusterStruct.DesiredCapacity = args.ClusterSize
		clusterStruct.MaxSize = args.MaxSize
		clusterStruct.InstanceType = args.InstanceType
		//create cluster
		err := aws.Ecs.CreateCluster(args.ClusterName, clusterStruct)
		if err != nil {
			println("error creating cluster ", err.Error())
			return err
		}
		//now get the cluster based on the stack name provided
		ecs, err := aws.Ecs.GetCluster(args.ClusterName)

		if err != nil {
			fmt.Printf("error retrieving stack %s", args.ClusterName)
			os.Exit(1)
		}
		//let's get the priority for the next service
		// priority, err := cf.LoadBalancer.GetHighestPriority()
		//create the service struct, this is the struct that defines everything we need to create a container service
		//(note that for the time being only ECS is supported)
		serviceStruct := service.EcsService{}
		serviceStruct.Vpc = args.VPC
		serviceStruct.Image = args.Image
		serviceStruct.ServiceName = args.ServiceName
		serviceStruct.ContainerName = args.ContainerName
		serviceStruct.HostedZoneName = args.HostedZoneName
		serviceStruct.PortMapping = args.PortMapping
		//attempt to create the service
		fmt.Printf("Creating service %s ...", args.ServiceName)
		err = aws.Serv.CreateService(&ecs, serviceStruct, args.ServiceName)

		if err != nil {
			fmt.Printf("error creating service")
			return err
		}

	}
	return nil
}
