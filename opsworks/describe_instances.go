package opsworks

import "bytes"
import "encoding/json"
import "github.com/bwilkins/aws"

type DescribeInstancesResponse struct {
  Instances []Instance
}

type DescribeInstancesRequest struct {
  StackId string `json:",omitempty"`
  LayerId string `json:",omitempty"`
  InstanceIds []string `json:",omitempty"`
}

type Instance struct {
  Architecture,
  AutoScalingType,
  AvailabilityZone,
  CreatedAt,
  EbsOptimized,
  Ec2InstanceId,
  Hostname,
  InstallUpdatesOnBoot,
  InstanceId,
  InstanceProfileArn,
  InstanceType,
  Os,
  PrivateDns,
  PrivateIp,
  PublicDns,
  PublicIp,
  RootDeviceType,
  RootDeviceVolumeId,
  SshHostDsaKeyFingerprint,
  SshHostRsaKeyFingerprint,
  StackId,
  Status,
  VirtualizationType string
  LayerIds,
  SecurityGroupIds []string
}

func DescribeInstances(request DescribeInstancesRequest) (*DescribeInstancesResponse, error) {
  bodyEncoded, err := json.Marshal(request)
  if err != nil {
    return nil, err
  }

  r, _ := aws.NewRequest("POST", EndpointDefinition, bytes.NewReader(bodyEncoded))

  v := new(DescribeInstancesResponse)
  return v, r.Do(v)
}
