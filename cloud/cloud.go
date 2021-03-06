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

//Cloud is an interface representing actions to take on a Cloud provider
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

	vpc, err := aws.retrieveOrCreateVPC(args)

	if err != nil {
		fmt.Println("Error retrieving or creating VPC ", args.VPCName)
		return errors.New("error retrieving or creating vpc")
	}

	ecs, err := aws.retrieveOrCreateCluster(&vpc, args)

	if err != nil {
		fmt.Println("Error retrieving or creating cluster ", args.ClusterName)
		return errors.New("error retrieving or creating cluster")
	}

	err = aws.deployService(&vpc, &ecs, args)

	if err != nil {
		fmt.Println("Error deploying service ", args.ServiceName)

		return errors.New("error deploying service")
	}
	return nil
}

//createVPC is a private method used to create an AWS VPC based on passed in argument
func (aws AWS) retrieveOrCreateVPC(args *commandlineargs.CommandLineArgs) (network.VPCOutput, error) {
	//if the cluster exists no need to look it up sicne it will also have to exist
	if args.ECSClusterExists {
		return network.VPCOutput{VPCID: args.VPC}, nil
	}
	//attempt to fetch the VPC by it's name
	vpcStack, err := aws.getVPC(&args.VPCName, &args.Tenant)
	//if stack name is the empty string then assume stack doesn't exist so create it

	if vpcStack.VPCID == "" {
		fmt.Printf("VPC %s doesn't exist, creating...\n", args.VPCName)
		//VPC doesn't exist so let's create a VPC with default secure values
		vpcInput := network.CreateDefaultVPC(args.VPCName, args.Tenant)
		//attempt to create VPC
		err = aws.Vpc.CreateNetwork(vpcInput)
		if err != nil {
			return network.VPCOutput{}, err
		}
		//attempt to grab the VPC now that we created it
		vpcStack, err = aws.getVPC(&args.VPCName, &args.Tenant)
	}

	if err != nil {
		println("error creating or retrieving vpc ", err.Error())
		return network.VPCOutput{}, err
	}

	return vpcStack, nil
}

//getVPC is a method to return vpcID, wsSubnetsIDs, and clusterSubnetIDs
func (aws AWS) getVPC(vpcName *string, tenant *string) (network.VPCOutput, error) {
	//let's grab the vpc to get needed output values
	vpcStack, err := aws.Stack.GetStack(vpcName)
	//if there was an error log and return
	if err != nil {
		fmt.Println("Get VPC returned an error")
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
func (aws AWS) retrieveOrCreateCluster(output *network.VPCOutput, args *commandlineargs.CommandLineArgs) (cluster.EcsOutput, error) {
	//attempt to lookup cluster, right now we're not handling this error so
	//we won't pretend to set it
	ecsOuput, _ := aws.Ecs.GetCluster(args.ClusterName)
	//check if the ECSClusterARN value exists
	if ecsOuput.ECSClusterArn == "" {
		//cluster doesn't exist, check if ClusterARN exists, if it does use input values
		//we're assuming validation has already been done
		if args.ECSClusterExists {
			ecsOuput.ECSDNSName = args.ECSDNSName
			ecsOuput.ECSHostedZoneID = args.ECSHostedZoneID
			ecsOuput.ECSLbArn = args.ECSALBArn
			ecsOuput.ECSLbFullName = args.ECSALBFullName
			ecsOuput.ECSClusterArn = args.ECSClusterARN
			ecsOuput.ECSAlbListener = args.ECSALBListener
		} else {
			//otherwise create the cluster
			fmt.Printf("cluster %s does not exist, creating...\n", args.ClusterName)
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
			//attempt to fetch cluster again
			ecsOuput, err = aws.Ecs.GetCluster(args.ClusterName)
			if err != nil {
				println("error creating cluster ", err.Error())
				return cluster.EcsOutput{}, err
			}
		}

	}

	return ecsOuput, nil
}

//deployService is a private method to deploy an AWS ECS Service
func (aws AWS) deployService(vpc *network.VPCOutput, ecs *cluster.EcsOutput, args *commandlineargs.CommandLineArgs) error {
	fmt.Printf("Creating service %s ...\n", args.ServiceName)

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
	return aws.Serv.DeployService(ecs, &serviceStruct)

}
