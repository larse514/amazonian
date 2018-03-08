package cloud

import (
	"errors"
	"fmt"

	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/commandlineargs"
	"github.com/larse514/amazonian/service"

	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/network"
)

//Cloud is an interface representating actions to take on a Cloud provider
type Cloud interface {
	CreateDeployment(args *commandlineargs.CommandLineArgs)
}

//AWS is an implementation of Cloud interface representing AWS cloud
type AWS struct {
	Vpc   network.Network
	Stack cf.Resource
	Ecs   cluster.Cluster
	Serv  service.Service
}

//CreateDeployment is method to create networking (VPC, Subnets, etc) for AWS
func (aws AWS) CreateDeployment(args *commandlineargs.CommandLineArgs) error {
	//check if vpc exists, if not attempt to create it
	if !args.VPCExists {
		fmt.Printf("VPC doesn't exist, creating %s...\n", args.VPCName)

		err := aws.createVPC(args)

		if err != nil {
			return errors.New("Failed to create VPC")
		}
	}

	//grab the VPC outputs we need to hook into our cluster and deployment
	output, err := aws.getVPC(&args.VPCName, &args.Tenant)
	if err != nil {
		return errors.New("Error retrieving vpc " + args.VPCName)

	}
	//check if the cluster exists, if not create it
	if !args.ClusterExists {
		fmt.Printf("Cluster doesn't exist, creating %s...", args.ClusterName)

		//create aws ECS cluster
		err = aws.createCluster(&output, args)
		if err != nil {
			return errors.New("Error creating cluster " + args.ClusterName)

		}

	}
	//now get the cluster based on the stack name provided
	ecs, err := aws.Ecs.GetCluster(args.ClusterName)

	if err != nil {
		fmt.Printf("error retrieving stack %s", args.ClusterName)
		return errors.New("error retrieving stack " + args.ClusterName)
	}
	fmt.Printf("Creating service %s ...", args.ServiceName)

	err = aws.deployService(&output, &ecs, args)
	if err != nil {
		return errors.New("error deploying service " + args.ServiceName)
	}
	return nil
}

//createVPC is a private method used to create an AWS VPC based on passed in argument
func (aws AWS) createVPC(args *commandlineargs.CommandLineArgs) error {
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
func (aws AWS) getVPC(vpcName *string, tenant *string) (network.VPCOutput, error) {
	//let's grab the vpc to get needed output values
	vpcStack, err := aws.Stack.GetStack(vpcName)
	if err != nil {
		fmt.Println("error creating vpc ", err.Error())
		return network.VPCOutput{}, err
	}
	//i'm sorry, need to really refactor this whole block
	vpcID := cf.GetOutputValue(vpcStack, "VPC-"+*tenant)
	wsSubnetIDs := cf.GetOutputValue(vpcStack, "WSSubnet1-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+*tenant)
	//todo-get VPC private subnets to work
	clusterSubnetIDs := cf.GetOutputValue(vpcStack, "WSSubnet1-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet2-"+*tenant) + "," + cf.GetOutputValue(vpcStack, "WSSubnet3-"+*tenant)

	return network.VPCOutput{VPCID: vpcID, WSSubnetIDs: wsSubnetIDs, CLSubnetIDs: clusterSubnetIDs}, nil
}

//createCluster is a private method to create a cluster if it doesn't exist
func (aws AWS) createCluster(output *network.VPCOutput, args *commandlineargs.CommandLineArgs) error {

	//create cluster
	ecsInput := cluster.EcsInput{}
	ecsInput.DomainName = args.HostedZoneName
	ecsInput.KeyName = args.KeyName
	ecsInput.VpcID = output.VPCID
	ecsInput.ClusterSubnetIds = output.CLSubnetIDs
	ecsInput.WSSubnetIds = output.WSSubnetIDs
	ecsInput.DesiredCapacity = args.ClusterSize
	ecsInput.MaxSize = args.MaxSize
	ecsInput.InstanceType = args.InstanceType
	ecsInput.ClusterName = args.ClusterName
	//create cluster
	err := aws.Ecs.CreateCluster(ecsInput)
	if err != nil {
		println("error creating cluster ", err.Error())
		return err
	}
	return nil
}

//deployService is a private method to deploy an AWS ECS Service
func (aws AWS) deployService(vpc *network.VPCOutput, ecs *cluster.EcsOutput, args *commandlineargs.CommandLineArgs) error {
	//create the service struct, this is the struct that defines everything we need to create a container service
	//(note that for the time being only ECS is supported)
	serviceStruct := service.EcsServiceInput{}
	serviceStruct.Vpc = vpc.VPCID
	serviceStruct.Image = args.Image
	serviceStruct.ServiceName = args.ServiceName
	serviceStruct.ContainerName = args.ContainerName
	serviceStruct.HostedZoneName = args.HostedZoneName
	serviceStruct.PortMapping = args.PortMapping
	//attempt to create the service
	return aws.Serv.CreateService(ecs, &serviceStruct)

}
