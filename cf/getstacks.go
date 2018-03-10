package cf

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

//Resource is a generic interface for retrieving information on a infrastruture stack
type Resource interface {
	GetStack(stackName *string) (cloudformation.Stack, error)
}

//Stack is a struct representing an AWS Cloudformation Stack
type Stack struct {
	Client cloudformationiface.CloudFormationAPI
}

//GetStack is a method to retrieve an AWS stack by stack name
func (stack Stack) GetStack(stackName *string) (cloudformation.Stack, error) {
	input := &cloudformation.DescribeStacksInput{StackName: stackName}

	output, err := stack.Client.DescribeStacks(input)

	if err != nil {
		//todo- remove this until i handle errors better
		// println("error: ", err.Error(), " received when trying to find stack: ", *stackName)
		emptyStackName := ""
		return cloudformation.Stack{StackName: &emptyStackName}, nil
	}

	stackLength := len(output.Stacks)

	if stackLength != 1 {
		println("invalid number of stacks returned.  Number was: ", stackLength, " should be 1")
		emptyStackName := ""
		return cloudformation.Stack{StackName: &emptyStackName}, err
	}

	return *output.Stacks[0], nil
}

//GetOutputValue method will retrieve an output value from Output array
func GetOutputValue(stack cloudformation.Stack, key string) string {
	for i := range stack.Outputs {
		if *stack.Outputs[i].ExportName == key {
			// Found!
			return *stack.Outputs[i].OutputValue
		}
	}
	return ""
}
