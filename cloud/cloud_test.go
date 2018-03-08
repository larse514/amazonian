package cloud

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/cluster"
	"github.com/larse514/amazonian/service"

	"github.com/larse514/amazonian/commandlineargs"

	"github.com/larse514/amazonian/network"
)

const (
	tenant      = "TENANT"
	vpc         = "VPC"
	vpcOutput   = "VPCOUTPUT"
	ws1         = "WSSubnet1"
	ws1Output   = "WSOUTPUT1"
	ws2         = "WSSubnet2"
	ws2Output   = "WSOUTPUT2"
	ws3         = "WSSubnet3"
	ws3Ouput    = "WSOUTPUT3"
	wsSubnetIds = ws1Output + "," + ws2Output + "," + ws3Ouput
)

type mockGoodVPC struct {
	network.Network
}

func (vpc mockGoodVPC) CreateNetwork(input *network.VPCInput) error {
	return nil
}

type mockBadVPC struct {
	network.Network
}

func (vpc mockBadVPC) CreateNetwork(input *network.VPCInput) error {
	return errors.New("ERROR")
}

type mockGoodStack struct {
	cf.Resource
}

func (stack mockGoodStack) GetStack(stackName *string) (cloudformation.Stack, error) {
	cfStack := cloudformation.Stack{}
	outputs := make([]*cloudformation.Output, 0)

	outputVpc := cloudformation.Output{}
	outputVpc.SetExportName(vpc + "-" + tenant)
	outputVpc.SetOutputValue(vpcOutput)
	outputs = append(outputs, &outputVpc)

	outputWs1 := cloudformation.Output{}
	outputWs1.SetExportName(ws1 + "-" + tenant)
	outputWs1.SetOutputValue(ws1Output)
	outputs = append(outputs, &outputWs1)

	outputWs2 := cloudformation.Output{}
	outputWs2.SetExportName(ws2 + "-" + tenant)
	outputWs2.SetOutputValue(ws2Output)
	outputs = append(outputs, &outputWs2)

	outputWs3 := cloudformation.Output{}
	outputWs3.SetExportName(ws3 + "-" + tenant)
	outputWs3.SetOutputValue(ws3Ouput)
	outputs = append(outputs, &outputWs3)

	cfStack.SetOutputs(outputs)

	return cfStack, nil
}

type mockBadStack struct {
	cf.Resource
}

func (stack mockBadStack) GetStack(stackName *string) (cloudformation.Stack, error) {
	return cloudformation.Stack{}, errors.New("STACK ERROR")
}

type mockGoodCluster struct {
	Resource cf.Resource
	Executor cf.Executor
}

func (client mockGoodCluster) GetCluster(stackName string) (cluster.EcsOutput, error) {
	return cluster.EcsOutput{}, nil
}
func (client mockGoodCluster) GetParameters() (map[string]string, error) {
	return nil, nil
}
func (client mockGoodCluster) CreateCluster(input cluster.EcsInput) error {
	args := createCommandLineArgs()
	output := createVPCOutput()
	//todo- could check all params
	if input.ClusterName != args.ClusterName {
		return errors.New("invalid ClusterName")
	}
	if output.VPCID != input.VpcID {
		return errors.New("invalid VpcID expected " + input.VpcID + " got " + output.VPCID)

	}
	return nil
}

type mockBadCluster struct {
	Resource cf.Resource
	Executor cf.Executor
}

func (client mockBadCluster) GetCluster(stackName string) (cluster.EcsOutput, error) {
	return cluster.EcsOutput{}, nil
}
func (client mockBadCluster) GetParameters() (map[string]string, error) {
	return nil, nil
}
func (client mockBadCluster) CreateCluster(input cluster.EcsInput) error {
	return errors.New("UNIT TEST ERROR")
}

type mockBadGetCluster struct {
	Resource cf.Resource
	Executor cf.Executor
}

func (client mockBadGetCluster) GetCluster(stackName string) (cluster.EcsOutput, error) {
	return cluster.EcsOutput{}, errors.New("GET CLUSTER FAILED")
}
func (client mockBadGetCluster) GetParameters() (map[string]string, error) {
	return nil, nil
}
func (client mockBadGetCluster) CreateCluster(input cluster.EcsInput) error {
	args := createCommandLineArgs()
	output := createVPCOutput()
	//todo- could check all params
	if input.ClusterName != args.ClusterName {
		return errors.New("invalid ClusterName")
	}
	if output.VPCID != input.VpcID {
		return errors.New("invalid VpcID expected " + input.VpcID + " got " + output.VPCID)

	}
	return nil
}

type mockGoodService struct {
	Executor     cf.Executor
	LoadBalancer cf.LoadBalancer
}

func (service mockGoodService) DeployService(ecs *cluster.EcsOutput, input *service.EcsServiceInput) error {
	expectedVPC := createVPCOutput().VPCID
	if input.Vpc != expectedVPC {
		return errors.New("Invlid VPCID returned.  Expected: " + expectedVPC + " but received: " + input.Vpc)
	}
	return nil
}

type mockBadService struct {
	Executor     cf.Executor
	LoadBalancer cf.LoadBalancer
}

func (service mockBadService) DeployService(ecs *cluster.EcsOutput, input *service.EcsServiceInput) error {
	return errors.New("UNIT TEST ERROR")
}

//helper methods
func createCommandLineArgs() *commandlineargs.CommandLineArgs {
	args := commandlineargs.CommandLineArgs{}

	args.VPC = "VPC"
	args.Image = "IMAGE"
	args.ServiceName = "SERVICENAME"
	args.ContainerName = "CONTAINRENAME"
	args.HostedZoneName = "HOSTEDNAME"
	args.PortMapping = "PORT"
	return &args
}

