package commandlineargs

import "testing"

//Basic Test Testxxx
//Benchmark Benchmarkxxx
//Example Examplexxxx
func TestValidateArguments_fail(t *testing.T) {
	arg1 := "argument one"
	arg2 := ""
	err := validateArguments(arg1, arg2)

	if err == nil {
		t.Log("Error should not be nil")
		t.Fail()
	}
}

func TestValidateArguments_failNoArgs(t *testing.T) {

	err := validateArguments()

	if err == nil {
		t.Log("No args should result in error")
		t.Fail()
	}
}
func TestValidateArguments_pass(t *testing.T) {
	arg1 := "argument one"
	arg2 := "argument two"
	err := validateArguments(arg1, arg2)

	if err != nil {
		t.Log("Error should be nil if args are specified")
		t.Fail()
	}
}

func TestCreateRandomString(t *testing.T) {
	val1 := createRandomString("STARTER")
	val2 := createRandomString("STARTER")

	if val1 == val2 {
		t.Log("Val ", val1, " should not equal ", val2)
		t.Fail()
	}
}
