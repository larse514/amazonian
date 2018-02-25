package cf

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
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
)

//Parameter is an interface to defined methods to retrieve various Cloudformation template
//parameter value
// type Parameter interface {
// }

//CreateCloudformationParameters is a method to convert a map of parameter
//key value pairs into AWS Parameters
func createCloudformationParameters(parameterMap map[string]string) []*cloudformation.Parameter {
	//initialize parameter slice
	parameters := make([]*cloudformation.Parameter, 0)

	//loop over the map and create a parameter object, then add it to the slice
	for key, value := range parameterMap {

		parameter := new(cloudformation.Parameter)
		parameter.SetParameterKey(key)
		parameter.SetParameterValue(value)
		parameters = append(parameters, parameter)
	}
	return parameters
}

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

//EcsService is a struct which defines required fields for an ECS Service
type EcsService struct {
	Vpc            string
	Priority       string
	HostedZoneName string
	Image          string
	ServiceName    string
	ContainerName  string
}

//CreateClusterParameters will create the Parameter list to generate an ecs cluster
//todo- unit tests!!!
func CreateClusterParameters(cluster EcsCluster) []*cloudformation.Parameter {
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
	return createCloudformationParameters(parameterMap)

}

//CreateServiceParameters will create the Parameter list to generate a cluster service
//todo- unit tests!!!
func CreateServiceParameters(outputs map[string]string, service EcsService, clusterStackName string) []*cloudformation.Parameter {
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
	parameterMap[clusterARNParam] = outputs[clusterStackName]
	parameterMap[eLBHostedZoneIDParam] = outputs[ecsHostedZoneID+"-"+clusterStackName]
	parameterMap[eLBDNSNameParam] = outputs[ecsDNSName+"-"+clusterStackName]
	parameterMap[eLBARNParam] = outputs[ecsLbArn+"-"+clusterStackName]
	parameterMap[aLBListenerARNParam] = outputs[albListener+"-"+clusterStackName]
	//now convert the key value map to a list of cloudformation.Parameter 's
	return createCloudformationParameters(parameterMap)
}
