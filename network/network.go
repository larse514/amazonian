package network

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/larse514/amazonian/assets"
	"github.com/larse514/amazonian/cf"
)

const (
	vpcTemplate    = "ias/cloudformation/vpc.yml"
	tenantParam    = "Tenant"
	cidrBlockParam = "CidrBlock"
	ws1CidrParam   = "WSSN1Cidr"
	ws2CidrParam   = "WSSN2Cidr"
	ws3CidrParam   = "WSSN3Cidr"
	app1CidrParam  = "APPSN1Cidr"
	app2CidrParam  = "APPSN2Cidr"
	app3CidrParam  = "APPSN3Cidr"
	db1CidrParam   = "DBSN1Cidr"
	db2CidrParam   = "DBSN2Cidr"

	//default cidr blocks
	cidrBlock = "192.168.0.0/16"
	ws1Cidr   = "192.168.4.0/24"
	ws2Cidr   = "192.168.5.0/24"
	ws3Cidr   = "192.168.6.0/24"

	app1Cidr = "192.168.14.0/24"
	app2Cidr = "192.168.15.0/24"
	app3Cidr = "192.168.16.0/24"

	db1Cidr = "192.168.34.0/24"
	db2Cidr = "192.168.35.0/24"
)

//Network is an interface to define operations with which to create
//cloud provider networks
type Network interface {
	CreateNetwork() error
}

//VPC is a struct representing AWS VPC object
type VPC struct {
	Executor   cf.Executor
	Name       string
	Tenant     string
	CIDRBlock  string
	WSSubnets  []Subnet
	APPSubnets []Subnet
	DBSubnets  []Subnet
}

//Subnet is a struct representing a subnet (not the best description i'll admit)
type Subnet struct {
	cidrBlock string
}

//CreateNetwork is a method to create a VPC network in AWS
func (vpc VPC) CreateNetwork() error {
	//Grab VPC template
	vpcTemplateBody, err := assets.GetAsset(vpcTemplate)
	if err != nil {
		fmt.Println("Error retrieving vpc template ", err.Error())
		return errors.New("Error creating template body for vpc")
	}
	err = vpc.Executor.CreateStack(vpcTemplateBody, vpc.Name, vpc.createVPCParameters())
	if err != nil {
		fmt.Println("Error creating vpc stack ", err.Error())
		return errors.New("Error creating vpc")
	}
	return vpc.Executor.PauseUntilFinished(vpc.Name)
}

//CreateClusterParameters will create the Parameter list to generate an ecs cluster
//todo- unit tests!!!
func (vpc VPC) createVPCParameters() []*cloudformation.Parameter {
	//we need to convert this (albiet awkwardly for the time being) to Cloudformation Parameters
	//we do as such first by converting everything to a key value map
	//key being the CF Param name, value is the value to provide to the cloudformation template
	//todo- refactor this approach
	parameterMap := make(map[string]string, 0)
	parameterMap[tenantParam] = vpc.Tenant
	parameterMap[cidrBlockParam] = vpc.CIDRBlock
	parameterMap[ws1CidrParam] = vpc.WSSubnets[0].cidrBlock
	parameterMap[ws2CidrParam] = vpc.WSSubnets[1].cidrBlock
	parameterMap[ws3CidrParam] = vpc.WSSubnets[2].cidrBlock
	parameterMap[app1CidrParam] = vpc.APPSubnets[0].cidrBlock
	parameterMap[app2CidrParam] = vpc.APPSubnets[1].cidrBlock
	parameterMap[app3CidrParam] = vpc.APPSubnets[2].cidrBlock
	parameterMap[db1CidrParam] = vpc.DBSubnets[0].cidrBlock
	parameterMap[db2CidrParam] = vpc.DBSubnets[1].cidrBlock
	return cf.CreateCloudformationParameters(parameterMap)

}

//CreateDefaultVPC is a method to create default VPC with Default subn
func CreateDefaultVPC(name string, tenant string) *VPC {
	vpc := VPC{}
	vpc.Name = name
	vpc.Tenant = tenant
	vpc.CIDRBlock = cidrBlock
	//app subnets
	appSN1 := Subnet{cidrBlock: app1Cidr}
	appSN2 := Subnet{cidrBlock: app2Cidr}
	appSN3 := Subnet{cidrBlock: app3Cidr}
	vpc.APPSubnets = []Subnet{appSN1, appSN2, appSN3}

	wsSn1 := Subnet{cidrBlock: ws1Cidr}
	wsSn2 := Subnet{cidrBlock: ws2Cidr}
	wsSn3 := Subnet{cidrBlock: ws3Cidr}
	vpc.WSSubnets = []Subnet{wsSn1, wsSn2, wsSn3}

	dbSn1 := Subnet{cidrBlock: db1Cidr}
	dbSn2 := Subnet{cidrBlock: db2Cidr}
	vpc.DBSubnets = []Subnet{dbSn1, dbSn2}

	return &vpc
}
