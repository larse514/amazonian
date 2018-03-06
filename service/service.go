package service

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/assets"
	"github.com/larse514/amazonian/cf"
	"github.com/larse514/amazonian/cluster"
)

const (
	//ECS Container Service consts
	priorityParam        = "Priority"
	hostedZoneNameParam  = "HostedZoneName"
	eLBHostedZoneIDParam = "ecslbhostedzoneid"
	eLBDNSNameParam      = "ecslbdnsname"
	eLBARNParam          = "ecslbarn"
	clusterARNParam      = "EcsInput"
	aLBListenerARNParam  = "alblistener"
	ecsLBFullNameParam   = "ecslbfullname"
	ecsClusterParam      = "ecscluster"

	imageParam         = "image"
	serviceNameParam   = "ServiceName"
	containerNameParam = "ContainerName"
	portMappingParam   = "PortMapping"

	//ecs cluster consts
	domainNameParam      = "DomainName"
	keyNameParam         = "KeyName"
	subnetIDParam        = "SubnetId"
	desiredCapacityParam = "DesiredCapacity"
	maxSizeParam         = "MaxSize"
	instanceTypeParam    = "InstanceType"
	//shared consts
	vpcParam = "VpcId"

	containerTemplatePath = "ias/cloudformation/containertemplate.yml"
)

//Service is used to create a generic Container Service
type Service interface {
	CreateService(ecs *cluster.EcsOutput, input *EcsServiceInput) error
}

//EcsService is used to create a ECS Container Service
type EcsService struct {
	Executor     cf.Executor
	LoadBalancer cf.LoadBalancer
}

//EcsServiceInput is used as input to create a ECS Container Service
type EcsServiceInput struct {
	Vpc            string
	Priority       string
	HostedZoneName string
	Image          string
	ServiceName    string
	ContainerName  string
	PortMapping    string
}

//Leaving this here as a demonstration of my plan

// type FargateService stuct {
// 	executor cf.Executor
// }

//CreateService is a method that creates a service for an ecs service
func (service EcsService) CreateService(ecs *cluster.EcsOutput, input *EcsServiceInput) error {
	//Now grab the priority
	priority, err := service.LoadBalancer.GetHighestPriority(&ecs.AlbListener)
	if err != nil {
		println("error retrieving latest priority ", err.Error())
		return err
	}

	//think through and refactor how to test this
	input.Priority = strconv.Itoa(priority + 1)

	//get the parameters
	parameters := createServiceParameters(ecs, input)
	//grab the template
	containerTemplate, err := assets.GetAsset(containerTemplatePath)
	if err != nil {
		println("error retrieving container service template ", err.Error())
		return err
	}

	//create the stack
	err = service.Executor.CreateStack(containerTemplate, input.ServiceName, parameters)
	if err != nil {
		println("Error processing create stack request ", err.Error())
		return err
	}
	//then wait
	err = service.Executor.PauseUntilFinished(input.ServiceName)
	if err != nil {
		println("Error while attempting to wait for stack to finish processing ", err.Error())
		return err
	}
	return nil
}

//CreateServiceParameters will create the Parameter list to generate a cluster service
//todo- unit tests!!!
func createServiceParameters(ecs *cluster.EcsOutput, service *EcsServiceInput) []*cloudformation.Parameter {
	//we need to convert this (albiet awkwardly for the time being) to Cloudformation Parameters
	//we do as such first by converting everything to a key value map
	//key being the CF Param name, value is the value to provide to the cloudformation template
	parameterMap := make(map[string]string, 0)
	//todo-refactor this bloody hardcoded mess
	parameterMap[vpcParam] = service.Vpc
	parameterMap[priorityParam] = service.Priority
	parameterMap[imageParam] = service.Image
	parameterMap[hostedZoneNameParam] = service.HostedZoneName
	parameterMap[serviceNameParam] = service.ServiceName
	parameterMap[containerNameParam] = service.ContainerName
	parameterMap[portMappingParam] = service.PortMapping
	//cluster mappings
	parameterMap[clusterARNParam] = ecs.ClusterArn
	parameterMap[eLBHostedZoneIDParam] = ecs.ECSHostedZoneID
	parameterMap[eLBDNSNameParam] = ecs.ECSDNSName
	parameterMap[eLBARNParam] = ecs.ECSLbArn
	parameterMap[aLBListenerARNParam] = ecs.AlbListener
	parameterMap[ecsLBFullNameParam] = ecs.ECSLbFullName
	parameterMap[ecsClusterParam] = ecs.StackName

	//now convert the key value map to a list of cloudformation.Parameter 's
	return cf.CreateCloudformationParameters(parameterMap)
}
