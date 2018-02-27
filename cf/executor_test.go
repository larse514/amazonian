package cf

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

const (
	stackName    = "STACKNAME"
	templateBody = "URL"
)

//CreateStack tests, mocks, and methods

// Define a mock to return a basic success
type mockGoodCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockGoodCloudFormationClient) CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return nil, nil
}

// Define a mock to return an error.
type mockBadCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockBadCloudFormationClient) CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return nil, errors.New("Bad Error")
}
func TestCloudformationCreateStack(t *testing.T) {
	executor := CFExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.CreateStack(templateBody, stackName, nil)
	if err != nil {
		t.Log("Successful stack request return error")
		t.Fail()
	}

}

func TestCloudformationCreateStackFails(t *testing.T) {
	executor := CFExecutor{Client: &mockBadCloudFormationClient{}}

	err := executor.CreateStack(templateBody, stackName, nil)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}

//PauseUntilFinished tests, mocks, and methods

func (m *mockGoodCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return nil
}

func (m *mockBadCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return errors.New("THIS IS AN ERROR")
}

func TestCloudformationWaitUntilStackCreateComplete(t *testing.T) {
	executor := CFExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.PauseUntilFinished(stackName)
	if err != nil {
		t.Log("Successful stack request return error")
		t.Fail()
	}

}
func TestCloudformationWaitUntilStackCreateCompleteFails(t *testing.T) {
	executor := CFExecutor{Client: &mockBadCloudFormationClient{}}

	err := executor.PauseUntilFinished(stackName)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}
