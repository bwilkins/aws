package ec2

import "github.com/bwilkins/aws"


var EndpointDefinition = aws.EndpointDefinition{
	SignatureVersion: 2,
	DataInterchangeFormat: "xml",
	Host: "ec2.ap-southeast-1.amazonaws.com",
	Region: "ap-southeast-1",
	ServiceName: "ec2",
	Version: "2013-02-18",
	TargetPrefix: "OpsWorks_20130218",
}
