package cf

import "github.com/aws/aws-sdk-go/service/cloudformation"

//CreateCloudformationParameters is a method to convert a map of parameter
//key value pairs into AWS Parameters
func CreateCloudformationParameters(parameterMap map[string]string) []*cloudformation.Parameter {
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
