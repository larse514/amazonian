package network

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/cf"
)

const (
	expectedCreateStackError = "Error creating vpc"
)

type mockGoodExecutor struct {
	cf.Executor
}
type mockBadStackExecutor struct {
	cf.Executor
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor mockGoodExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return nil
}
func (executor mockGoodExecutor) PauseUntilFinished(stackName string) error {
	return nil
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor mockBadStackExecutor) CreateStack(templateBody string, stackName string, parameters []*cloudformation.Parameter) error {
	return errors.New("ERROR")
}
func (executor mockBadStackExecutor) PauseUntilFinished(stackName string) error {
	return nil
}
func TestCreateNetwork(t *testing.T) {
}

func TestCreateVPCParameters(t *testing.T) {
	vpc := createVPC(mockGoodExecutor{})

	parms := vpc.createVPCParameters()
	if parms == nil {
		t.Fail()
	}
	len := len(parms)
	if len != 10 {
		t.Log("invalid num parms, expected 10, found ", len)
	}
}

func TestCreateVPC(t *testing.T) {
	vpc := createVPC(mockGoodExecutor{})
	err := vpc.CreateNetwork()
	if err != nil {
		t.Log("error returned ", err.Error())
		t.Fail()
	}

}
func TestCreateNetworkCreateStackErrors(t *testing.T) {
	vpc := createVPC(mockBadStackExecutor{})
	err := vpc.CreateNetwork()
	if err == nil {
		t.Log("error not return ")
		t.Fail()
	}
	if err.Error() != expectedCreateStackError {
		t.Log("correct error not return. expected ", expectedCreateStackError, " found ", err.Error())
		t.Fail()
	}

}
func TestCreateDefaultVPCAppSubnets(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	appC1 := vpc.APPSubnets[0].cidrBlock
	if appC1 != app1Cidr {
		t.Log("incorrect cidr block for app subnet expected ", app1Cidr, " got ", appC1)
		t.Fail()
	}
	appC2 := vpc.APPSubnets[1].cidrBlock
	if appC2 != app2Cidr {
		t.Log("incorrect cidr block for app subnet expected ", app2Cidr, " got ", appC2)
		t.Fail()

	}
	appC3 := vpc.APPSubnets[2].cidrBlock
	if appC3 != app3Cidr {
		t.Log("incorrect cidr block for app subnet expected ", app3Cidr, " got ", appC3)
	}
	//ws tests

}
func TestCreateDefaultVPCWSSubnets(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	wsC1 := vpc.WSSubnets[0].cidrBlock
	if wsC1 != ws1Cidr {
		t.Log("incorrect cidr block for ws subnet expected ", ws1Cidr, " got ", wsC1)
		t.Fail()
	}
	wsC2 := vpc.WSSubnets[1].cidrBlock
	if wsC2 != ws2Cidr {
		t.Log("incorrect cidr block for ws subnet expected ", ws2Cidr, " got ", wsC2)
		t.Fail()

	}
	wsC3 := vpc.WSSubnets[2].cidrBlock
	if wsC3 != ws3Cidr {
		t.Log("incorrect cidr block for ws subnet expected ", ws3Cidr, " got ", wsC3)
	}

}
func TestCreateDefaultVPCDBSubnets(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	dbC1 := vpc.DBSubnets[0].cidrBlock
	if dbC1 != db1Cidr {
		t.Log("incorrect cidr block for db subnet expected ", db1Cidr, " got ", dbC1)
		t.Fail()
	}
	dbC2 := vpc.DBSubnets[1].cidrBlock
	if dbC2 != db2Cidr {
		t.Log("incorrect cidr block for db subnet expected ", db2Cidr, " got ", dbC2)
		t.Fail()

	}

}
func TestCreateDefaultVPCCidr(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	cidr := vpc.CIDRBlock
	if cidr != cidrBlock {
		t.Log("incorrect cidr block for  expected ", cidrBlock, " got ", cidr)
		t.Fail()
	}

}
func TestCreateDefaultVPCName(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	name := vpc.Name
	if name != "NAME" {
		t.Log("incorrect name  expected NAME but got ", name)
		t.Fail()
	}

}
func TestCreateDefaultVPCTenat(t *testing.T) {
	vpc := CreateDefaultVPC("NAME", "TENANT")
	tenant := vpc.Tenant
	if tenant != "TENANT" {
		t.Log("incorrect tenant  expected TENANT but got ", tenant)
		t.Fail()
	}

}

//helper method
func createVPC(executor cf.Executor) VPC {
	vpc := VPC{Executor: executor}
	vpc.Tenant = "TENANT"
	vpc.CIDRBlock = "10.0.0.0/16"

	vpc.WSSubnets = []Subnet{}
	vpc.WSSubnets = append(vpc.WSSubnets, Subnet{cidrBlock: "10.0.0.0/32"})
	vpc.WSSubnets = append(vpc.WSSubnets, Subnet{cidrBlock: "10.0.0.1/32"})
	vpc.WSSubnets = append(vpc.WSSubnets, Subnet{cidrBlock: "10.0.0.2/32"})
	vpc.APPSubnets = []Subnet{}
	vpc.APPSubnets = append(vpc.APPSubnets, Subnet{cidrBlock: "10.0.0.0/0"})
	vpc.APPSubnets = append(vpc.APPSubnets, Subnet{cidrBlock: "10.0.0.0/0"})
	vpc.APPSubnets = append(vpc.APPSubnets, Subnet{cidrBlock: "10.0.0.0/0"})
	vpc.DBSubnets = []Subnet{}
	vpc.DBSubnets = append(vpc.DBSubnets, Subnet{cidrBlock: "10.0.0.0/0"})
	vpc.DBSubnets = append(vpc.DBSubnets, Subnet{cidrBlock: "10.0.0.0/0"})
	return vpc
}
