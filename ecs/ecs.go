package ecs

import (
	"errors"
)

type ecs struct {
	stackName       string
	domainName      string
	keyName         string
	vpcID           string
	subnetID        string
	desiredCapacity int
	maxSize         int
	InstanceType    string
}

type ecsBuilder struct {
	ecs ecs
}

//constructor method for ecsBuilder
func newEcsBuilder(ecs ecs) *ecsBuilder {
	builder := ecsBuilder{ecs}
	return &builder
}

func (builder *ecsBuilder) IsValid() error {
	if builder.ecs.domainName == "" || builder.ecs.keyName == "" {
		return errors.New("Invalid ECS Cluster")
	}

	return nil
}