func createVPCOutput() *network.VPCOutput {
	return &network.VPCOutput{VPCID: "VPCOUTPUT"}
}
func createECSServiceOutput() *cluster.EcsOutput {
	return &cluster.EcsOutput{StackName: "StackName"}
}

//public method tests
func TestCreateDeployment(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockGoodCluster{}, Serv: mockGoodService{}}

	args := createCommandLineArgs()
	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err != nil {
		t.Log("Error returned when it shouldn't, ", err.Error())
		t.Fail()
	}

}
func TestCreateDeploymentVPCFails(t *testing.T) {
	cloud := AWS{Vpc: mockBadVPC{}, Stack: mockGoodStack{}, Ecs: mockGoodCluster{}, Serv: mockGoodService{}}
	expectedError := "Failed to create VPC"

	args := createCommandLineArgs()
	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}
func TestCreateDeploymentFetchStackFails(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockBadStack{}, Ecs: mockGoodCluster{}, Serv: mockGoodService{}}
	args := createCommandLineArgs()

	expectedError := "Error retrieving vpc " + args.VPCName

	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}
func TestCreateDeploymentCreateClusterFails(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockBadCluster{}, Serv: mockGoodService{}}
	args := createCommandLineArgs()

	expectedError := "Error creating cluster " + args.ClusterName

	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}
func TestCreateDeploymentGetClusterFails(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockBadGetCluster{}, Serv: mockGoodService{}}
	args := createCommandLineArgs()

	expectedError := "error retrieving stack " + args.ClusterName

	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}
func TestCreateDeploymentDeployServiceFails(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockGoodCluster{}, Serv: mockBadService{}}
	args := createCommandLineArgs()

	expectedError := "error deploying service " + args.ServiceName

	args.Tenant = tenant
	args.VPCExists = false
	args.ClusterExists = false
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}
func TestCreateVPC(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}}

	err := cloud.createVPC(createCommandLineArgs())

	if err != nil {
		t.Log("error returned ", err.Error())
		t.Fail()
	}
}
func TestCreateVPCFails(t *testing.T) {
	cloud := AWS{Vpc: mockBadVPC{}}

	err := cloud.createVPC(createCommandLineArgs())

	if err == nil {
		t.Log("error not returned ")
		t.Fail()
	}
	if err.Error() != "ERROR" {
		t.Log("invalid error message returned expected ERROR got ", err.Error())
		t.Fail()
	}
}

func TestGetVpc(t *testing.T) {
	cloud := AWS{Stack: mockGoodStack{}}
	vpcForPtr := vpc
	tenantForVpc := tenant
	output, err := cloud.getVPC(&vpcForPtr, &tenantForVpc)

	if err != nil {
		t.Log("retrieved unexepected error ", err.Error())
		t.Fail()
	}

	if output.VPCID != vpcOutput {
		t.Log("invalid vpcId expected ", vpcOutput, " got ", output.VPCID)
	}
	if output.WSSubnetIDs != wsSubnetIds {
		t.Log("invalid wsSubnetIds expected ", wsSubnetIds, " got ", output.WSSubnetIDs)
		t.Fail()

	}
	if output.CLSubnetIDs != wsSubnetIds {
		t.Log("invalid clSubnetIds expected ", wsSubnetIds, " got ", output.CLSubnetIDs)
		t.Fail()
	}

}

func TestGetVpcFails(t *testing.T) {
	cloud := AWS{Stack: mockBadStack{}}
	vpcForPtr := vpc
	tenantForVpc := tenant
	_, err := cloud.getVPC(&vpcForPtr, &tenantForVpc)

	if err == nil {
		t.Log("no error returned")
		t.Fail()
	}

}

func TestCreateCluster(t *testing.T) {
	cloud := AWS{Ecs: mockGoodCluster{}}
	args := createCommandLineArgs()
	output := createVPCOutput()
	err := cloud.createCluster(output, args)

	if err != nil {
		t.Log("error returned ", err.Error())
		t.Fail()
	}

}

func TestCreateClusterFails(t *testing.T) {
	cloud := AWS{Ecs: mockBadCluster{}}
	args := createCommandLineArgs()
	output := createVPCOutput()
	err := cloud.createCluster(output, args)

	if err == nil {
		t.Log("error not returned when it should return UNIT TEST ERROR ")
		t.Fail()
	}

	if err.Error() != "UNIT TEST ERROR" {
		t.Log("Incorrect error returned expected UNIT TEST ERROR but got: ", err.Error())
		t.Fail()
	}

}

func TestDeployService(t *testing.T) {
	cloud := AWS{Serv: mockGoodService{}}
	args := createCommandLineArgs()
	output := createECSServiceOutput()
	vpcOutput := createVPCOutput()
	err := cloud.deployService(vpcOutput, output, args)

	if err != nil {
		t.Log("error returned ", err.Error())
		t.Fail()
	}

}
func TestDeployServiceFails(t *testing.T) {
	cloud := AWS{Serv: mockBadService{}}
	args := createCommandLineArgs()
	output := createECSServiceOutput()
	vpcOutput := createVPCOutput()
	err := cloud.deployService(vpcOutput, output, args)

	if err == nil {
		t.Log("error not returned when it should return UNIT TEST ERROR ")
		t.Fail()
	}

	if err.Error() != "UNIT TEST ERROR" {
		t.Log("Incorrect error returned expected UNIT TEST ERROR but got: ", err.Error())
		t.Fail()
	}

}
