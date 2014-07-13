package aws

type EndpointDefinition struct {
  SignatureAlgorithm,
  Host,
  Region,
  ServiceName,
  Version,
  TargetPrefix string
}

type Credentials struct {
  AccessId,
  SecretKey string
}

var AccessCredentials Credentials

func SetAccessCredentials(creds Credentials) {
  AccessCredentials = creds
}
