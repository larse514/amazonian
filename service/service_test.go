package service

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

//CreateStack tests, mocks, and methods

// Define a mock to return a basic success
type mockGoodExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodExecutor) CreateStack() error {
	return nil
}
func (m mockGoodExecutor) PauseUntilFinished() error {
	return nil
}

// Define a mock to fail on pause
type mockGoodCreateStackFailedPauseExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodCreateStackFailedPauseExecutor) CreateStack() error {
	return nil
}
func (m mockGoodCreateStackFailedPauseExecutor) PauseUntilFinished() error {
	return errors.New("ERROR")
}

// Define a mock to fail on Create Stack
type mockBadCreateStackExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockBadCreateStackExecutor) CreateStack() error {
	return errors.New("ERROR")
}
func (m mockBadCreateStackExecutor) PauseUntilFinished() error {
	return nil
}
func TestCreateServicePasses(t *testing.T) {
	serv := ECSService{Executor: mockGoodExecutor{}}

	err := serv.CreateService()

	if err != nil {
		t.Log("Error returned when both methods returned successfully")
		t.Fail()
	}

}
func TestCreateServiceCreateStackFails(t *testing.T) {
	serv := ECSService{Executor: mockBadCreateStackExecutor{}}

	err := serv.CreateService()

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}

func TestCreateServicePauseFails(t *testing.T) {
	serv := ECSService{Executor: mockGoodCreateStackFailedPauseExecutor{}}

	err := serv.CreateService()

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}
