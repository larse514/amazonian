package cluster

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
)

const (
	stackName  = "STACKNAME"
	stackValue = "VALUE-" + stackName
	clusterArn = "EcsInput"
)

var outputParams = []struct {
	exportName  string
	outputValue string
}{
	{stackName, clusterArn},
	{ecsHostedZoneID + "-" + stackName, ecsHostedZoneID},
	{albListener + "-" + stackName, albListener},
	{ecsDNSName + "-" + stackName, ecsDNSName},
	{ecsLbArn + "-" + stackName, ecsLbArn},
}

type mockGoodResource struct {
	cf.Resource
}
type mockGoodExecutor struct {
	cf.Executor
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor mockGoodExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (executor mockGoodExecutor) PauseUntilCreateFinished(stackName string) error {
	return nil
}
func (m mockGoodResource) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := createStack(*stackName)
	return *outputs, nil

}

type mockGoodResourceNoStacks struct {
	cf.Resource
}

func (m mockGoodResourceNoStacks) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := cloudformation.Stack{}
	return outputs, nil

}

type mockBadResource struct {
	cf.Resource
}

func (m mockBadResource) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := createStack(*stackName)
	return *outputs, errors.New("THIS IS AN ERROR")

}
func createStack(stackName string) *cloudformation.Stack {

	outputs := make([]*cloudformation.Output, 0)

	for i := 0; i < len(outputParams); i++ {
		output := cloudformation.Output{}
		output.SetExportName(outputParams[i].exportName)
		output.SetOutputValue(outputParams[i].outputValue)
		outputs = append(outputs, &output)
	}

	stack := cloudformation.Stack{}
	stack.SetOutputs(outputs)
	return &stack
}

func TestEcsGetCluster(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResource{}}

	output, _ := ecs.GetCluster(stackName)

	if output.ClusterArn != outputParams[0].outputValue {
		t.Log("mismatch in output key ", output.ClusterArn, " with ", outputParams[0].outputValue)
		t.Fail()
	}
	if output.ECSHostedZoneID != outputParams[1].outputValue {
		t.Log("mismatch in output key ", output.ECSHostedZoneID, " with ", outputParams[1].outputValue)
		t.Fail()
	}
	if output.AlbListener != outputParams[2].outputValue {
		t.Log("mismatch in output key ", output.AlbListener, " with ", outputParams[2].outputValue)
		t.Fail()
	}
	if output.ECSDNSName != outputParams[3].outputValue {
		t.Log("mismatch in output key ", output.ECSDNSName, " with ", outputParams[3].outputValue)
		t.Fail()
	}
	if output.ECSLbArn != outputParams[4].outputValue {
		t.Log("mismatch in output key ", output.ECSLbArn, " with ", outputParams[4].outputValue)
		t.Fail()
	}

}

func TestEcsGetClusterGetParametersNoStacks(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResourceNoStacks{}}

	output, _ := ecs.GetCluster(stackName)
	if output.ClusterArn != "" {
		t.Log("no stacks have been returned, ClusterArn should be \"\"", output.ClusterArn)
		t.Fail()
	}

}
func TestEcsGetClusterFails(t *testing.T) {

	ecs := Ecs{Resource: mockBadResource{}}

	output, _ := ecs.GetCluster(stackName)

	if output.ClusterArn != "" {
		t.Log("ClusterARN not equal to emptry string")
		t.Fail()
	}

}
func TestEcsGetClusterFailsErrorNil(t *testing.T) {

	ecs := Ecs{Resource: mockBadResource{}}

	_, err := ecs.GetCluster(stackName)

	if err != nil {
		t.Log("Error not nil")
		t.Fail()
	}

}
func TestEcsCreateCluster(t *testing.T) {
	clusterStruct := EcsInput{ClusterName: stackName, DomainName: "DOMAIN", KeyName: "KEY", VpcID: "VPC", WSSubnetIds: "WSSUBNETS", ClusterSubnetIds: "ClusterSubnetIds", DesiredCapacity: "CAPACITY", MaxSize: "MAXSIZE", InstanceType: "INSTANCETYPE"}

	ecs := Ecs{Executor: mockGoodExecutor{}}

	err := ecs.CreateCluster(clusterStruct)

	if err != nil {
		t.Log("error is not nil when it should be ", err.Error())
		t.Fail()
	}

}

//CreateClusterParameters tests
func TestCreateClusterParameters(t *testing.T) {
	clusterStruct := EcsInput{DomainName: "DOMAIN", KeyName: "KEY", VpcID: "VPC", WSSubnetIds: "WSSUBNETS", ClusterSubnetIds: "ClusterSubnetIds", DesiredCapacity: "CAPACITY", MaxSize: "MAXSIZE", InstanceType: "INSTANCETYPE"}

	params := createClusterParameters(clusterStruct)
	len := len(params)
	if len != 8 {
		t.Log("expected parameter length to be 8 but found ", len)
		t.Fail()

	}

}
