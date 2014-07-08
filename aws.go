package aws

type EndpointDefinition struct {
  SignatureAlgorithm,
  Host,
  Region,
  ServiceName,
  Version,
  TargetPrefix string
}
