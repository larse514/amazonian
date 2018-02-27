# amazonian [![CircleCI](https://circleci.com/gh/larse514/amazonian.svg?style=svg)](https://circleci.com/gh/larse514/amazonian)
Opensource tool to faciliate AWS container based application development.  

## what does amazonian do?
![ECS](docs/ecs.png) <br />
The goal of amazonian is to abstract away the complexity of deploying continers in AWS.  You can use your own cluster, or let amazonian create one for you.  Either way, just provide a few configuration values, and amazonian will deploy, run, monitor, and secure your containers for you.  Its your infrastructure so there's no cost, other than the infrastructure required to run a cluster.

## installation
Eventually the plan is to add this as a commandline tool and distribute is to various targets.  I am thinking at least targeting MacOS (Homebrew) and at least one Linux OS (maybe ubuntu?).
<br />
<br />
In the meantime there are two options to use amazonian.  <br />
1) Pull the binary from the temporary S3 distribution bucket here: [amazonian](https://s3.amazonaws.com/amazonian.package.release/latest/amazonian) <br />
2) A full install and build <br />

# running amazonian

once the amazonian binary is ready to run, it can be executed from the command line.  There is currently no UI associated with amazonian as the main target it to facilitate CI/CD pipelines

### parameters
The following describes the parameters amazonian uses.

| Paramater      | Description                                                            | Required | Default   | Note                                                    |
|----------------|------------------------------------------------------------------------|----------|-----------|---------------------------------------------------------|
| VPC            | Target VPC to deploy your containers                                   | Yes      | None      |                                                         |
| Priority       | Priority to use for Load Balancer Listener rules                       | Yes      | None      | This will be removed soon for a dynamic lookup          |
| HostedZoneName | Route 53 hosted zone name to use for cluster and container deployments | Yes      | None      |                                                         |
| Image          | Docker Repository Image to be deployed as a container                  | Yes      | None      |                                                         |
| ServiceName    | Name of container service to be deployed                               | Yes      | None      |                                                         |
| ContainerName  | Name of container to be deployed                                       | Yes      | None      |                                                         |
| ClusterName    | Name ECS Cluster to use                                                | Yes      | None      | This will be expanded to include Fargate and Kubernetes |
| ClusterExists  | Specify whether to use an existing cluster                             | No       | false     |                                                         |
| Subnets        | List of VPC Subnets to deploy cluster to.                              | Sometimes       |           | Required if ClusterExists is false                      |
| KeyName        | Key name to use for EC2 instances within ECS cluster.                  | No       |           |                                                         |
| ClusterSize    | Number of host machines for cluster.                                   | No       | 1         |                                                         |
| MaxSize        | Max number of host machines cluster can scale to                       | No       | 1         |                                                         |
| InstanceType   | Type of EC2 machine                                                    | No       | t2.medium | Required if ClusterExists is false                      |

An example command of how you might run amazonian can be seen below:

`./amazonian --VPC=vpc-c7aa77be --Priority=12 --HostedZoneName=vssdevelopment.com --Image=willejs/go-hello-world --ServiceName=Node --ContainerName=Hello --ClusterName=amazonian-ecs-dev --ClusterExists=false --Subnets=subnet-b61d81fe,subnet-0202dc58 --KeyName=dummy_key1 ClusterSize=1 mazSizePrt=1 instanceTypePrt=t2.medium`

## environment setup
amazonian leverages the AWS SDKs in order to build the necessary infrastructure to support your containers.  This requires the environment with which amazonian is executed to be setup with appropriate IAM credentials and AWS configuration.  AWS provides documentation here:
[AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) <br /> 

