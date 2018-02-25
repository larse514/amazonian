package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"fmt"
)

//Cloudformation struct used to create cloudformation stacks
type cfExecutor struct {
	client      cloudformationiface.CloudFormationAPI
	stackName   string
	templateURL string
	parameters  []*cloudformation.Parameter
}

//ExecuteStack is a general method to create aws cloudformation stacks
func (executor *cfExecutor) CreateStack() error {
	//generate cloudformation CreateStackInput to be used to create stack
	input := &cloudformation.CreateStackInput{TemplateURL: aws.String(executor.templateURL), StackName: aws.String(executor.stackName), Parameters: executor.parameters}

	_, err := executor.client.CreateStack(input)
	//if there's an error return it
	if err != nil {
		fmt.Println("Got error creating stack:")
		fmt.Println(err.Error())
	}
	return err

}

//DescribeStack is a method to wait on the status of a cloudformation stack until it finishes
func (executor *cfExecutor) PauseUntilFinished() error {
	fmt.Println("Waiting for stack to be created")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(executor.stackName)}
	err := executor.client.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be created")
		fmt.Println(err)
	}
	return err
}
