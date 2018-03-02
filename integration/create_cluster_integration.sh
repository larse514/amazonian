#!/bin/bash

set -e
##Test creation of vpc, cluster, and container with defaults

./workdir/amazonian --HostedZoneName=vssdevelopment.com --Image=willejs/go-hello-world

source amazonian-output

curl --fail https://${ServiceName}.vssdevelopment.com/

aws cloudformation delete-stack --stack-name "${ServiceName}"

##Now test and ensure it can reuse the same cluster and vpc
# echo | Test cluster reuse |

CONTAINER_NAME=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f' | head -c 5`

./workdir/amazonian --VPCId=${VPCId} --VpcExists=true --HostedZoneName=vssdevelopment.com \
--Image=willejs/go-hello-world --ServiceName=${CONTAINER_NAME} --ContainerName=${CONTAINER_NAME} \
--ClusterName=${ClusterName} --ClusterExists=true

curl --fail https://${CONTAINER_NAME}.vssdevelopment.com/
# echo | Cleaning up ${CONTAINER_NAME} and ${CLUSTER_NAME} |

source amazonian-output

echo "about to delete ${CONTAINER_NAME} ${ClusterName} and ${VPCName}"
aws cloudformation delete-stack --stack-name "${CONTAINER_NAME}"
aws cloudformation delete-stack --stack-name "${ClusterName}"
aws cloudformation delete-stack --stack-name "${VPCName}"
