/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

import (
	"fmt"
	"reflect"

	tgw "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

const pageSize = 50

type HasGetNextStart[T any] interface {
	GetNextStart() (*string, error)
}

func iteratePagedAPI[T any, Q HasGetNextStart[T]](
	apiFunc func(int64, *string) (Q, any, error),
	getArray func(Q) []T) ([]T, error) {
	var start *string = nil
	// We can use collection.TotalCount for efficiency, but usually it's a single page
	res := make([]T, 0)
	for {
		collection, _, err := apiFunc(pageSize, start)
		if err != nil {
			return nil, fmt.Errorf("[iteratePagedAPI] error getting item: %w", err)
		}
		res = append(res, getArray(collection)...)
		start, err = collection.GetNextStart()
		if err != nil {
			return nil, fmt.Errorf("[iteratePagedAPI] error getting next page: %w", err)
		}
		if start == nil {
			break
		}
	}
	return res, nil
}

// getResources is a generic function for collecting resources, where the conversion to internal resource is straightforward.
// R is the type of the IBM-Cloud-API resource to collect, C is the corresponding collection, I is the internal type
func getResources[R any, C HasGetNextStart[R], I any](
	listFunc func(pageSize int64, next *string) (C, any, error),
	getArrayFunc func(collection C) []R,
	convertFunc func(*R) *I,
) ([]*I, error) {
	resources, err := iteratePagedAPI(listFunc, getArrayFunc)
	if err != nil {
		var zero [0]R
		return nil, fmt.Errorf("error listing resources of type %s: %w", reflect.TypeOf(zero).Elem(), err)
	}
	res := make([]*I, len(resources))
	for i := range resources {
		res[i] = convertFunc(&resources[i])
	}

	return res, nil
}

func getVPCs(vpcService *vpcv1.VpcV1, region, resourceGroupID string) ([]*datamodel.VPC, error) {
	APIFunc := func(pageSize int64, next *string) (*vpcv1.VPCCollection, any, error) {
		return vpcService.ListVpcs(&vpcv1.ListVpcsOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	getArray := func(collection *vpcv1.VPCCollection) []vpcv1.VPC {
		return collection.Vpcs
	}
	vpcs, err := iteratePagedAPI(APIFunc, getArray)
	if err != nil {
		return nil, fmt.Errorf("[getVPCs] error getting VPCs: %w", err)
	}
	res := make([]*datamodel.VPC, len(vpcs))

	getArrayPrefixes := func(collection *vpcv1.AddressPrefixCollection) []vpcv1.AddressPrefix {
		return collection.AddressPrefixes
	}
	for i := range vpcs {
		APIFuncPrefixes := func(pageSize int64, next *string) (*vpcv1.AddressPrefixCollection, any, error) {
			return vpcService.ListVPCAddressPrefixes(&vpcv1.ListVPCAddressPrefixesOptions{Limit: &pageSize, Start: next, VPCID: vpcs[i].ID})
		}
		addressPrefixes, err := iteratePagedAPI(APIFuncPrefixes, getArrayPrefixes)
		if err != nil {
			return nil, fmt.Errorf("[getVPCs] error getting Address Prefixes: %w", err)
		}

		res[i] = datamodel.NewVPC(&vpcs[i], region, addressPrefixes)
	}
	return res, nil
}

func getSubnets(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.Subnet, error) {
	subnetAPIFunc := func(pageSize int64, next *string) (*vpcv1.SubnetCollection, any, error) {
		return vpcService.ListSubnets(&vpcv1.ListSubnetsOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	subnetGetArray := func(collection *vpcv1.SubnetCollection) []vpcv1.Subnet {
		return collection.Subnets
	}
	subnets, err := iteratePagedAPI(subnetAPIFunc, subnetGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getSubnets] error getting Subnets: %w", err)
	}
	res := make([]*datamodel.Subnet, len(subnets))
	for i := range subnets {
		reservedIPs, err := getReservedIps(vpcService, *subnets[i].ID, *subnets[i].Name)
		if err != nil {
			return nil, err
		}
		res[i] = datamodel.NewSubnet(&subnets[i], reservedIPs)
	}

	return res, nil
}

// getReservedIps is a second API call to get the list of reserved IPs in a subnet
func getReservedIps(vpcService *vpcv1.VpcV1, subnetID, name string) ([]vpcv1.ReservedIP, error) {
	reservedIPAPIFunc := func(pageSize int64, next *string) (*vpcv1.ReservedIPCollection, any, error) {
		options := vpcService.NewListSubnetReservedIpsOptions(subnetID)
		options.Limit = &pageSize
		options.Start = next
		return vpcService.ListSubnetReservedIps(options)
	}
	reservedIPGetArray := func(collection *vpcv1.ReservedIPCollection) []vpcv1.ReservedIP {
		return collection.ReservedIps
	}
	reservedIPs, err := iteratePagedAPI(reservedIPAPIFunc, reservedIPGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getReservedIps]: error getting reserved IPs for %s", name)
	}
	return reservedIPs, nil
}

func getPublicGateways(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.PublicGateway, error) {
	gatewayAPIFunc := func(pageSize int64, next *string) (*vpcv1.PublicGatewayCollection, any, error) {
		return vpcService.ListPublicGateways(&vpcv1.ListPublicGatewaysOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	gatewayGetArray := func(collection *vpcv1.PublicGatewayCollection) []vpcv1.PublicGateway {
		return collection.PublicGateways
	}

	return getResources(gatewayAPIFunc, gatewayGetArray, datamodel.NewPublicGateway)
}

func getFloatingIPs(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.FloatingIP, error) {
	floatingIPAPIFunc := func(pageSize int64, next *string) (*vpcv1.FloatingIPCollection, any, error) {
		return vpcService.ListFloatingIps(&vpcv1.ListFloatingIpsOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	floatingIPGetArray := func(collection *vpcv1.FloatingIPCollection) []vpcv1.FloatingIP {
		return collection.FloatingIps
	}

	return getResources(floatingIPAPIFunc, floatingIPGetArray, datamodel.NewFloatingIP)
}

func getNetworkACLs(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.NetworkACL, error) {
	networkACLAPIFunc := func(pageSize int64, next *string) (*vpcv1.NetworkACLCollection, any, error) {
		return vpcService.ListNetworkAcls(&vpcv1.ListNetworkAclsOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	networkACLGetArray := func(collection *vpcv1.NetworkACLCollection) []vpcv1.NetworkACL {
		return collection.NetworkAcls
	}

	return getResources(networkACLAPIFunc, networkACLGetArray, datamodel.NewNetworkACL)
}

func getSecurityGroups(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.SecurityGroup, error) {
	securityGroupAPIFunc := func(pageSize int64, next *string) (*vpcv1.SecurityGroupCollection, any, error) {
		return vpcService.ListSecurityGroups(&vpcv1.ListSecurityGroupsOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	securityGroupGetArray := func(collection *vpcv1.SecurityGroupCollection) []vpcv1.SecurityGroup {
		return collection.SecurityGroups
	}

	return getResources(securityGroupAPIFunc, securityGroupGetArray, datamodel.NewSecurityGroup)
}

// Get all Endpoint Gateways (VPEs)
func getEndpointGateways(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.EndpointGateway, error) {
	endpointGatewayAPIFunc := func(pageSize int64, next *string) (*vpcv1.EndpointGatewayCollection, any, error) {
		return vpcService.ListEndpointGateways(&vpcv1.ListEndpointGatewaysOptions{Limit: &pageSize, Start: next,
			ResourceGroupID: &resourceGroupID})
	}
	endpointGatewayGetArray := func(collection *vpcv1.EndpointGatewayCollection) []vpcv1.EndpointGateway {
		return collection.EndpointGateways
	}

	return getResources(endpointGatewayAPIFunc, endpointGatewayGetArray, datamodel.NewEndpointGateway)
}

func getInstances(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.Instance, error) {
	instanceAPIFunc := func(pageSize int64, next *string) (*vpcv1.InstanceCollection, any, error) {
		return vpcService.ListInstances(&vpcv1.ListInstancesOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID})
	}
	instanceGetArray := func(collection *vpcv1.InstanceCollection) []vpcv1.Instance {
		return collection.Instances
	}
	instances, err := iteratePagedAPI(instanceAPIFunc, instanceGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getInstances] error getting Instances: %w", err)
	}
	res := make([]*datamodel.Instance, len(instances))
	for i := range instances {
		id := *instances[i].ID
		name := *instances[i].Name

		networkInterfaces, err := getNetworkInterface(vpcService, id, name)
		if err != nil {
			return nil, err
		}
		res[i] = datamodel.NewInstance(&instances[i], networkInterfaces)
	}

	return res, nil
}

// Second API call to get detailed network interfaces information
func getNetworkInterface(vpcService *vpcv1.VpcV1, id, name string) ([]vpcv1.NetworkInterface, error) {
	options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
	options.SetInstanceID(id)
	networkInterfaces, _, err := vpcService.ListInstanceNetworkInterfaces(options)
	if err != nil {
		return nil, fmt.Errorf("[getInstances] error getting NW Interfaces for %s", name)
	}
	return networkInterfaces.NetworkInterfaces, nil
}

func getVirtualNIs(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.VirtualNI, error) {
	vniAPIFunc := func(pageSize int64, next *string) (*vpcv1.VirtualNetworkInterfaceCollection, any, error) {
		opts := vpcv1.ListVirtualNetworkInterfacesOptions{Limit: &pageSize, Start: next, ResourceGroupID: &resourceGroupID}
		return vpcService.ListVirtualNetworkInterfaces(&opts)
	}
	vniGetArray := func(collection *vpcv1.VirtualNetworkInterfaceCollection) []vpcv1.VirtualNetworkInterface {
		return collection.VirtualNetworkInterfaces
	}

	return getResources(vniAPIFunc, vniGetArray, datamodel.NewVirtualNI)
}

func getRoutingTables(vpcService *vpcv1.VpcV1, vpcList []*datamodel.VPC) ([]*datamodel.RoutingTable, error) {
	var res []*datamodel.RoutingTable

	routingTableAPIFunc := func(vpcID string) func(pageSize int64, next *string) (*vpcv1.RoutingTableCollection, any, error) {
		return func(pageSize int64, next *string) (*vpcv1.RoutingTableCollection, any, error) {
			options := vpcService.NewListVPCRoutingTablesOptions(vpcID)
			options.Limit = &pageSize
			options.Start = next
			return vpcService.ListVPCRoutingTables(options)
		}
	}
	routingTableGetArray := func(collection *vpcv1.RoutingTableCollection) []vpcv1.RoutingTable {
		return collection.RoutingTables
	}

	vpcResourceType := vpcv1.VPCReferenceResourceTypeVPCConst
	for i := range vpcList {
		vpc := vpcList[i]
		vpcID := *vpc.ID
		vpcRef := &vpcv1.VPCReference{
			CRN:          vpc.CRN,
			Href:         vpc.Href,
			ID:           vpc.ID,
			Name:         vpc.Name,
			ResourceType: &vpcResourceType,
		}

		routingTables, err := iteratePagedAPI(routingTableAPIFunc(vpcID), routingTableGetArray)
		if err != nil {
			return nil, fmt.Errorf("[getRoutingTables] error getting Routing Tables for %s: %w", vpcID, err)
		}
		for j := range routingTables {
			routes, err := getRoutes(vpcService, vpcID, *routingTables[j].ID)
			if err != nil {
				return nil, err
			}
			res = append(res, datamodel.NewRoutingTable(&routingTables[j], routes, vpcRef))
		}
	}

	return res, nil
}

func getRoutes(vpcService *vpcv1.VpcV1, vpcID, rtID string) ([]vpcv1.Route, error) {
	routingTableRouteAPIFunc := func(pageSize int64, next *string) (*vpcv1.RouteCollection, any, error) {
		options := vpcService.NewListVPCRoutingTableRoutesOptions(vpcID, rtID)
		options.Limit = &pageSize
		options.Start = next
		return vpcService.ListVPCRoutingTableRoutes(options)
	}
	routingTableRouteGetArray := func(collection *vpcv1.RouteCollection) []vpcv1.Route {
		return collection.Routes
	}
	routes, err := iteratePagedAPI(routingTableRouteAPIFunc, routingTableRouteGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getRoutes] error getting routes for %s: %w", rtID, err)
	}
	return routes, nil
}

// Get all Load Balancers
func getLoadBalancers(vpcService *vpcv1.VpcV1, resourceGroupID string) ([]*datamodel.LoadBalancer, error) {
	loadBalancerAPIFunc := func(pageSize int64, next *string) (*vpcv1.LoadBalancerCollection, any, error) {
		return vpcService.ListLoadBalancers(&vpcv1.ListLoadBalancersOptions{Limit: &pageSize, Start: next})
	}
	loadBalancerGetArray := func(collection *vpcv1.LoadBalancerCollection) []vpcv1.LoadBalancer {
		return collection.LoadBalancers
	}
	loadBalancers, err := iteratePagedAPI(loadBalancerAPIFunc, loadBalancerGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getLoadBalancers] error getting Load Balancer: %w", err)
	}
	res := make([]*datamodel.LoadBalancer, 0, len(loadBalancers))
	for i := range loadBalancers {
		if resourceGroupID != "" && *(loadBalancers[i].ResourceGroup.ID) != resourceGroupID {
			continue
		}

		// get all the listeners
		lbID := *loadBalancers[i].ID
		listenerOptions := &vpcv1.ListLoadBalancerListenersOptions{}
		listenerOptions.SetLoadBalancerID(lbID)
		listenersCollection, _, err := vpcService.ListLoadBalancerListeners(listenerOptions)
		if err != nil {
			return nil, fmt.Errorf("[getLoadBalancers] error getting listeners for %s: %w", lbID, err)
		}

		listeners := make([]datamodel.LoadBalancerListener, len(listenersCollection.Listeners))
		for j := range listenersCollection.Listeners {
			listenerID := *listenersCollection.Listeners[j].ID
			policiesOptions := &vpcv1.ListLoadBalancerListenerPoliciesOptions{}
			policiesOptions.SetLoadBalancerID(lbID)
			policiesOptions.SetListenerID(listenerID)
			policiesCollection, _, polErr := vpcService.ListLoadBalancerListenerPolicies(policiesOptions)
			if polErr != nil {
				return nil, fmt.Errorf("[getLoadBalancers] error getting policies for %s: %w", listenerID, err)
			}

			policies := make([]datamodel.LoadBalancerListenerPolicy, len(policiesCollection.Policies))
			for k := range policiesCollection.Policies {
				policies[k], polErr = getPolicyRules(vpcService, lbID, listenerID, &policiesCollection.Policies[k])
				if polErr != nil {
					return nil, polErr
				}
			}
			listeners[j] = datamodel.NewLoadBalancerListener(&listenersCollection.Listeners[j], policies)
		}

		// get all the pools
		poolOptions := &vpcv1.ListLoadBalancerPoolsOptions{}
		poolOptions.SetLoadBalancerID(lbID)
		poolsCollection, _, err := vpcService.ListLoadBalancerPools(poolOptions)
		if err != nil {
			return nil, fmt.Errorf("[getLoadBalancers] error getting pools for %s: %w", lbID, err)
		}
		pools := make([]datamodel.LoadBalancerPool, len(poolsCollection.Pools))
		for j := range pools {
			pools[j], err = getPoolMembers(vpcService, lbID, &poolsCollection.Pools[j])
			if err != nil {
				return nil, err
			}
		}
		res = append(res, datamodel.NewLoadBalancer(&loadBalancers[i], listeners, pools))
	}

	return res, nil
}

func getPoolMembers(vpcService *vpcv1.VpcV1, lbID string, vpcPool *vpcv1.LoadBalancerPool) (datamodel.LoadBalancerPool, error) {
	options := &vpcv1.ListLoadBalancerPoolMembersOptions{}
	options.SetLoadBalancerID(lbID)
	options.SetPoolID(*vpcPool.ID)
	members, _, err := vpcService.ListLoadBalancerPoolMembers(options)
	if err != nil {
		return datamodel.LoadBalancerPool{}, fmt.Errorf("[getPoolMembers] error getting pool members for %s: %w", *vpcPool.ID, err)
	}
	pool := datamodel.NewLoadBalancerPool(vpcPool, members.Members)
	return pool, nil
}

func getPolicyRules(vpcService *vpcv1.VpcV1, lbID, listenerID string,
	lbPolicy *vpcv1.LoadBalancerListenerPolicy) (datamodel.LoadBalancerListenerPolicy, error) {
	options := &vpcv1.ListLoadBalancerListenerPolicyRulesOptions{}
	options.SetLoadBalancerID(lbID)
	options.SetListenerID(listenerID)
	options.SetPolicyID(*lbPolicy.ID)
	rules, _, ruleErr := vpcService.ListLoadBalancerListenerPolicyRules(options)
	if ruleErr != nil {
		return datamodel.LoadBalancerListenerPolicy{},
			fmt.Errorf("[getPolicyRules] error getting rules for %s: %w", *lbPolicy.ID, ruleErr)
	}
	policy := datamodel.NewLoadBalancerListenerPolicy(lbPolicy, rules.Rules)
	return policy, nil
}

func getTransitConnections(tgwService *tgw.TransitGatewayApisV1,
	tgwList []*datamodel.TransitGateway) ([]*datamodel.TransitConnection, error) {
	APIFunc := func(pageSize int64, next *string) (*tgw.TransitConnectionCollection, any, error) {
		return tgwService.ListConnections(&tgw.ListConnectionsOptions{Limit: &pageSize, Start: next})
	}
	getArray := func(collection *tgw.TransitConnectionCollection) []tgw.TransitConnection {
		return collection.Connections
	}

	transitCons, err := iteratePagedAPI(APIFunc, getArray)
	if err != nil {
		return nil, fmt.Errorf("[getTransitConnections] error getting transit connections: %w", err)
	}
	var res []*datamodel.TransitConnection
	for i := range transitCons {
		for j := range tgwList {
			if *(transitCons[i].TransitGateway.ID) == *(tgwList[j].ID) {
				res = append(res, datamodel.NewTransitConnection(&transitCons[i]))
			}
		}
	}
	return res, nil
}

func getTransitGateways(tgwService *tgw.TransitGatewayApisV1, resourceGroupID string) ([]*datamodel.TransitGateway, error) {
	APIFunc := func(pageSize int64, next *string) (*tgw.TransitGatewayCollection, any, error) {
		return tgwService.ListTransitGateways(&tgw.ListTransitGatewaysOptions{Limit: &pageSize, Start: next})
	}
	getArray := func(collection *tgw.TransitGatewayCollection) []tgw.TransitGateway {
		return collection.TransitGateways
	}

	transitGws, err := iteratePagedAPI(APIFunc, getArray)
	if err != nil {
		return nil, fmt.Errorf("[getTransitGateways] error getting transit gateways: %w", err)
	}
	var res []*datamodel.TransitGateway
	for i := range transitGws {
		if resourceGroupID == "" || *(transitGws[i].ResourceGroup.ID) == resourceGroupID {
			res = append(res, datamodel.NewTransitGateway(&transitGws[i]))
		}
	}
	return res, nil
}
