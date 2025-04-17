/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package datamodel

import (
	"encoding/json"

	tgw "github.com/IBM/networking-go-sdk/transitgatewayapisv1"

	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// Helper function for unmarshalling

func jsonToMap(jsonStr []byte) (map[string]json.RawMessage, error) {
	var result map[string]json.RawMessage
	err := json.Unmarshal(jsonStr, &result)
	return result, err
}

func basicUnmarshal[A any](data []byte, unmarshalFunc func(map[string]json.RawMessage, any) error,
	objRef *A, tags *BaseTaggedResource) error {
	asMap, err := jsonToMap(data)
	if err != nil {
		return err
	}

	asRef := new(A)
	err = unmarshalFunc(asMap, &asRef)
	if err != nil {
		return err
	}
	*objRef = *asRef

	if tags != nil {
		err = json.Unmarshal(data, tags)
		if err != nil {
			return err
		}
	}

	return nil
}

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
	Region          string                `json:"region"`
	AddressPrefixes []vpcv1.AddressPrefix `json:"address_prefixes"`
	BaseTaggedResource
}

func NewVPC(sdkVPC *vpcv1.VPC, region string, addressPrefixes []vpcv1.AddressPrefix) *VPC {
	return &VPC{VPC: *sdkVPC, Region: region, AddressPrefixes: addressPrefixes}
}

func (res *VPC) UnmarshalJSON(data []byte) error {
	asMap, err := jsonToMap(data)
	if err != nil {
		return err
	}
	asObj := &vpcv1.VPC{}
	err = vpcv1.UnmarshalVPC(asMap, &asObj)
	if err != nil {
		return err
	}
	res.VPC = *asObj

	val, ok := asMap["address_prefixes"]
	if ok {
		if err := json.Unmarshal(val, &res.AddressPrefixes); err != nil {
			return err
		}
	}

	val, ok = asMap["region"]
	if ok {
		if err := json.Unmarshal(val, &res.Region); err != nil {
			return err
		}
	}

	return json.Unmarshal(data, &res.BaseTaggedResource)
}

func (res *VPC) GetCRN() *string { return res.CRN }

// ReservedIPWrapper is an alias to vpcv1.ReservedIP that allows us to override the implementation of UnmarshalJSON
type ReservedIPWrapper struct {
	vpcv1.ReservedIP
}

func (res *ReservedIPWrapper) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalReservedIP, &res.ReservedIP, nil)
}

// Subnet configuration object
type Subnet struct {
	vpcv1.Subnet
	ReservedIps []ReservedIPWrapper `json:"reserved_ips"`
	BaseTaggedResource
}

func NewSubnet(subnet *vpcv1.Subnet, reservedIPs []vpcv1.ReservedIP) *Subnet {
	reservedIPWraps := make([]ReservedIPWrapper, len(reservedIPs))
	for i := range reservedIPs {
		reservedIPWraps[i] = ReservedIPWrapper{reservedIPs[i]}
	}
	return &Subnet{Subnet: *subnet, ReservedIps: reservedIPWraps}
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

func (res *FloatingIP) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalFloatingIP, &res.FloatingIP, &res.BaseTaggedResource)
}

// NetworkACL configuration object
type NetworkACL struct {
	vpcv1.NetworkACL
	BaseTaggedResource
}

func NewNetworkACL(networkACL *vpcv1.NetworkACL) *NetworkACL {
	return &NetworkACL{NetworkACL: *networkACL}
}

func (res *NetworkACL) GetCRN() *string { return res.CRN }

func (res *NetworkACL) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalNetworkACL, &res.NetworkACL, &res.BaseTaggedResource)
}

// SecurityGroup configuration object
type SecurityGroup struct {
	vpcv1.SecurityGroup
	BaseTaggedResource
}

func NewSecurityGroup(securityGroup *vpcv1.SecurityGroup) *SecurityGroup {
	return &SecurityGroup{SecurityGroup: *securityGroup}
}

func (res *SecurityGroup) GetCRN() *string { return res.CRN }

func (res *SecurityGroup) UnmarshalJSON(data []byte) error {
	return basicUnmarshal[vpcv1.SecurityGroup](data, vpcv1.UnmarshalSecurityGroup, &res.SecurityGroup, &res.BaseTaggedResource)
}

// EndpointGateway configuration object
type EndpointGateway struct {
	vpcv1.EndpointGateway
	BaseTaggedResource
}

func NewEndpointGateway(endpointGateway *vpcv1.EndpointGateway) *EndpointGateway {
	return &EndpointGateway{EndpointGateway: *endpointGateway}
}

func (res *EndpointGateway) GetCRN() *string { return res.CRN }

func (res *EndpointGateway) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalEndpointGateway, &res.EndpointGateway, &res.BaseTaggedResource)
}

// Instance configuration object
type Instance struct {
	vpcv1.Instance
	NetworkInterfaces []vpcv1.NetworkInterface `json:"network_interfaces"`
	BaseTaggedResource
}

func NewInstance(instance *vpcv1.Instance, networkInterfaces []vpcv1.NetworkInterface) *Instance {
	return &Instance{Instance: *instance, NetworkInterfaces: networkInterfaces}
}

func (res *Instance) GetCRN() *string { return res.CRN }

// Virtual Network Interface object
type VirtualNI struct {
	vpcv1.VirtualNetworkInterface
	BaseTaggedResource
}

func NewVirtualNI(vni *vpcv1.VirtualNetworkInterface) *VirtualNI {
	return &VirtualNI{VirtualNetworkInterface: *vni}
}

