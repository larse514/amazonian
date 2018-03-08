package cf

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

//GetStack tests, mocks, and methods

// Define a mock to return a basic success
type mockGoodICloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}
type mockGoodICloudFormationClientError struct {
	cloudformationiface.CloudFormationAPI
}
type mockGoodICloudFormationClient2StacksReturned struct {
	cloudformationiface.CloudFormationAPI
}
type mockGoodICloudFormationClient0StacksReturned struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodICloudFormationClient) DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	outputs := createOutputs(1)
	return outputs, nil

}
func (m mockGoodICloudFormationClientError) DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	outputs := createOutputs(0)
	return outputs, errors.New("error")

}
func (m mockGoodICloudFormationClient2StacksReturned) DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	outputs := createOutputs(2)
	return outputs, nil

}
func (m mockGoodICloudFormationClient0StacksReturned) DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	outputs := createOutputs(0)
	return outputs, nil

}

//helper method
func createOutputs(numOutputs int) *cloudformation.DescribeStacksOutput {
	stack := cloudformation.Stack{}
	outputs := make([]*cloudformation.Output, 0)
	output := cloudformation.Output{}
	output.SetExportName("OUTPUTKEY")
	output.SetOutputValue("OUTPUTVALUE")
	outputs = append(outputs, &output)
	stack.SetOutputs(outputs)
	stacks := make([]*cloudformation.Stack, 0)

	for i := 0; i < numOutputs; i++ {
		stacks = append(stacks, &stack)

	}

	describeStackOutput := &cloudformation.DescribeStacksOutput{}
	describeStackOutput.SetStacks(stacks)
	return describeStackOutput
}

func TestGetStack(t *testing.T) {
	stacks := Stack{Client: mockGoodICloudFormationClient{}}
	stackName := "STACK"
	output, err := stacks.GetStack(&stackName)

	println(output.GoString(), " ", err)
	if err != nil {
		t.Log("unexpected error encountered ", err.Error())
		t.Fail()
	}
	outputKey := output.Outputs[0].ExportName
	outputValue := output.Outputs[0].OutputValue

	if *outputKey != "OUTPUTKEY" {
		t.Log("output key ", *outputKey, " is not equal to OUTPUTKEY")
		t.Fail()
	}
	if *outputValue != "OUTPUTVALUE" {
		t.Log("output key ", *outputKey, " is not equal to OUTPUTKEY")
		t.Fail()
	}
}

func TestGetStack2StacksReturned(t *testing.T) {
	stacks := Stack{Client: mockGoodICloudFormationClient2StacksReturned{}}
	stackName := "STACK"
	stack, _ := stacks.GetStack(&stackName)

	if *stack.StackName != "" {
		t.Log("Error not return as it should have for multiple stacks returned")
		t.Fail()
	}

}
func TestGetStack2StacksReturnedNoError(t *testing.T) {
	stacks := Stack{Client: mockGoodICloudFormationClientError{}}
	stackName := "STACK"
	_, err := stacks.GetStack(&stackName)

	if err != nil {
		t.Log("Error ", err.Error(), " returned when there shouldn't be")
		t.Fail()
	}

}

func TestGetStack0StacksReturned(t *testing.T) {
	stacks := Stack{Client: mockGoodICloudFormationClient0StacksReturned{}}
	stackName := "STACK"
	stack, _ := stacks.GetStack(&stackName)

	if *stack.StackName != "" {
		t.Log("Error not return as it should have for multiple stacks returned")
		t.Fail()
	}

}

// func TestGetStackError(t *testing.T) {
// 	stacks := Stack{Client: mockGoodICloudFormationClientError{}}
// 	stackName := "STACK"
// 	_, err := stacks.GetStack(&stackName)

// 	if err == nil || err.Error() != "error" {
// 		t.Log("Error not return as it should have for multiple stacks returned")
// 		t.Fail()
// 	}

// }

func TestGetOutputValue(t *testing.T) {
	stacks := createOutputs(1)

	val := GetOutputValue(*stacks.Stacks[0], "OUTPUTKEY")

	if val != "OUTPUTVALUE" {
		t.Log("invalid value ", val)
		t.Fail()
	}
}
func TestGetOutputValueNoValue(t *testing.T) {
	stacks := createOutputs(1)

	val := GetOutputValue(*stacks.Stacks[0], "NOT FOUND")

	if val != "" {
		t.Log("invalid value returned, expected empty string but got", val)
		t.Fail()
	}
}

func TestGetOutputValueMulit(t *testing.T) {
	stacks := createOutputs(1)

	val := GetOutputValue(*stacks.Stacks[0], "OUTPUTKEY") + "," + GetOutputValue(*stacks.Stacks[0], "OUTPUTKEY")

	if val != "OUTPUTVALUE,OUTPUTVALUE" {
		t.Log("invalid value ", val)
		t.Fail()
	}
}
