package opsworks

import "github.com/bwilkins/aws"


var EndpointDefinition = aws.EndpointDefinition{
  SignatureVersion: 4,
  DataInterchangeFormat: "json",
  Host: "opsworks.us-east-1.amazonaws.com",
  Region: "us-east-1",
  ServiceName: "opsworks",
  Version: "2013-02-18",
  TargetPrefix: "OpsWorks_20130218",
}
