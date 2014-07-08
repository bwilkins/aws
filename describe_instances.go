package aws

import "bytes"
import "encoding/json"

// Sugar
type DescribeInstancesResponse struct {
  Instances []Instance
}

type DescribeInstancesRequest struct {
  StackId string `json:",omitempty"`
  LayerId string `json:",omitempty"`
  InstanceIds []string `json:",omitempty"`
}

type Instance struct {
  AmiId string
  Status  string
  Hostname    string
  IpAddress  string
}

func DescribeInstances(request DescribeInstancesRequest) (*DescribeInstancesResponse, error) {
  bodyEncoded, err := json.Marshal(request)
  if err != nil {
    return nil, err
  }

  r, _ := NewRequest("POST", "https://"+OpsWorksHost+"/", bytes.NewReader(bodyEncoded))

  v := new(DescribeInstancesResponse)
  return v, Do(r, v)
}
