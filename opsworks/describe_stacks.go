package opsworks

import "github.com/bwilkins/aws/request"

type DescribeStacksRequest struct {
  StackIds []string `json:",omitempty"`
}

type DescribeStacksResponse struct {
  Stacks []Stack
}

type Stack struct {
  Arn,
  CreatedAt,
  CustomJson,
  DefaultAvailabilityZone,
  DefaultInstanceProfileArn,
  DefaultOs,
  DefaultRootDeviceType,
  DefaultSshKeyName,
  DefaultSubnetId,
  HostnameTheme,
  Name,
  Region,
  ServiceRoleArn,
  StackId,
  VpcId string
  Attributes map[string]string
  CustomCookbooksSource TCustomCookbooksSource
  ChefConfiguration TChefConfiguration
  UseCustomCookbooks,
  UseOpsworksSecurityGroups bool
}

type TChefConfiguration struct {
  BerkshelfVersion string
  ManageBerkshelf bool
}
type TConfigurationManager struct {
  Name, Version string
}

type TCustomCookbooksSource struct {
  Password,
  Revision,
  SshKey,
  Type,
  Url,
  Username string
}

func DescribeStacks(req DescribeStacksRequest) (*DescribeStacksResponse, error) {
  r, _ := request.NewRequest("POST", "DescribeStacks", EndpointDefinition, req)

  v := new(DescribeStacksResponse)
  return v, r.Do(v)
}
