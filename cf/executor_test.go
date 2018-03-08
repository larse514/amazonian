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

func (m *mockGoodCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return nil
}
func (m *mockGoodCloudFormationClient) UpdateStack(*cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return &cloudformation.UpdateStackOutput{}, nil
}

// Define a mock to return an error.
type mockBadCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockBadCloudFormationClient) CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return nil, errors.New("Bad Error")
}
func (m *mockBadCloudFormationClient) UpdateStack(*cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return &cloudformation.UpdateStackOutput{}, errors.New("Bad Error")
}
func (m *mockBadCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return errors.New("THIS IS AN ERROR")
}

func TestCloudformationCreateStack(t *testing.T) {
	executor := CFExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.CreateStack(templateBody, stackName, nil)
	if err != nil {
		t.Log("Successful stack request return error ", err.Error())
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
func TestCloudformationUpdateStack(t *testing.T) {
	executor := CFExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.UpdateStack(templateBody, stackName, nil)
	if err != nil {
		t.Log("Successful stack request return error ", err.Error())
		t.Fail()
	}

}

func TestCloudformationUpdateStackFails(t *testing.T) {
	executor := CFExecutor{Client: &mockBadCloudFormationClient{}}

	err := executor.UpdateStack(templateBody, stackName, nil)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}

//PauseUntilCreateFinished tests, mocks, and methods

func TestCloudformationWaitUntilStackCreateComplete(t *testing.T) {
	executor := CFExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.PauseUntilCreateFinished(stackName)
	if err != nil {
		t.Log("Successful stack request return error")
		t.Fail()
	}

}
func TestCloudformationWaitUntilStackCreateCompleteFails(t *testing.T) {
	executor := CFExecutor{Client: &mockBadCloudFormationClient{}}

	err := executor.PauseUntilCreateFinished(stackName)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}

func TestCreateTagsLength(t *testing.T) {
	tags := createTags()

	if len(tags) != 1 {
		t.Log("invalid number of tags")
		t.Fail()
	}
}
func TestCreateTagsKey(t *testing.T) {
	tags := createTags()
	key := *tags[0].Key
	if key != amazonianKey {
		t.Log("invalid number key expected ", amazonianKey, " found ", key)
		t.Fail()
	}
}
func TestCreateTagsValue(t *testing.T) {
	tags := createTags()
	value := *tags[0].Value
	if value != amazonianValue {
		t.Log("invalid number value expected ", amazonianValue, " found ", value)
		t.Fail()
	}
}
