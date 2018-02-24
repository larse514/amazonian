package commandlineargs

import "testing"

//Basic Test Testxxx
//Benchmark Benchmarkxxx
//Example Examplexxxx
func TestValidateArguments_fail(t *testing.T) {
	arg1 := "argument one"
	arg2 := ""
	err := ValidateArguments(arg1, arg2)

	if err == nil {
		t.Log("Error should not be nil")
		t.Fail()
	}
}

func TestValidateArguments_failNoArgs(t *testing.T) {

	err := ValidateArguments()

	if err == nil {
		t.Log("No args should result in error")
		t.Fail()
	}
}
func TestValidateArguments_pass(t *testing.T) {
	arg1 := "argument one"
	arg2 := "argument two"
	err := ValidateArguments(arg1, arg2)

	if err != nil {
		t.Log("Error should be nil if args are specified")
		t.Fail()
	}
}
