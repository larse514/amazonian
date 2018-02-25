#!/bin/bash

set -e
##Test creation of container from existing cluster
CONTAINER_NAME=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`
PRIORITY=`cat /dev/urandom | env LC_CTYPE=C tr -cd '0-9' | head -c 2`
./workdir/amazonian --VPC=vpc-c7aa77be --Priority=${PRIORITY} --HostedZoneName=vssdevelopment.com \
--Image=willejs/go-hello-world --StackName=${CONTAINER_NAME} --ServiceName=${CONTAINER_NAME} --ContainerName=${CONTAINER_NAME} \
--ClusterName=amazonian-ecs --ClusterExists

curl --fail https://${CONTAINER_NAME}.vssdevelopment.com/



aws cloudformation delete-stack --stack-name "${CONTAINER_NAME}"
