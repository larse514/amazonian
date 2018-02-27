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
func (executor mockGoodExecutor) PauseUntilFinished(stackName string) error {
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

	ecs, _ = ecs.GetCluster(stackName)

	if ecs.ClusterArn != outputParams[0].outputValue {
		t.Log("mismatch in output key ", ecs.ClusterArn, " with ", outputParams[0].outputValue)
		t.Fail()
	}
	if ecs.ECSHostedZoneID != outputParams[1].outputValue {
		t.Log("mismatch in output key ", ecs.ECSHostedZoneID, " with ", outputParams[1].outputValue)
		t.Fail()
	}
	if ecs.AlbListener != outputParams[2].outputValue {
		t.Log("mismatch in output key ", ecs.AlbListener, " with ", outputParams[2].outputValue)
		t.Fail()
	}
	if ecs.ECSDNSName != outputParams[3].outputValue {
		t.Log("mismatch in output key ", ecs.ECSDNSName, " with ", outputParams[3].outputValue)
		t.Fail()
	}
	if ecs.ECSLbArn != outputParams[4].outputValue {
		t.Log("mismatch in output key ", ecs.ECSLbArn, " with ", outputParams[4].outputValue)
		t.Fail()
	}

}

func TestEcsGetClusterGetParametersNoStacks(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResourceNoStacks{}}

	ecs, _ = ecs.GetCluster(stackName)
	if ecs.ClusterArn != "" {
		t.Log("no stacks have been returned, ClusterArn should be \"\"", ecs.ClusterArn)
		t.Fail()
	}

}
func TestEcsGetClusterFails(t *testing.T) {

	ecs := Ecs{Resource: mockBadResource{}}

	ecs, err := ecs.GetCluster(stackName)

	if err == nil {
		t.Log("error is nil when it shouldn't be")
		t.Fail()
	}
	if err.Error() != "THIS IS AN ERROR" {
		t.Log("Error message: ", err.Error(), " is invalid")
		t.Fail()
	}

}
func TestEcsCreateCluster(t *testing.T) {
	clusterStruct := EcsCluster{DomainName: "DOMAIN", KeyName: "KEY", VpcID: "VPC", SubnetIDs: "SUBNETS", DesiredCapacity: "CAPACITY", MaxSize: "MAXSIZE", InstanceType: "INSTANCETYPE"}

	ecs := Ecs{Executor: mockGoodExecutor{}}

	err := ecs.CreateCluster(stackName, clusterStruct)
	println("jere")

	if err != nil {
		t.Log("error is not nil when it should be ", err.Error())
		t.Fail()
	}

}

//CreateClusterParameters tests
func TestCreateClusterParameters(t *testing.T) {
	cluster := EcsCluster{DomainName: "DOMAINAME"}

	params := createClusterParameters(cluster)

	if *params[0].ParameterKey != "DomainName" {
		t.Log("paramkey ", params[0].ParameterKey, " did not get set to correct constant value")
	}
	if *params[0].ParameterValue != "DOMAINAME" {
		t.Log("paramvalue ", params[0].ParameterValue, " did not get set to correct constant value")
	}
}
