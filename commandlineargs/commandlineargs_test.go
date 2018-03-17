package commandlineargs

import (
	"strings"
	"testing"
)

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

func TestCreateRandomStringRandomSeed(t *testing.T) {
	val1 := createRandomString("service")

	if val1 == "service698081" {
		t.Log("Val ", val1, " should not equal ")
		t.Fail()
	}
}

func TestCreateArgs(t *testing.T) {
	args := createArgs()
	if args.ServiceName == "" {
		t.Log("service is nil")
		t.Fail()
	}

	if args.ClusterName == "" || !strings.HasPrefix(args.ClusterName, cluster) {
		t.Log("cluster name is invalid ", args.ClusterName)
		t.Fail()
	}
	if args.ServiceName != args.ContainerName {
		t.Log("ContainerName", args.ContainerName, " not defaulted to ServiceName of ", args.ServiceName)
		t.Fail()
	}
}

func TestValidateECSParametersErrorTrue(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}

}
func TestValidateECSParametersECSClusterShouldNotExist(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBArn = "NOT EMPTY String"

	exists, _ := doesECSExist(args)
	if exists {
		t.Log("ECS Cluster should not exist")
		t.Fail()
	}

}
func TestValidateECSParametersALBListenerError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBArn = "NOT EMPTY String"
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}
}

func TestValidateECSParametersECSALBArnError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}

}
func TestValidateECSParametersECSHostedZoneIDError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}

}
func TestValidateECSParametersECSDNSNameError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}

}
func TestValidateECSParametersClusterArnError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err == nil {
		t.Log("No error found when one of expected")
		t.Fail()
	}

}
func TestValidateECSParametersNoError(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	_, err := doesECSExist(args)
	if err != nil {
		t.Log("Error returned when non-expected ", err.Error())
		t.Fail()
	}

}
func TestValidateECSParametersECSClusterExists(t *testing.T) {
	args := CommandLineArgs{}
	args.ECSALBFullName = "NOT EMPTY String"
	args.ECSALBListener = "NOT EMPTY String"
	args.ECSClusterARN = "NOT EMPTY String"
	args.ECSDNSName = "NOT EMPTY String"
	args.ECSHostedZoneID = "NOT EMPTY String"
	args.ECSALBArn = "NOT EMPTY String"

	exists, _ := doesECSExist(args)
	if !exists {
		t.Log("ECS Cluster shoudl exist")
		t.Fail()
	}

}
