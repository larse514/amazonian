#!/bin/bash

set -e
##Test creation of vpc, cluster, and container with defaults

image=$1
image2=$2

./workdir/amazonian --HostedZoneName=vssdevelopment.com --Image=${image} --PortMapping=8080

source amazonian-output

curl --fail https://${ServiceName}.vssdevelopment.com/

aws cloudformation delete-stack --stack-name "${ServiceName}"

##Now test and ensure it can reuse the same cluster and vpc
# echo | Test cluster reuse |

CONTAINER_NAME=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCName=${VPCName} --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${CONTAINER_NAME} --ClusterName=${ClusterName} --PortMapping=8080

# ./workdir/amazonian --VPCName=${VPCName} --HostedZoneName=vssdevelopment.com \
# --Image=${image} --ServiceName=${CONTAINER_NAME} --ContainerName=${CONTAINER_NAME} \
# --ClusterName=${ClusterName} --PortMapping=8080

curl --fail https://${CONTAINER_NAME}.vssdevelopment.com/
# echo | Cleaning up ${CONTAINER_NAME} and ${CLUSTER_NAME} |
CONTAINER_NAME2=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCName=${VPCName} --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${CONTAINER_NAME2} --ClusterName=${ClusterName} --PortMapping=8080

curl --fail https://${CONTAINER_NAME2}.vssdevelopment.com/

#now we want to test if we can deploy a new version

./workdir/amazonian --VPCName=${VPCName} --HostedZoneName=vssdevelopment.com \
--Image=${image2} --ServiceName=${CONTAINER_NAME2} --ClusterName=${ClusterName} --PortMapping=8080

curl --fail https://${CONTAINER_NAME2}.vssdevelopment.com/

source amazonian-output

echo "about to delete ${CONTAINER_NAME} ${CONTAINER_NAME2} ${ClusterName} and ${VPCName}"

aws cloudformation delete-stack --stack-name "${CONTAINER_NAME}"
aws cloudformation delete-stack --stack-name "${CONTAINER_NAME2}"
aws cloudformation delete-stack --stack-name "${ClusterName}"

#next, test if we can create a new cluster in the same vpc
NEW_CLUSTER=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCName=${VPCName} --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${CONTAINER_NAME2} --ClusterName=${NEW_CLUSTER} --PortMapping=8080


curl --fail https://${CONTAINER_NAME2}.vssdevelopment.com/

# now check if we can grab the existing cluster values and create something with that
#Get output values, this is a soemwhat naive approach since it is a lot of api calls
ecsclusterarn=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecscluster')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbarn=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbarn')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbdnsname=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbdnsname')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbhostedzoneid=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbhostedzoneid')]][0][*].{OutputValue:OutputValue}" --output text`
alblistener=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'alblistener')]][0][*].{OutputValue:OutputValue}" --output text`
ecslbfullname=`aws cloudformation describe-stacks --stack-name "${NEW_CLUSTER}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'ecslbfullname')]][0][*].{OutputValue:OutputValue}" --output text`
vpcId=`aws cloudformation describe-stacks --stack-name "${VPCName}" --query "Stacks[0].[Outputs[? starts_with(OutputKey, 'VPC')]][0][*].{OutputValue:OutputValue}" --output text`

echo "retrived ${vpcId}"

EXISTING_CLUSTER_CONTAINER=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCId=${vpcId} --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${EXISTING_CLUSTER_CONTAINER} --ClusterName=${NEW_CLUSTER} --PortMapping=8080 \
--ECSClusterARN=${ecsclusterarn} --ECSALBHostedZoneID=${ecslbhostedzoneid} --ECSALNDNSName=${ecslbdnsname} \
--ECSALBArn=${ecslbarn} --ECSALBListener=${alblistener} --ECSALBFullName=${ecslbfullname}

curl --fail https://${EXISTING_CLUSTER_CONTAINER}.vssdevelopment.com/

source amazonian-output
aws cloudformation delete-stack --stack-name "${CONTAINER_NAME2}"
aws cloudformation delete-stack --stack-name "${EXISTING_CLUSTER_CONTAINER}"
aws cloudformation delete-stack --stack-name "${NEW_CLUSTER}"
aws cloudformation delete-stack --stack-name "${VPCName}"

