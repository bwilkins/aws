package aws

type EndpointDefinition struct {
  DataInterchangeFormat,
  Host,
  Region,
  ServiceName,
  Version,
  TargetPrefix string
  SignatureVersion int
}

type Credentials struct {
  AccessId,
  SecretKey string
}

var AccessCredentials Credentials

func SetAccessCredentials(creds Credentials) {
  AccessCredentials = creds
}
