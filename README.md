# amazonian [![CircleCI](https://circleci.com/gh/larse514/amazonian.svg?style=svg)](https://circleci.com/gh/larse514/amazonian)
Opensource tool to faciliate AWS application resource development

## environment setup
Step 1: _Install Go_ <br />
Follow the installation process in [Golang](https://golang.org/doc/install#install) docs <br /> 
Step 2: _Install Go SDK_ <br />
Follow the setup process defined in the [AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

## installation
Eventually the plan is to add this as a commandline tool and distribute it to various targets.  I am thinking at least targeting MacOS (Homebrew) and at least one Linux OS (maybe ubuntu?)

Step 1: _clone from git_ <br />
`$ git clone https://github.com/larse514/amazonian.git`  <br />
Step 2: build amazonian_ <br />
`$ make` <br />

### Example command
`./workdir/amazonian --VPC=vpc-c7aa77be --Priority=12 --HostedZoneName=vssdevelopment.com --Image=willejs/go-hello-world --ServiceName=Node --ContainerName=Hello --ClusterName=amazonian-ecs-dev --ClusterExists=false --Subnets=subnet-b61d81fe,subnet-0202dc58 --KeyName=dummy_key1 ClusterSize=1 mazSizePrt=1 instanceTypePrt=t2.medium`