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
	Vpc   network.Network
	Stack cf.Resource
	Ecs   *cluster.Ecs
	Serv  *service.EcsService
}

//CreateDeployment is method to create networking (VPC, Subnets, etc) for AWS
func (aws AWS) CreateDeployment(args *commandlineargs.CommandLineArgs) error {
	//check if vpc exists, if not attempt to create it
	if !args.VPCExists {
		err := aws.createVPC(args)

		if err != nil {
			return err
		}
	}

	//check if the cluster exists, if not create it
	if !args.ClusterExists {
		fmt.Printf("Cluster doesn't exist, creating %s...", args.ClusterName)

		//grab the VPC outputs we need to hook into our cluster
		//todo-do I mirror the getCluster operation and output a struct rather than a bunch of strings?
		//it might make sense since these are specific string, not the same thing.  That is to say
		//order is important, therefore they should be treated differently..okay i talked myself into
		//it, adding as refactor in backlog
		vpcID, wsSubnetIDs, clusterSubnetIDs, err := aws.getVPC(&args.VPCName, &args.Tenant)

		if err != nil {
			return err

		}
		//create aws ECS cluster
		err = aws.createCluster(&vpcID, &wsSubnetIDs, &clusterSubnetIDs, args)
		if err != nil {
			return err

		}

	}
	//now get the cluster based on the stack name provided
	ecs, err := aws.Ecs.GetCluster(args.ClusterName)

	if err != nil {
		fmt.Printf("error retrieving stack %s", args.ClusterName)
		os.Exit(1)
	}

	err = aws.deployService(&ecs, args)

	if err != nil {
		fmt.Printf("error creating service")
		return err
	}
	return nil
}

//createVPC is a private method used to create an AWS VPC based on passed in argument
func (aws AWS) createVPC(args *commandlineargs.CommandLineArgs) error {
	fmt.Printf("VPC doesn't exist, creating %s...", args.VPCName)
	//VPC doesn't exist so let's create a VPC with default secure values
	vpcInput := network.CreateDefaultVPC(args.VPCName, args.Tenant)
	//attempt to create VPC
	err := aws.Vpc.CreateNetwork(vpcInput)

	if err != nil {
		println("error creating vpc ", err.Error())
		return err
	}

	return nil
}

//getVPC is a method to return vpcID, wsSubnetsIDs, and clusterSubnetIDs
func (aws AWS) getVPC(vpcName *string, tenant *string) (string, string, string, error) {
	//let's grab the vpc to get needed output values
	vpcStack, err := aws.Stack.GetStack(vpcName)
	if err != nil {
		fmt.Println("error creating vpc ", err.Error())
		return "", "", "", err
	}
	//i'm sorry, need to really refactor this whole block
	vpcID := cf.GetOutputValue(vpcStack, "VPC-"+*tenant)
	wsSubnetIDs := cf.GetOutputValue(vpcStack, "WSSubnet1-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+*tenant)
	//todo-get VPC private subnets to work
	clusterSubnetIDs := cf.GetOutputValue(vpcStack, "WSSubnet1-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+*tenant)

	return vpcID, wsSubnetIDs, clusterSubnetIDs, nil
}

//createCluster is a private method to create a cluster if it doesn't exist
func (aws AWS) createCluster(vpcID *string, wsSubnetIds *string, clSubnetIds *string, args *commandlineargs.CommandLineArgs) error {

	//create cluster
	clusterStruct := cluster.EcsCluster{}
	clusterStruct.DomainName = args.HostedZoneName
	clusterStruct.KeyName = args.KeyName
	clusterStruct.VpcID = *vpcID
	clusterStruct.ClusterSubnetIds = *clSubnetIds
	clusterStruct.WSSubnetIds = *wsSubnetIds
	clusterStruct.DesiredCapacity = args.ClusterSize
	clusterStruct.MaxSize = args.MaxSize
	clusterStruct.InstanceType = args.InstanceType
	//create cluster
	err := aws.Ecs.CreateCluster(args.ClusterName, clusterStruct)
	if err != nil {
		println("error creating cluster ", err.Error())
		return err
	}
	return nil
}

//deployService is a private method to deploy an AWS ECS Service
func (aws AWS) deployService(ecs *cluster.Ecs, args *commandlineargs.CommandLineArgs) error {
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
	return aws.Serv.CreateService(ecs, serviceStruct, args.ServiceName)

}
