# amazonian [![CircleCI](https://circleci.com/gh/larse514/amazonian.svg?style=svg)](https://circleci.com/gh/larse514/amazonian)
Opensource tool to faciliate AWS application resource development.  Amazonian will either let you provide your own cluster, or create a new one for you.  

## installation
Eventually the plan is to add this as a commandline tool and distribute it to various targets.  I am thinking at least targeting MacOS (Homebrew) and at least one Linux OS (maybe ubuntu?).
<br />
<br />
In the meantime there are two options to use amazonian.  <br />
1) Pull the binary from the temporary S3 distribution bucket here: [amazonian](https://s3.amazonaws.com/amazonian.package.release/latest/amazonian) <br />
2) A full install and build<br />

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
| Subnets        | List of VPC Subnets to deploy cluster to.                              | No       |           | Required if ClusterExists is false                      |
| KeyName        | Key name to use for EC2 instances within ECS cluster.                  | No       |           |                                                         |
| ClusterSize    | Number of host machines for cluster.                                   | No       | 1         |                                                         |
| MaxSize        | Max number of host machines cluster can scale to                       | No       | 1         |                                                         |
| InstanceType   | Type of EC2 machine                                                    | No       | t2.medium | Required if ClusterExists is false                      |

An example command of how you might run amazonian can be seen below:

`./amazonian --VPC=vpc-c7aa77be --Priority=12 --HostedZoneName=vssdevelopment.com --Image=willejs/go-hello-world --ServiceName=Node --ContainerName=Hello --ClusterName=amazonian-ecs-dev --ClusterExists=false --Subnets=subnet-b61d81fe,subnet-0202dc58 --KeyName=dummy_key1 ClusterSize=1 mazSizePrt=1 instanceTypePrt=t2.medium`

# contributing
If you would like to contribute to amazonian feel free to create a pull request or to fork the project itself. While amazonian is still under active development, and has not been released in any form, also feel free to raise issues as that will aide the development process itself.

## environment setup
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
