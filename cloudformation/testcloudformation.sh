#!/bin/bash

set -e

ecstemplate=$1
container=$2

aws cloudformation package --template-file cloudformation/ecs.yml --output-template-file ecs-output.yml --s3-bucket circleci.deployables
aws cloudformation deploy --template-file ecs-output.yml --capabilities CAPABILITY_IAM --stack-name "${ecstemplate}" --parameter-overrides KeyName=dummy_key1 VpcId=vpc-c7aa77be SubnetId=subnet-b61d81fe,subnet-0202dc58 InstanceType=t2.medium DomainName=vssdevelopment.com

#Get output values, this is a soemwhat naive approach since it is a lot of api calls
ecscluster=`aws cloudformation describe-stacks --stack-name "${ecstemplate}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecscluster')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbarn=`aws cloudformation describe-stacks --stack-name "${ecstemplate}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbarn')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbdnsname=`aws cloudformation describe-stacks --stack-name "${ecstemplate}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbdnsname')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbhostedzoneid=`aws cloudformation describe-stacks --stack-name "${ecstemplate}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbhostedzoneid')]][0][*].{OutputValue:OutputValue}" --output text`
alblistener=`aws cloudformation describe-stacks --stack-name "${ecstemplate}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'alblistener')]][0][*].{OutputValue:OutputValue}" --output text`

#package and deploy an instance
aws cloudformation package --template-file cloudformation/hellocontainer.yml --output-template-file hellocontainer-output.yml --s3-bucket circleci.deployables
aws cloudformation deploy --template-file hellocontainer-output.yml --capabilities CAPABILITY_IAM --stack-name "${container}" --parameter-overrides VpcId=vpc-c7aa77be Priority=1 ecscluster="${ecscluster}" ecslbarn="${ecslbarn}" ecslbdnsname="${ecslbdnsname}"	 ecslbhostedzoneid="${ecslbhostedzoneid}" alblistener="${alblistener}" HostedZoneName=vssdevelopment.com  ServiceName=hello ContainerName=hello
#just send basic test
sleep 15
curl --fail https://hello.vssdevelopment.com/