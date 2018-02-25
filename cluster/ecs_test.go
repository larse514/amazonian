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

type mockGoodResource struct {
	cf.Resource
}

func (m mockGoodResource) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := createStack(1, *stackName)
	return *outputs, nil

}

type mockGoodResourceNoStacks struct {
	cf.Resource
}

func (m mockGoodResourceNoStacks) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := createStack(0, *stackName)
	return *outputs, nil

}

type mockBadResource struct {
	cf.Resource
}

func (m mockBadResource) GetStack(stackName *string) (cloudformation.Stack, error) {
	outputs := createStack(1, *stackName)
	return *outputs, errors.New("THIS IS AN ERROR")

}
func createStack(numOutputs int, stackName string) *cloudformation.Stack {
	output := cloudformation.Output{}
	output.SetOutputKey("ecscluster-" + stackName)
	output.SetOutputValue(stackValue)
	outputs := make([]*cloudformation.Output, 0)

	for i := 0; i < numOutputs; i++ {
		outputs = append(outputs, &output)

	}

	stack := cloudformation.Stack{}
	stack.SetOutputs(outputs)
	return &stack
}

func TestEcsGetCluster(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResource{}}

	ecs, _ = ecs.GetCluster(stackName)
	if *ecs.stack.Outputs[0].OutputKey != "ecscluster-STACKNAME" {
		t.Log("mismatch in output key ", ecs.stack.GoString())
		t.Fail()
	}
	if *ecs.stack.Outputs[0].OutputValue != stackValue {
		t.Log("mismatch in output key ", ecs.stack.GoString())
		t.Fail()
	}

}
func TestEcsGetClusterGetParameters(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResource{}}

	ecs, _ = ecs.GetCluster(stackName)
	if *ecs.stack.Outputs[0].OutputKey != "ecscluster-STACKNAME" {
		t.Log("mismatch in output key ", ecs.stack.GoString())
		t.Fail()
	}
	if *ecs.stack.Outputs[0].OutputValue != stackValue {
		t.Log("mismatch in output key ", ecs.stack.GoString())
		t.Fail()
	}
	paramMap := ecs.GetOutputParameters()

	if paramMap["ecscluster-STACKNAME"] != stackValue {
		t.Log("Param map incorrectly created ", paramMap)
		t.Fail()
	}

}
func TestEcsGetClusterGetParametersNoStacks(t *testing.T) {

	ecs := Ecs{Resource: mockGoodResourceNoStacks{}}

	ecs, _ = ecs.GetCluster(stackName)
	if len(ecs.stack.Outputs) != 0 {
		t.Log("no stacks should have been returned ", ecs.stack.GoString())
		t.Fail()
	}
	paramMap := ecs.GetOutputParameters()

	if len(paramMap) != 0 {
		t.Log("Param map incorrectly created ", paramMap)
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
