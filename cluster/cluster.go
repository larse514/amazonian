package cluster

import (
	"os"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/assets"
	"github.com/larse514/amazonian/cf"
)

const (
	//ECS Container Service consts
	priorityParam        = "Priority"
	hostedZoneNameParam  = "HostedZoneName"
	eLBHostedZoneIDParam = "ecslbhostedzoneid"
	eLBDNSNameParam      = "ecslbdnsname"
	eLBARNParam          = "ecslbarn"
	clusterARNParam      = "ecscluster"
	aLBListenerARNParam  = "alblistener"
	imageParam           = "image"
	serviceNameParam     = "ServiceName"
	containerNameParam   = "ContainerName"

	//ecs cluster consts
	domainNameParam      = "DomainName"
	keyNameParam         = "KeyName"
	subnetIDParam        = "SubnetId"
	desiredCapacityParam = "DesiredCapacity"
	maxSizeParam         = "MaxSize"
	instanceTypeParam    = "InstanceType"
	//shared consts
	vpcParam = "VpcId"

	//export param names
	clusterArn      = "ecscluster"
	ecsHostedZoneID = "ecslbhostedzoneid"
	albListener     = "alblistener"
	ecsDNSName      = "ecslbdnsname"
	ecsLbArn        = "ecslbarn"

	//
	containerTemplatePath = "ias/cloudformation/containertemplate.yml"
	ecsTemplatePath       = "ias/cloudformation/ecs.yml"
)

//Cluster interface to expose operations to work with various conainter clusters
type Cluster interface {
	GetCluster(stackName string) (Cluster, error)
	GetParameters() (map[string]string, error)
	CreateClusterNew(clusterName string, cluster EcsCluster) error
}

//Ecs is an implementation of an ECS cluster
type Ecs struct {
	Resource        cf.Resource
	Executor        cf.Executor
	StackName       string
	ClusterArn      string
	ECSHostedZoneID string
	AlbListener     string
	ECSDNSName      string
	ECSLbArn        string
}

//Parameter is an interface to defined methods to retrieve various Cloudformation template
//parameter value
// type Parameter interface {
// }

//EcsCluster is a struct which defines required files for an ECS Cluster
type EcsCluster struct {
	DomainName      string
	KeyName         string
	VpcID           string
	SubnetIDs       string
	DesiredCapacity string
	MaxSize         string
	//todo- could make this first class citizen
	InstanceType string
}

//GetCluster is a method to return an ECS cluster
//todo- should this just be refactored to a constructor-like implementation?
func (ecs Ecs) GetCluster(stackName string) (Ecs, error) {
	stack, err := ecs.Resource.GetStack(&stackName)

	if err != nil {
		println("error retrieving stack ", err.Error())
		return ecs, err
	}

	outputMap := getOutputParameters(&stack)
	//todo- I know, hard coded convention =/
	ecs.ClusterArn = outputMap[stackName]
	ecs.ECSHostedZoneID = outputMap[ecsHostedZoneID+"-"+stackName]
	ecs.ECSDNSName = outputMap[ecsDNSName+"-"+stackName]
	ecs.ECSLbArn = outputMap[ecsLbArn+"-"+stackName]
	ecs.AlbListener = outputMap[albListener+"-"+stackName]

	return ecs, nil
}

//getOutputParameters will retrieve the Ecs Cluster exported parameters
func getOutputParameters(stack *cloudformation.Stack) map[string]string {
	outputMap := make(map[string]string, 0)
	for _, output := range stack.Outputs {
		outputMap[*output.ExportName] = *output.OutputValue
	}

	return outputMap
}

//CreateCluster will create an ECS cluster
func (ecs Ecs) CreateCluster(clusterName string, cluster EcsCluster) error {
	//create the parameters
	clusterParams := createClusterParameters(cluster)
	//grab template
	ecsTemplate, err := assets.GetAsset(ecsTemplatePath)

	if err != nil {
		os.Exit(1)
	}
	ecs.StackName = clusterName
	//create the stack
	err = ecs.Executor.CreateStack(ecsTemplate, clusterName, clusterParams)
	if err != nil {
		println("Error processing create stack request ", err.Error())
		return err
	}
	//then wait
	err = ecs.Executor.PauseUntilFinished(clusterName)

	if err != nil {
		println("Error while attempting to wait for stack to finish processing ", err.Error())
		return err
	}

	return nil
}

//CreateClusterParameters will create the Parameter list to generate an ecs cluster
//todo- unit tests!!!
func createClusterParameters(cluster EcsCluster) []*cloudformation.Parameter {
	//we need to convert this (albiet awkwardly for the time being) to Cloudformation Parameters
	//we do as such first by converting everything to a key value map
	//key being the CF Param name, value is the value to provide to the cloudformation template
	parameterMap := make(map[string]string, 0)
	parameterMap[vpcParam] = cluster.VpcID
	parameterMap[domainNameParam] = cluster.DomainName
	parameterMap[keyNameParam] = cluster.KeyName
	parameterMap[subnetIDParam] = cluster.SubnetIDs
	parameterMap[desiredCapacityParam] = cluster.DesiredCapacity
	parameterMap[maxSizeParam] = cluster.MaxSize
	parameterMap[instanceTypeParam] = cluster.InstanceType
	return cf.CreateCloudformationParameters(parameterMap)

}
