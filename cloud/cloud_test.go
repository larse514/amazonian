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
	cfStack.SetStackName(*stackName)

	return cfStack, nil
}

type mockBadStack struct {
	cf.Resource
}

func (stack mockBadStack) GetStack(stackName *string) (cloudformation.Stack, error) {
	return cloudformation.Stack{}, errors.New("STACK ERROR")
}

type mockEmptyStringStack struct {
	cf.Resource
}

func (stack mockEmptyStringStack) GetStack(stackName *string) (cloudformation.Stack, error) {
	emptyString := ""
	return cloudformation.Stack{StackName: &emptyString}, nil
}

type mockEmptyStackName struct {
	cf.Resource
}

func (stack mockEmptyStackName) GetStack(stackName *string) (cloudformation.Stack, error) {
	emptyStack := ""
	return cloudformation.Stack{StackName: &emptyStack}, nil
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
	err := cloud.CreateDeployment(args)

	if err != nil {
		t.Log("Error returned when it shouldn't, ", err.Error())
		t.Fail()
	}

}
func TestCreateDeploymentECSProvided(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockBadGetCluster{}, Serv: mockGoodService{}}

	args := createCommandLineArgs()

	args.ECSClusterExists = true

	args.Tenant = tenant
	err := cloud.CreateDeployment(args)

	if err != nil {
		t.Log("Error returned when it shouldn't, ", err.Error())
		t.Fail()
	}

}
func TestCreateDeploymentECSParamsProvided(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockGoodCluster{}, Serv: mockGoodService{}}

	args := createCommandLineArgs()
	args.Tenant = tenant
	err := cloud.CreateDeployment(args)

	if err != nil {
		t.Log("Error returned when it shouldn't, ", err.Error())
		t.Fail()
	}

}
func TestCreateDeploymentVPCFails(t *testing.T) {
	cloud := AWS{Vpc: mockBadVPC{}, Stack: mockEmptyStringStack{}, Ecs: mockGoodCluster{}, Serv: mockGoodService{}}
	expectedError := "error retrieving or creating vpc"

	args := createCommandLineArgs()
	args.Tenant = tenant
	err := cloud.CreateDeployment(args)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()

	} else {
		if err.Error() != expectedError {
			t.Log("Error returned ", err.Error(), " expected ", expectedError)
			t.Fail()
		}
	}

}
func TestCreateDeploymentClusterFails(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockBadCluster{}, Serv: mockGoodService{}}
	args := createCommandLineArgs()

	args.Tenant = tenant
	err := cloud.CreateDeployment(args)

	if err != nil {
		t.Log("Error returned ", err.Error(), " when it shouldn't have")
		t.Fail()
	}

}
func TestCreateDeploymentDeployServiceFailss(t *testing.T) {
	cloud := AWS{Vpc: mockGoodVPC{}, Stack: mockGoodStack{}, Ecs: mockGoodCluster{}, Serv: mockBadService{}}
	args := createCommandLineArgs()

	expectedError := "error deploying service"

	args.Tenant = tenant
	err := cloud.CreateDeployment(args)

	if err.Error() != expectedError {
		t.Log("Error returned ", err.Error(), " expected ", expectedError)
		t.Fail()
	}

}

//retrieveOrCreateVPC tests
func TestRetrieveOrCreateVPC(t *testing.T) {
	cloud := AWS{}
	args := commandlineargs.CommandLineArgs{ECSClusterExists: true, VPC: vpc}
	expected := vpc
	got, _ := cloud.retrieveOrCreateVPC(&args)

	if got.VPCID != expected {
		t.Log("got ", got, " expcted ", expected)
		t.Fail()
	}

}