amazonian itself requires the following minimum permissions to execute:
TODO- the below is not fully correct!!! it is currently a place holder.  You have been warned!

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "application-autoscaling:DeleteScalingPolicy",
                "application-autoscaling:DeregisterScalableTarget",
                "application-autoscaling:DescribeScalableTargets",
                "application-autoscaling:DescribeScalingActivities",
                "application-autoscaling:DescribeScalingPolicies",
                "application-autoscaling:PutScalingPolicy",
                "application-autoscaling:RegisterScalableTarget",
                "autoscaling:UpdateAutoScalingGroup",
                "autoscaling:CreateAutoScalingGroup",
                "autoscaling:CreateLaunchConfiguration",
                "autoscaling:DeleteAutoScalingGroup",
                "autoscaling:DeleteLaunchConfiguration",
                "autoscaling:Describe*",
                "cloudformation:CreateStack",
                "cloudformation:DeleteStack",
                "cloudformation:DescribeStack*",
                "cloudformation:UpdateStack",
                "cloudwatch:DescribeAlarms",
                "cloudwatch:GetMetricStatistics",
                "cloudwatch:PutMetricAlarm",
                "ec2:AssociateRouteTable",
                "ec2:AttachInternetGateway",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:CancelSpotFleetRequests",
                "ec2:CreateInternetGateway",
                "ec2:CreateRoute",
                "ec2:CreateRouteTable",
                "ec2:CreateSecurityGroup",
                "ec2:CreateSubnet",
                "ec2:CreateVpc",
                "ec2:DeleteSubnet",
                "ec2:DeleteVpc",
                "ec2:Describe*",
                "ec2:DetachInternetGateway",
                "ec2:DisassociateRouteTable",
                "ec2:ModifySubnetAttribute",
                "ec2:ModifyVpcAttribute",
                "ec2:RequestSpotFleet",
                "elasticloadbalancing:CreateListener",
                "elasticloadbalancing:CreateLoadBalancer",
                "elasticloadbalancing:CreateRule",
                "elasticloadbalancing:CreateTargetGroup",
                "elasticloadbalancing:DeleteListener",
                "elasticloadbalancing:DeleteLoadBalancer",
                "elasticloadbalancing:DeleteRule",
                "elasticloadbalancing:DeleteTargetGroup",
                "elasticloadbalancing:DescribeListeners",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancing:DescribeRules",
                "elasticloadbalancing:DescribeTargetGroups",
                "ecs:*",
                "events:DescribeRule",
                "events:DeleteRule",
                "events:ListRuleNamesByTarget",
                "events:ListTargetsByRule",
                "events:PutRule",
                "events:PutTargets",
                "events:RemoveTargets",
                "iam:ListAttachedRolePolicies",
                "iam:ListInstanceProfiles",
                "iam:ListRoles",
                "logs:CreateLogGroup",
                "logs:DescribeLogGroups",
                "logs:FilterLogEvents"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DeleteInternetGateway",
                "ec2:DeleteRoute",
                "ec2:DeleteRouteTable",
                "ec2:DeleteSecurityGroup"
            ],
            "Resource": [
                "*"
            ],
            "Condition": {
                "StringLike": {
                    "ec2:ResourceTag/aws:cloudformation:stack-name": "EC2ContainerService-*"
                }
            }
        },
        {
            "Action": "iam:PassRole",
            "Effect": "Allow",
            "Resource": [
                "*"
            ],
            "Condition": {
                "StringLike": {
                    "iam:PassedToService": "ecs-tasks.amazonaws.com"
                }
            }
        },
        {
            "Action": "iam:PassRole",
            "Effect": "Allow",
            "Resource": [
                "arn:aws:iam::*:role/ecsInstanceRole*"
            ],
            "Condition": {
                "StringLike": {
                    "iam:PassedToService": ["ec2.amazonaws.com", "ec2.amazonaws.com.cn"]
                }
            }
        },
        {
            "Action": "iam:PassRole",
            "Effect": "Allow",
            "Resource": [
                "arn:aws:iam::*:role/ecsAutoscaleRole*"
            ],
            "Condition": {
                "StringLike": {
                    "iam:PassedToService": ["application-autoscaling.amazonaws.com", "application-autoscaling.amazonaws.com.cn"]
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": "iam:CreateServiceLinkedRole",
            "Resource": "*",
            "Condition": {
                "StringLike": {
                    "iam:AWSServiceName": ["ecs.amazonaws.com", "spot.amazonaws.com", "spotfleet.amazonaws.com"]
                }
            }
        }
    ]
}
```

# contributing
If you would like to contribute to amazonian feel free to create a pull request or to fork the project itself. While amazonian is still under active development, and has not been released in any form, also feel free to raise issues as that will aide the development process.

## development environment setup
amazonian requires Golang and the aws go SDK installed

Step 1: _Install Go_ <br />
Follow the installation process in [Golang](https://golang.org/doc/install#install) docs <br /> 
Step 2: _Install Go SDK_ <br />
Follow the setup process defined in the [AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) <br />
Step 3: _IAM Role Setup_ <br />
TODO- add required role permissions<br />

## build 
Step 1: _clone from git_ <br />
`$ git clone https://github.com/larse514/amazonian.git`  <br />
Step 2: _build amazonian_ <br />
`$ make` <br />
Step 3: test amazonian_ <br />
`$ make test` <br />
or if you don't want to see the verbose logs 
`$ make test-min` <br />
