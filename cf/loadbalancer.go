package cf

import (
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

//LoadBalancer is an interface to define services
type LoadBalancer interface {
	GetHighestPriority(listenerArn *string) (int, error)
}

//AWSElb is a struct to interact with ELBV2 API
type AWSElb struct {
	Client elbv2iface.ELBV2API
}

//GetHighestPriority is a method to retrieve highest priority for ELBV2 Listener Rules
func (lb AWSElb) GetHighestPriority(listenerArn *string) (int, error) {
	input := &elbv2.DescribeRulesInput{ListenerArn: listenerArn}

	output, err := lb.Client.DescribeRules(input)
	if err != nil {
		println("Error retrieving Rules Input ", err.Error())
		return 0, errors.New("Error retrieving Rules Input")
	}

	return getHighestPriorty(output), nil
}

//helper method to grab highest priority from Rules slice
func getHighestPriorty(output *elbv2.DescribeRulesOutput) int {
	priority := 1

	for _, rule := range output.Rules {

		i, err := strconv.Atoi(*rule.Priority)
		if err != nil {
			println("Error processing priority ", err.Error())
			continue
		}
		if i > priority {
			priority = i
		}
	}

	return priority
}
