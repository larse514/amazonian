package cf

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

//GetStack tests, mocks, and methods

// Define a mock to return a basic success
type mockGoodelbv2ifaceClient struct {
	elbv2iface.ELBV2API
}

func (m mockGoodelbv2ifaceClient) DescribeRules(*elbv2.DescribeRulesInput) (*elbv2.DescribeRulesOutput, error) {
	return createDescribeRules(), nil
}

type mockBadelbv2ifaceClient struct {
	elbv2iface.ELBV2API
}

func (m mockBadelbv2ifaceClient) DescribeRules(*elbv2.DescribeRulesInput) (*elbv2.DescribeRulesOutput, error) {
	return nil, errors.New("ERROR")
}
func createDescribeRules() *elbv2.DescribeRulesOutput {
	output := &elbv2.DescribeRulesOutput{}

	rules := make([]*elbv2.Rule, 0)
	priority1 := "1"
	priority4 := "4"
	priority7 := "7"

	priority2 := "2"

	rules = append(rules, &elbv2.Rule{Priority: &priority1})
	rules = append(rules, &elbv2.Rule{Priority: &priority2})
	rules = append(rules, &elbv2.Rule{Priority: &priority4})
	rules = append(rules, &elbv2.Rule{Priority: &priority7})

	output.SetRules(rules)
	return output
}

func TestGetHighestPriorty(t *testing.T) {
	lb := AWSElb{Client: mockGoodelbv2ifaceClient{}}

	listenerARN := "listenerarn"

	priority, _ := lb.GetHighestPriority(&listenerARN)

	if priority != 7 {
		t.Log("priority ", priority, " invalid")
		t.Fail()
	}

}
func TestGetHighestPriortyError(t *testing.T) {
	lb := AWSElb{Client: mockBadelbv2ifaceClient{}}

	listenerARN := "listenerarn"

	_, err := lb.GetHighestPriority(&listenerARN)

	if err == nil {
		t.Log("error not returned")
		t.Fail()
	}

}

func TestGetHighestPriortyLowerCase(t *testing.T) {

	priority := getHighestPriorty(createDescribeRules())

	if priority != 7 {
		t.Log("priority ", priority, " invalid")
		t.Fail()
	}
}
