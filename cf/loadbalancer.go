package cf

import (
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

//LoadBalancer is a struct to interact with ELBV2 API
type LoadBalancer struct {
	Client elbv2iface.ELBV2API
}

//GetHighestPriority is a method to retrieve highest priority for ELBV2 Listener Rules
func (lb LoadBalancer) GetHighestPriority(listenerArn *string) (string, error) {
	input := &elbv2.DescribeRulesInput{ListenerArn: listenerArn}

	output, err := lb.Client.DescribeRules(input)
	if err != nil {
		println("Error retrieving Rules Input ", err.Error())
		return "", errors.New("Error retrieving Rules Input")
	}

	return getHighestPriorty(output), nil
}

//helper method to grab highest priority from Rules slice
func getHighestPriorty(output *elbv2.DescribeRulesOutput) string {
	priority := 1

	for _, rule := range output.Rules {
		i, err := strconv.Atoi(*rule.Priority)
		if err != nil {
			println("Error processing priority")
		}
		if i > priority {
			priority = i
		}
	}

	return strconv.Itoa(priority)
}
