package service

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/larse514/amazonian/cluster"
)

//CreateStack tests, mocks, and methods
const (
	templateBody = "BODY"
	stackName    = "STACKNAME"
)

// Define a mock to return a basic success
type mockGoodExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (m mockGoodExecutor) PauseUntilFinished(stackName string) error {
	return nil
}

// Define a mock to fail on pause
type mockGoodCreateStackFailedPauseExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodCreateStackFailedPauseExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (m mockGoodCreateStackFailedPauseExecutor) PauseUntilFinished(stackName string) error {
	return errors.New("ERROR")
}

// Define a mock to fail on Create Stack
type mockBadCreateStackExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockBadCreateStackExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return errors.New("ERROR")
}
func (m mockBadCreateStackExecutor) PauseUntilFinished(stackName string) error {
	return nil
}
func TestCreateServicePasses(t *testing.T) {
	serv := EcsService{Executor: mockGoodExecutor{}}
	ecs := cluster.Ecs{}
	service := EcsService{}
	err := serv.CreateService(&ecs, service, stackName)

	if err != nil {
		t.Log("Error returned when both methods returned successfully")
		t.Fail()
	}

}
func TestCreateServiceCreateStackFails(t *testing.T) {
	serv := EcsService{Executor: mockBadCreateStackExecutor{}}
	ecs := cluster.Ecs{}
	service := EcsService{}
	err := serv.CreateService(&ecs, service, stackName)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}

func TestCreateServicePauseFails(t *testing.T) {
	serv := EcsService{Executor: mockGoodCreateStackFailedPauseExecutor{}}
	ecs := cluster.Ecs{}
	service := EcsService{}
	err := serv.CreateService(&ecs, service, stackName)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}
