# amazonian [![CircleCI](https://circleci.com/gh/larse514/amazonian.svg?style=svg)](https://circleci.com/gh/larse514/amazonian)
Opensource tool to faciliate AWS application resource development

## installation
Eventually the plan is to add this as a commandline tool and distribute it to various targets.  I am thinking at least targeting MacOS (Homebrew) and at least one Linux OS (maybe ubuntu?).
<br />
<br />
In the meantime a full install and build is required to run


## environment setup
amazonian requires Golang and the aws go SDK installed

Step 1: _Install Go_ <br />
Follow the installation process in [Golang](https://golang.org/doc/install#install) docs <br /> 
Step 2: _Install Go SDK_ <br />
Follow the setup process defined in the [AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)
## build 
Step 1: _clone from git_ <br />
`$ git clone https://github.com/larse514/amazonian.git`  <br />
Step 2: _build amazonian_ <br />
`$ make` <br />

# running amazonian
### parameters
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
| Subnets        | List of VPC Subnets to deploy cluster to.                              | No       |           | Required if ClusterExists is false                      |
| KeyName        | Key name to use for EC2 instances within ECS cluster.                  | No       |           |                                                         |
| ClusterSize    | Number of host machines for cluster.                                   | No       | 1         |                                                         |
| MaxSize        | Max number of host machines cluster can scale to                       | No       | 1         |                                                         |
| InstanceType   | Type of EC2 machine                                                    | No       | t2.medium | Required if ClusterExists is false                      |
`./workdir/amazonian --VPC=vpc-c7aa77be --Priority=12 --HostedZoneName=vssdevelopment.com --Image=willejs/go-hello-world --ServiceName=Node --ContainerName=Hello --ClusterName=amazonian-ecs-dev --ClusterExists=false --Subnets=subnet-b61d81fe,subnet-0202dc58 --KeyName=dummy_key1 ClusterSize=1 mazSizePrt=1 instanceTypePrt=t2.medium`