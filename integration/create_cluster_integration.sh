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

./workdir/amazonian --VPCId=${VPCId} --VpcExists=true --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${CONTAINER_NAME} --ContainerName=${CONTAINER_NAME} \
--ClusterName=${ClusterName} --ClusterExists=true --PortMapping=8080

curl --fail https://${CONTAINER_NAME}.vssdevelopment.com/
# echo | Cleaning up ${CONTAINER_NAME} and ${CLUSTER_NAME} |
CONTAINER_NAME2=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCId=${VPCId} --VpcExists=true --HostedZoneName=vssdevelopment.com \
--Image=${image} --ServiceName=${CONTAINER_NAME2} --ContainerName=${CONTAINER_NAME2} \
--ClusterName=${ClusterName} --ClusterExists=true --PortMapping=8080

curl --fail https://${CONTAINER_NAME2}.vssdevelopment.com/

#now we want to test if we can deploy a new version

# ./workdir/amazonian --VPCId=${VPCId} --VpcExists=true --HostedZoneName=vssdevelopment.com \
# --Image=${image2} --ServiceName=${CONTAINER_NAME2} --ContainerName=${CONTAINER_NAME2} \
# --ClusterName=${ClusterName} --ClusterExists=true --PortMapping=8080

# curl --fail https://${CONTAINER_NAME2}.vssdevelopment.com/

source amazonian-output

echo "about to delete ${CONTAINER_NAME} ${CONTAINER_NAME2} ${ClusterName} and ${VPCName}"

aws cloudformation delete-stack --stack-name "${CONTAINER_NAME}"
aws cloudformation delete-stack --stack-name "${CONTAINER_NAME2}"
aws cloudformation delete-stack --stack-name "${ClusterName}"
aws cloudformation delete-stack --stack-name "${VPCName}"
