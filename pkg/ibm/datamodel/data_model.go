package datamodel

import (
	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	vpcv1 "github.com/IBM/vpc-go-sdk/vpcv1"
)

// The following types define the "canonical data model" for IBM resources.
// For the most part, these are the SDK types extended with extra information like tags or info from multiple calls

type TaggedResource interface {
	SetTags([]string)
	GetCRN() *string
}

// BaseTaggedResource type is used as an abstraction for all resources that IBM allows tagging
type BaseTaggedResource struct {
	Tags []string `json:"tags"`
}

func (res *BaseTaggedResource) SetTags(tags []string) {
	res.Tags = tags
}

// VPC configuration object
type VPC struct {
	vpcv1.VPC
	BaseTaggedResource
}

func NewVPC(sdkVPC *vpcv1.VPC) *VPC {
	return &VPC{VPC: *sdkVPC}
}

func (res *VPC) GetCRN() *string { return res.VPC.CRN }

// Subnet configuration object
type Subnet struct {
	vpcv1.Subnet
	ReservedIps []vpcv1.ReservedIP `json:"reserved_ips"`
	BaseTaggedResource
}

func NewSubnet(subnet *vpcv1.Subnet) *Subnet {
	return &Subnet{Subnet: *subnet}
}

func (res *Subnet) GetCRN() *string { return res.CRN }

// PublicGateway configuration object
type PublicGateway struct {
	vpcv1.PublicGateway
	BaseTaggedResource
}

func NewPublicGateway(publicGateway *vpcv1.PublicGateway) *PublicGateway {
	return &PublicGateway{PublicGateway: *publicGateway}
}

func (res *PublicGateway) GetCRN() *string { return res.CRN }

// FloatingIP configuration object
type FloatingIP struct {
	vpcv1.FloatingIP
	BaseTaggedResource
}

func NewFloatingIP(floatingIP *vpcv1.FloatingIP) *FloatingIP {
	return &FloatingIP{FloatingIP: *floatingIP}
}

func (res *FloatingIP) GetCRN() *string { return res.CRN }

// NetworkACL configuration object
type NetworkACL struct {
	vpcv1.NetworkACL
	BaseTaggedResource
}

func NewNetworkACL(networkACL *vpcv1.NetworkACL) *NetworkACL {
	return &NetworkACL{NetworkACL: *networkACL}
}

func (res *NetworkACL) GetCRN() *string { return res.CRN }

// SecurityGroup configuration object
type SecurityGroup struct {
	vpcv1.SecurityGroup
	BaseTaggedResource
}

func NewSecurityGroup(securityGroup *vpcv1.SecurityGroup) *SecurityGroup {
	return &SecurityGroup{SecurityGroup: *securityGroup}
}

func (res *SecurityGroup) GetCRN() *string { return res.CRN }

// EndpointGateway configuration object
type EndpointGateway struct {
	vpcv1.EndpointGateway
	BaseTaggedResource
}

func NewEndpointGateway(endpointGateway *vpcv1.EndpointGateway) *EndpointGateway {
	return &EndpointGateway{EndpointGateway: *endpointGateway}
}

func (res *EndpointGateway) GetCRN() *string { return res.CRN }

// Instance configuration object
type Instance struct {
	vpcv1.Instance
	NetworkInterfaces []vpcv1.NetworkInterface `json:"network_interfaces"`
	BaseTaggedResource
}

func NewInstance(instance *vpcv1.Instance) *Instance {
	return &Instance{Instance: *instance}
}

func (res *Instance) GetCRN() *string { return res.CRN }

// RoutingTable configuration object (not taggable)
type RoutingTable struct {
	vpcv1.RoutingTable
	Routes []vpcv1.Route `json:"routes"`
}

func NewRoutingTable(rt *vpcv1.RoutingTable, routes []vpcv1.Route) *RoutingTable {
	return &RoutingTable{RoutingTable: *rt, Routes: routes}
}

// LoadBalancer configuration objects

// LoadBalancerPool object with explicit members (not references)
type LoadBalancerPool struct {
	vpcv1.LoadBalancerPool
	Members []vpcv1.LoadBalancerPoolMember `json:"members"`
}

func NewLoadBalancerPool(loadBalancerPool *vpcv1.LoadBalancerPool,
	members []vpcv1.LoadBalancerPoolMember) LoadBalancerPool {
	return LoadBalancerPool{LoadBalancerPool: *loadBalancerPool, Members: members}
}

// LoadBalancerListenerPolicy configuration with explicit rules (not references)
type LoadBalancerListenerPolicy struct {
	vpcv1.LoadBalancerListenerPolicy
	Rules []vpcv1.LoadBalancerListenerPolicyRule `json:"rules"`
}

func NewLoadBalancerListenerPolicy(
	loadBalancerListenerPolicy *vpcv1.LoadBalancerListenerPolicy,
	rules []vpcv1.LoadBalancerListenerPolicyRule) LoadBalancerListenerPolicy {
	return LoadBalancerListenerPolicy{LoadBalancerListenerPolicy: *loadBalancerListenerPolicy, Rules: rules}
}

// LoadBalancerListener configuration object with explicit policies (not references)
type LoadBalancerListener struct {
	vpcv1.LoadBalancerListener
	Policies []LoadBalancerListenerPolicy `json:"policies"`
}

func NewLoadBalancerListener(
	loadBalancerListener *vpcv1.LoadBalancerListener, policies []LoadBalancerListenerPolicy) LoadBalancerListener {
	return LoadBalancerListener{LoadBalancerListener: *loadBalancerListener, Policies: policies}
}

// LoadBalancer configuration object with explicit listeners and pools (not references)
type LoadBalancer struct {
	vpcv1.LoadBalancer
	Listeners []LoadBalancerListener `json:"listeners"`
	Pools     []LoadBalancerPool     `json:"pools"`
	BaseTaggedResource
}

func NewLoadBalancer(
	lb *vpcv1.LoadBalancer,
	listeners []LoadBalancerListener,
	pools []LoadBalancerPool) *LoadBalancer {
	return &LoadBalancer{
		LoadBalancer: *lb,
		Listeners:    listeners,
		Pools:        pools,
	}
}

func (res *LoadBalancer) GetCRN() *string { return res.CRN }

type IKSWorkerNode struct {
	iksv1.GetWorkerResponse
}

func NewIKSWorkerNode(getWorkerResponse *iksv1.GetWorkerResponse) *IKSWorkerNode {
	return &IKSWorkerNode{GetWorkerResponse: *getWorkerResponse}
}

type IKSWorkerPool struct {
	iksv1.GetWorkerPoolResponse
}

func NewIKSWorkerPool(getWorkerPoolResponse *iksv1.GetWorkerPoolResponse) *IKSWorkerPool {
	return &IKSWorkerPool{GetWorkerPoolResponse: *getWorkerPoolResponse}
}
