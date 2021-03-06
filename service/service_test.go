package service

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
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
type mockGoodLoadBalancer struct {
	Client elbv2iface.ELBV2API
}
type mockBadLoadBalancer struct {
	Client elbv2iface.ELBV2API
}

func (lb mockGoodLoadBalancer) GetHighestPriority(listenerArn *string) (int, error) {
	return 10, nil
}
func (lb mockBadLoadBalancer) GetHighestPriority(listenerArn *string) (int, error) {
	return 10, errors.New("ERROR")
}
func (m mockGoodExecutor) CreateStack(templateBody string, sName string, parameters []*cloudformation.Parameter) error {
	if sName != stackName {
		return errors.New("INVALID STACK NAME")
	}
	return nil
}
func (m mockGoodExecutor) UpdateStack(templateBody string, sName string, parameters []*cloudformation.Parameter) error {
	if sName != stackName {
		return errors.New("INVALID STACK NAME")
	}
	return nil
}
func (m mockGoodExecutor) PauseUntilCreateFinished(stackName string) error {
	return nil
}
func (m mockGoodExecutor) PauseUntilUpdateFinished(stackName string) error {
	return nil
}

// Define a mock to fail on pause
type mockGoodCreateStackFailedPauseExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockGoodCreateStackFailedPauseExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (m mockGoodCreateStackFailedPauseExecutor) UpdateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (m mockGoodCreateStackFailedPauseExecutor) PauseUntilCreateFinished(stackName string) error {
	return errors.New("ERROR")
}
func (m mockGoodCreateStackFailedPauseExecutor) PauseUntilUpdateFinished(stackName string) error {
	return errors.New("ERROR")
}

// Define a mock to fail on Create Stack
type mockBadCreateStackExecutor struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockBadCreateStackExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return errors.New("ERROR")
}
func (m mockBadCreateStackExecutor) UpdateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return errors.New("ERROR")
}
func (m mockBadCreateStackExecutor) PauseUntilCreateFinished(stackName string) error {
	return nil
}
func (m mockBadCreateStackExecutor) PauseUntilUpdateFinished(stackName string) error {
	return nil
}

type mockGoodGetStack struct {
	Client cloudformationiface.CloudFormationAPI
}

func (m mockGoodGetStack) GetStack(stackName *string) (cloudformation.Stack, error) {
	return cloudformation.Stack{StackName: stackName}, nil
}
func TestCreateServicePasses(t *testing.T) {
	serv := EcsService{Executor: mockGoodExecutor{}, LoadBalancer: mockGoodLoadBalancer{}, Resource: mockGoodGetStack{}}
	ecs := cluster.EcsOutput{}
	service := EcsServiceInput{ServiceName: stackName}
	err := serv.DeployService(&ecs, &service)

	if err != nil {
		t.Log("Error returned when both methods returned successfully")
		t.Fail()
	}

}
func TestCreateServiceCreateStackFails(t *testing.T) {
	serv := EcsService{Executor: mockBadCreateStackExecutor{}, LoadBalancer: mockGoodLoadBalancer{}, Resource: mockGoodGetStack{}}
	ecs := cluster.EcsOutput{}
	service := EcsServiceInput{}
	err := serv.DeployService(&ecs, &service)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}
func TestCreateServicePriorityFails(t *testing.T) {
	serv := EcsService{Executor: mockGoodExecutor{}, LoadBalancer: mockBadLoadBalancer{}, Resource: mockGoodGetStack{}}
	ecs := cluster.EcsOutput{}
	service := EcsServiceInput{}
	err := serv.DeployService(&ecs, &service)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}

func TestCreateServicePauseFails(t *testing.T) {
	serv := EcsService{Executor: mockGoodCreateStackFailedPauseExecutor{}, LoadBalancer: mockGoodLoadBalancer{}, Resource: mockGoodGetStack{}}
	ecs := cluster.EcsOutput{}
	service := EcsServiceInput{}
	err := serv.DeployService(&ecs, &service)

	if err == nil {
		t.Log("Error not returned")
		t.Fail()
	}

}

// func TestUpdateServicePasses(t *testing.T) {
// 	serv := EcsService{Executor: mockGoodExecutor{}, LoadBalancer: mockGoodLoadBalancer{}, Resource: mockGoodGetStack{}}
// 	ecs := cluster.EcsOutput{}
// 	service := EcsServiceInput{ServiceName: stackName}
// 	err := serv.DeployService(&ecs, &service)

// 	if err != nil {
// 		t.Log("Error returned when both methods returned successfully")
// 		t.Fail()
// 	}

// }