func (res *VirtualNI) GetCRN() *string { return res.CRN }

func (res *VirtualNI) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalVirtualNetworkInterface, &res.VirtualNetworkInterface, &res.BaseTaggedResource)
}

// RoutingTable configuration object (not taggable)
type RoutingTable struct {
	vpcv1.RoutingTable
	Routes []RouteWrapper      `json:"routes"`
	VPC    *vpcv1.VPCReference `json:"vpc"`
}

func NewRoutingTable(rt *vpcv1.RoutingTable, routes []vpcv1.Route, vpcRef *vpcv1.VPCReference) *RoutingTable {
	routesWrapper := make([]RouteWrapper, len(routes))
	for i := range routes {
		routesWrapper[i] = RouteWrapper{routes[i]}
	}
	return &RoutingTable{RoutingTable: *rt, Routes: routesWrapper, VPC: vpcRef}
}

// RouteWrapper is an alias to vpcv1.Route that allows us to override
// the implementation of UnmarshalJSON
type RouteWrapper struct {
	vpcv1.Route
}

func (res *RouteWrapper) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalRoute, &res.Route, nil)
}

// LoadBalancer configuration objects

// LoadBalancerPoolMemberWrapper is an alias to vpcv1.LoadBalancerPoolMember that allows us to override
// the implementation of UnmarshalJSON
type LoadBalancerPoolMemberWrapper struct {
	vpcv1.LoadBalancerPoolMember
}

func (res *LoadBalancerPoolMemberWrapper) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalLoadBalancerPoolMember, &res.LoadBalancerPoolMember, nil)
}

// LoadBalancerPool object with explicit members (not references)
type LoadBalancerPool struct {
	vpcv1.LoadBalancerPool
	Members []LoadBalancerPoolMemberWrapper `json:"members"`
}

func NewLoadBalancerPool(loadBalancerPool *vpcv1.LoadBalancerPool,
	members []vpcv1.LoadBalancerPoolMember) LoadBalancerPool {
	LoadBalancerPoolMemberWraps := make([]LoadBalancerPoolMemberWrapper, len(members))
	for i := range members {
		LoadBalancerPoolMemberWraps[i] = LoadBalancerPoolMemberWrapper{members[i]}
	}
	return LoadBalancerPool{LoadBalancerPool: *loadBalancerPool, Members: LoadBalancerPoolMemberWraps}
}

// LoadBalancerListenerPolicyRuleWrapper is an alias to vpcv1.LoadBalancerListenerPolicyRule that allows us to override
// the implementation of UnmarshalJSON
type LoadBalancerListenerPolicyRuleWrapper struct {
	vpcv1.LoadBalancerListenerPolicyRule
}

func (res *LoadBalancerListenerPolicyRuleWrapper) UnmarshalJSON(data []byte) error {
	return basicUnmarshal(data, vpcv1.UnmarshalLoadBalancerListenerPolicyRule, &res.LoadBalancerListenerPolicyRule, nil)
}

// LoadBalancerListenerPolicy configuration with explicit rules (not references)
type LoadBalancerListenerPolicy struct {
	vpcv1.LoadBalancerListenerPolicy
	Rules []LoadBalancerListenerPolicyRuleWrapper `json:"rules"`
}

func NewLoadBalancerListenerPolicy(
	loadBalancerListenerPolicy *vpcv1.LoadBalancerListenerPolicy,
	rules []vpcv1.LoadBalancerListenerPolicyRule) LoadBalancerListenerPolicy {
	rulesWrap := make([]LoadBalancerListenerPolicyRuleWrapper, len(rules))
	for i := range rules {
		rulesWrap[i] = LoadBalancerListenerPolicyRuleWrapper{rules[i]}
	}

	return LoadBalancerListenerPolicy{LoadBalancerListenerPolicy: *loadBalancerListenerPolicy, Rules: rulesWrap}
}

func (res *LoadBalancerListenerPolicy) UnmarshalJSON(data []byte) error {
	asMap, err := jsonToMap(data)
	if err != nil {
		return err
	}
	asObj := &vpcv1.LoadBalancerListenerPolicy{}
	err = vpcv1.UnmarshalLoadBalancerListenerPolicy(asMap, &asObj)
	if err != nil {
		return err
	}
	res.LoadBalancerListenerPolicy = *asObj

	var rules []LoadBalancerListenerPolicyRuleWrapper
	err = json.Unmarshal(asMap["rules"], &rules)
	if err != nil {
		return err
	}
	res.Rules = rules

	return nil
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

// TransitConnection configuration object
type TransitConnection struct {
	tgw.TransitConnection
}

func NewTransitConnection(transitConnection *tgw.TransitConnection) *TransitConnection {
	return &TransitConnection{TransitConnection: *transitConnection}
}

type TransitGateway struct {
	tgw.TransitGateway
}

func NewTransitGateway(transitGateway *tgw.TransitGateway) *TransitGateway {
	return &TransitGateway{TransitGateway: *transitGateway}
}

// IKSWorkerNode configuration object
type IKSCluster struct {
	iksv1.GetClustersResponse
	WorkerNodes []iksv1.GetWorkerResponse
}

func NewCluster(cluster *iksv1.GetClustersResponse, getWorkerResponse []iksv1.GetWorkerResponse) *IKSCluster {
	return &IKSCluster{GetClustersResponse: *cluster, WorkerNodes: getWorkerResponse}
}
