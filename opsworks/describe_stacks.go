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
  UseCustomCookbooks,
  UseOpsworksSecurityGroups,
  VpcId string
  Attributes,
  ChefConfiguration,
  ConfigurationManager,
  CustomCookbooksSource map[string]string
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
