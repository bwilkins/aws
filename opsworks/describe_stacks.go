package opsworks

import "bytes"
import "encoding/json"
import "github.com/bwilkins/aws"

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

func DescribeStacks(request DescribeStacksRequest) (*DescribeStacksResponse, error) {
  bodyEncoded, err := json.Marshal(request)
  if err != nil {
    return nil, err
  }

  r, _ := aws.NewRequest("POST", "DescribeStacks", EndpointDefinition, bytes.NewReader(bodyEncoded))

  v := new(DescribeStacksResponse)
  return v, r.Do(v)
}
