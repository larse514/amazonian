#!/bin/bash

set -e

bucketname=$1
stackname=$2

aws cloudformation package --template-file cloudformation/deploymentbucket.yml --output-template-file bucket-output.yml --s3-bucket circleci.deployables
aws cloudformation deploy --template-file bucket-output.yml --capabilities CAPABILITY_IAM --stack-name "${stackname}" --parameter-overrides S3BucketName="${bucketname}" 
