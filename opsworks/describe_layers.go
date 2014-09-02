package opsworks

import "github.com/bwilkins/aws/request"

type DescribeLayersRequest struct {
  LayerIds []string `json:",omitempty"`
  StackId string `json:",omitempty"`
}

type DescribeLayersResponse struct {
  Layers []Layer
}

type Layer struct {
  Attributes map[string]string
  AutoAssignElasticIps bool
  AutoAssignPublicIps bool
  CreatedAt string
  CustomInstanceProfileArn string
  CustomRecipes RecipesBlock
  CustomSecurityGroupIds []string
  DefaultRecipes RecipesBlock
  DefaultSecurityGroupNames []string
  EnableAutoHealing bool
  InstallUpdatesOnBoot bool
  LayerId string
  Name string
  Packages []string
  Shortname string
  StackId string
  Type string
  UseEbsOptimizedInstances bool
  VolumeConfigurations []VolumeConfigBlock
}

type RecipesBlock struct {
  Configure,
  Deploy,
  Setup,
  Shutdown,
  Undeploy []string
}

type VolumeConfigBlock struct {
  Iops int
  MountPoint string
  NumberOfDisks int
  RaidLevel int
  Size int
  VolumeType string
}

func DescribeLayers(req DescribeLayersRequest) (*DescribeLayersResponse, error) {
  r, _ := request.NewRequest("POST", "DescribeLayers", EndpointDefinition, req)

  v := new(DescribeLayersResponse)
  return v, r.Do(v)
}
