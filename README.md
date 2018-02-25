# amazonian
Opensource tool to faciliate AWS application resource development

## setup
Follow the setup process defined in the [AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)
amazonian assumes the use of environment variables. 

### Example command
`./amazonian --VPC=vpc-c7aa77be --Priority=1 --HostedZoneName=vssdevelopment.com --ELBHostedZoneId="Z35SXDOTRQ7X7K" --ELBDNSName="ECSLB-524689478.us-east-1.elb.amazonaws.com" --ELBARN="arn:aws:elasticloadbalancing:us-east-1:417615409974:loadbalancer/app/ECSLB/bdbdba6216f30332" --ClusterARN="amazonian-ecs-ECSCluster-1CO3I8UEEG80H" --ALBListenerARN="arn:aws:elasticloadbalancing:us-east-1:417615409974:listener/app/ECSLB/bdbdba6216f30332/8ca78dfa2878578f" --image=willejs/go-hello-world`


// // Create CloudFormation client in region
// svc := cloudformation.New(sess)
// Set stack name, template url
// Initialize a session that the SDK uses to load
// credentials from the shared credentials file ~/.aws/credentials
// and configuration from the shared configuration file ~/.aws/config.
// sess := session.Must(session.NewSessionWithOptions(session.Options{
// 	SharedConfigState: session.SharedConfigEnable,
// }))
