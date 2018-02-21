#!/bin/bash
ecstemplate=$1

aws cloudformation delete-stack --stack-name "${ecstemplate}"
