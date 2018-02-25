package ecs

import "testing"

func TestEcsIsValid(t *testing.T) {
	builder := new(ecsBuilder)

	err := builder.IsValid()

	if err == nil {
		t.Log("ECS is invalid, error should not be nill")
		t.Fail()
	}
}
