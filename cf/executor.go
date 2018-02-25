package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"fmt"
)

//Executor is an interface to execute and create stacks
type Executor interface {
	CreateStack() error
	PauseUntilFinished() error
}

//CFExecutor struct used to create cloudformation stacks
type CFExecutor struct {
	Client      cloudformationiface.CloudFormationAPI
	StackName   string
	TemplateURL string
	Parameters  []*cloudformation.Parameter
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor CFExecutor) CreateStack() error {
	//generate cloudformation CreateStackInput to be used to create stack
	input := &cloudformation.CreateStackInput{}

	input.SetTemplateURL(*aws.String(executor.TemplateURL))
	input.SetStackName(*aws.String(executor.StackName))
	input.SetParameters(executor.Parameters)
	input.SetCapabilities(createCapability())

	fmt.Println("about to create stack with parameters: ", input)
	_, err := executor.Client.CreateStack(input)
	//if there's an error return it
	if err != nil {
		fmt.Println("Got error creating stack:")
		fmt.Println(err.Error())
	}
	return err

}

//PauseUntilFinished is a method to wait on the status of a cloudformation stack until it finishes
func (executor CFExecutor) PauseUntilFinished() error {
	fmt.Println("Waiting for stack to be created")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(executor.StackName)}
	err := executor.Client.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be created")
		fmt.Println(err)
	}
	return err
}

func createCapability() []*string {
	capabilities := make([]*string, 0)
	capIAM := "CAPABILITY_IAM"
	capabilities = append(capabilities, &capIAM)

	return capabilities
}
