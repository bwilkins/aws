package opsworks

import "bytes"
import "encoding/json"
import "github.com/bwilkins/aws"

type DescribeStackSummaryRequest struct {
  StackId string `json:",omitempty"`
}

type DescribeStackSummaryResponse struct {
  StackSummary StackSummary
}

type StackSummary struct {
  AppsCount int64
  Arn string
  InstanceCount InstanceCountBlock
  LayersCount int64
  Name string
  StackId string
}

type InstanceCountBlock struct {
  Booting,
  ConnectionLost,
  Online,
  Pending,
  Rebooting,
  Requested,
  RunningSetup,
  SetupFailed,
  ShuttingDown,
  StartFailed,
  Stopped,
  Stopping,
  Terminated,
  Terminating int64
}

func (i *InstanceCountBlock) Sum() int64 {
  return i.Online
  //return i.Booting + i.ConnectionLost + i.Online +
    //i.Pending + i.Rebooting + i.Requested + i.RunningSetup +
    //i.SetupFailed + i.ShuttingDown + i.Stopped + i.Stopping +
    //i.Terminated + i.Terminating
}


func DescribeStackSummary(request DescribeStackSummaryRequest) (*DescribeStackSummaryResponse, error) {
  bodyEncoded, err := json.Marshal(request)
  if err != nil {
    return nil, err
  }

  r, _ := aws.NewRequest("POST", "DescribeStackSummary", EndpointDefinition, bytes.NewReader(bodyEncoded))

  v := new(DescribeStackSummaryResponse)
  return v, r.Do(v)
}
