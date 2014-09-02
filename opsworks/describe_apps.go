package opsworks

import "bytes"
import "encoding/json"
import "github.com/bwilkins/aws/signing/v4"

type DescribeAppsResponse struct {
  Apps []App
}

type DescribeAppsRequest struct {
  StackId string `json:",omitempty"`
  AppIds []string `json:",omitempty"`
}

type App struct {
  AppId string
  AppSource AppSourceBlock
  Attributes AttributesBlock
  CreatedAt string
  DataSources []DataSourceBlock
  Description string
  Domains []string
  EnableSsl string
  Name string
  Shortname string
  StackId string
  Type string
  SslConfiguration SslConfigurationBlock
}

type AppSourceBlock struct {
  Password string
  Revision string
  SshKey string
  Type string
  Url string
  Username string
}

type AttributesBlock map[string]string

type DataSourceBlock struct {
  Arn string
  DatabaseName string
  Type string
}

type SslConfigurationBlock struct {
  Certificate string
  Chain string
  PrivateKey string
}

func DescribeApps(request DescribeAppsRequest) (*DescribeAppsResponse, error) {
  bodyEncoded, err := json.Marshal(request)
  if err != nil {
    return nil, err
  }

  r, _ := v4.NewRequest("POST", "DescribeApps", EndpointDefinition, bytes.NewReader(bodyEncoded))

  v := new(DescribeAppsResponse)
  return v, r.Do(v)
}
