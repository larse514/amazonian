package cloud

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"

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

//helper methods

func createCommandLineArgs() *commandlineargs.CommandLineArgs {
	return &commandlineargs.CommandLineArgs{VPC: "VPC"}
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
	vpcID, ws, cl, err := cloud.getVPC(&vpcForPtr, &tenantForVpc)

	if err != nil {
		t.Log("retrieved unexepected error ", err.Error())
		t.Fail()
	}

	if vpcID != vpcOutput {
		t.Log("invalid vpcId expected ", vpcOutput, " got ", vpcID)
	}
	if ws != wsSubnetIds {
		t.Log("invalid wsSubnetIds expected ", wsSubnetIds, " got ", ws)
		t.Fail()

	}
	if cl != wsSubnetIds {
		t.Log("invalid clSubnetIds expected ", wsSubnetIds, " got ", cl)
		t.Fail()
	}

}

func TestGetVpcFails(t *testing.T) {
	cloud := AWS{Stack: mockBadStack{}}
	vpcForPtr := vpc
	tenantForVpc := tenant
	_, _, _, err := cloud.getVPC(&vpcForPtr, &tenantForVpc)

	if err == nil {
		t.Log("no error returned")
		t.Fail()
	}

}
