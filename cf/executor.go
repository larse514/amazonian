package cf

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"fmt"
)

//Executor is an interface to execute and create stacks
type Executor interface {
	CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error
	PauseUntilFinished(stackName string) error
}

//CFExecutor struct used to create cloudformation stacks
type CFExecutor struct {
	Client cloudformationiface.CloudFormationAPI
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor CFExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	//generate cloudformation CreateStackInput to be used to create stack
	input := &cloudformation.CreateStackInput{}

	input.SetTemplateBody(*aws.String(templateBody))
	input.SetStackName(*aws.String(stackName))
	input.SetParameters(parameters)
	input.SetCapabilities(createCapability())
	//todo-refactor to return output
	_, err := executor.Client.CreateStack(input)
	//if there's an error return it
	if err != nil {
		fmt.Println("Got error creating stack: ", err.Error())
		return errors.New("Error creating stack")

	}
	return nil

}

//PauseUntilFinished is a method to wait on the status of a cloudformation stack until it finishes
func (executor CFExecutor) PauseUntilFinished(stackName string) error {
	fmt.Println("Waiting for stack to be created")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	err := executor.Client.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be created")
		fmt.Println(err)
	}
	return err
}

//helper method which statically generates CAPABILITY_IAM (a requirement for CloudFormation)
func createCapability() []*string {
	capabilities := make([]*string, 0)
	capIAM := "CAPABILITY_IAM"
	capabilities = append(capabilities, &capIAM)

	return capabilities
}
