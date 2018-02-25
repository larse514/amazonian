package cluster

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
)

const ()

//Cluster interface to expose operations to work with various conainter clusters
type Cluster interface {
	GetCluster(stackName string) (Cluster, error)
	GetParameters() (map[string]string, error)
	CreateCluster(clusterName string) error
}

//Ecs is an implementation of an ECS cluster
type Ecs struct {
	Resource  cf.Resource
	Executor  cf.Executor
	StackName string
	stack     cloudformation.Stack
}

//GetCluster is a method to return an ECS cluster
//todo- should this just be refactored to a constructor-like implementation?
func (ecs Ecs) GetCluster(stackName string) (Ecs, error) {
	stack, err := ecs.Resource.GetStack(&stackName)

	if err != nil {
		println("error retrieving stack ", err.Error())
		return ecs, err
	}
	println("fetched stack ", stack.GoString())
	ecs.stack = stack

	return ecs, nil
}

//GetOutputParameters will retrieve the Ecs Cluster exported parameters
func (ecs Ecs) GetOutputParameters() map[string]string {

	outputMap := make(map[string]string, 0)
	println("about to convert outputs from ", ecs.stack.GoString())
	for _, output := range ecs.stack.Outputs {
		println("output ", output.GoString())
		outputMap[*output.ExportName] = *output.OutputValue
	}

	return outputMap
}

//CreateCluster will create an ECS cluster
func (ecs Ecs) CreateCluster(clusterName string) error {
	ecs.StackName = clusterName

	//create the stack
	err := ecs.Executor.CreateStack()
	if err != nil {
		println("Error processing create stack request ", err.Error())
		return err
	}
	//then wait
	err = ecs.Executor.PauseUntilFinished()
	if err != nil {
		println("Error while attempting to wait for stack to finish processing ", err.Error())
		return err
	}

	return nil
}
