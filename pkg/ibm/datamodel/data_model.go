package datamodel

import "github.com/IBM/vpc-go-sdk/vpcv1"

// The following types define the "canonical data model" for IBM resources.
// For the most part, these are the SDK types extended with extra information like tags or info from multiple calls

// VPC configuration object
type VPC struct {
	vpcv1.VPC
	Tags []string `json:"tags"`
}

func NewVPC(sdkVPC *vpcv1.VPC) *VPC {
	return &VPC{VPC: *sdkVPC}
}

// Subnet configuration object
type Subnet struct {
	vpcv1.Subnet
	ReservedIps []vpcv1.ReservedIP `json:"reserved_ips"`
	Tags        []string           `json:"tags"`
}

func NewSubnet(subnet *vpcv1.Subnet) *Subnet {
	return &Subnet{Subnet: *subnet}
}

// PublicGateway configuration object
type PublicGateway struct {
	vpcv1.PublicGateway
	Tags []string `json:"tags"`
}

func NewPublicGateway(publicGateway *vpcv1.PublicGateway) *PublicGateway {
	return &PublicGateway{PublicGateway: *publicGateway}
}

// FloatingIP configuration object
type FloatingIP struct {
	vpcv1.FloatingIP
	Tags []string `json:"tags"`
}

func NewFloatingIP(floatingIP *vpcv1.FloatingIP) *FloatingIP {
	return &FloatingIP{FloatingIP: *floatingIP}
}

// NetworkACL configuration object
type NetworkACL struct {
	vpcv1.NetworkACL
	Tags []string `json:"tags"`
}

func NewNetworkACL(networkACL *vpcv1.NetworkACL) *NetworkACL {
	return &NetworkACL{NetworkACL: *networkACL}
}

// SecurityGroup configuration object
type SecurityGroup struct {
	vpcv1.SecurityGroup
	Tags []string `json:"tags"`
}

func NewSecurityGroup(securityGroup *vpcv1.SecurityGroup) *SecurityGroup {
	return &SecurityGroup{SecurityGroup: *securityGroup}
}

// EndpointGateway configuration object
type EndpointGateway struct {
	vpcv1.EndpointGateway
	Tags []string `json:"tags"`
}

func NewEndpointGateway(endpointGateway *vpcv1.EndpointGateway) *EndpointGateway {
	return &EndpointGateway{EndpointGateway: *endpointGateway}
}

// Instance configuration object
type Instance struct {
	vpcv1.Instance
	NetworkInterfaces []vpcv1.NetworkInterface `json:"network_interfaces"`
	Tags              []string                 `json:"tags"`
}

func NewInstance(instance *vpcv1.Instance) *Instance {
	return &Instance{Instance: *instance}
}
