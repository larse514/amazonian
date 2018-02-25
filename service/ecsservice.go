package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
)

//CreateService is a method that creates a service for an ecs service
func CreateService(parameterMap map[string]string, stackName string, templateURL string) error {

	// Set stack name, template url
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	//todo-maybe meove this out even further
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create CloudFormation client in region
	svc := cloudformation.New(sess)
	parameters := cf.CreateCloudformationParameters(parameterMap)

	executor := cf.CFExecutor{Client: svc, StackName: stackName, TemplateURL: templateURL, Parameters: parameters}
	//create the stack
	err := executor.CreateStack()
	if err != nil {
		println("Error processing create stack request ", err.Error)
		return err
	}
	//then wait
	err = executor.PauseUntilFinished()
	if err != nil {
		println("Error while attempting to wait for stack to finish processing ", err.Error)
		return err
	}
	return nil
}
