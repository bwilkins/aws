package opsworks

import "github.com/bwilkins/aws/signing/v4"

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
  Ec2InstanceId,
  Hostname,
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
  EbsOptimized,
  InstallUpdatesOnBoot bool
}

func DescribeInstances(request DescribeInstancesRequest) (*DescribeInstancesResponse, error) {
  r, _ := v4.NewRequest("POST", "DescribeInstances", EndpointDefinition, request)
  v := new(DescribeInstancesResponse)
  return v, r.Do(v)
}
