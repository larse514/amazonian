package service

import (
	"github.com/larse514/amazonian/cf"
)

//Service is used to create a generic Container Service
type Service interface {
	CreateService(parameterMap map[string]string) (string, error)
}

//ECSService is used to create a ECS Container Service
type ECSService struct {
	Executor cf.Executor
}

//Leaving this here as a demonstration of my plan

// type FargateService stuct {
// 	executor cf.Executor
// }

//CreateService is a method that creates a service for an ecs service
func (service ECSService) CreateService() error {
	//create the stack
	err := service.Executor.CreateStack()
	if err != nil {
		println("Error processing create stack request ", err.Error())
		return err
	}
	//then wait
	err = service.Executor.PauseUntilFinished()
	if err != nil {
		println("Error while attempting to wait for stack to finish processing ", err.Error())
		return err
	}
	return nil
}

// //CreateServiceParameters will create the Parameter list to generate a cluster service
// func (service ECSService) CreateServiceParameters() {
// outputs := getOutputsAsMapFromStack()
// ecs.StackName = stackName
// ecs.ClusterArn = outputs[clusterArn+"-"+stackName]
// ecs.HostedZoneID = outputs[hostedZoneID+"-"+stackName]
// ecs.EcsDNSName = outputs[ecsDNSName+"-"+stackName]
// ecs.EcsLBArn = outputs[ecsLbArn+"-"+stackName]
// ecs.AlbListener = outputs[albListener+"-"+stackName]
// }
