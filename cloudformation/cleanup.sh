#!/bin/bash
ecstemplate=$1
container=$2

aws cloudformation delete-stack --stack-name "${container}"

aws cloudformation delete-stack --stack-name "${ecstemplate}"
