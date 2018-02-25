package service

type service interface {
	createService(serviceDefinition service) (string, error)
}

type serviceDefinition struct {
	serviceName       string
	containerName     string
	vpcID             string
	priority          int
	hostedZoneName    string
	ecslbHostedZoneID string
	ecslbDNSName      string
	ecslbArn          string
	ecsCluster        string
	albListener       string
	image             string
}
