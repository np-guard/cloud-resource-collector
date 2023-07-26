package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

const pageSize = 4

type HasGetNextStart[T any] interface {
	GetNextStart() (*string, error)
}

func iteratePagedApi[T any, Q HasGetNextStart[T]](
	apiFunc func(int64, *string) (Q, any, error),
	getArray func(Q) []T) ([]T, error) {
	var start *string = nil
	// We can use collection.TotalCount for efficiency, but usually it's a single page
	res := make([]T, 0)
	for {
		collection, _, err := apiFunc(pageSize, start)
		if err != nil {
			return nil, fmt.Errorf("[getVPCs] error getting VPCs: %w", err)
		}
		res = append(res, getArray(collection)...)
		start, err = collection.GetNextStart()
		if err != nil {
			return nil, fmt.Errorf("[getVPCs] error getting next page: %w", err)
		}
		if start == nil {
			break
		}
	}
	return res, nil
}

// Get all VPCs
func getVPCs(vpcService *vpcv1.VpcV1) ([]*datamodel.VPC, error) {
	apiFunc := func(pageSize int64, next *string) (*vpcv1.VPCCollection, any, error) {
		return vpcService.ListVpcs(&vpcv1.ListVpcsOptions{Limit: &pageSize, Start: next})
	}
	getArray := func(collection *vpcv1.VPCCollection) []vpcv1.VPC {
		return collection.Vpcs
	}
	// We can use vpcCollection.TotalCount for efficiency, but usually it's a single page
	apiResult, err := iteratePagedApi(apiFunc, getArray)
	if err != nil {
		return nil, fmt.Errorf("[getVPCs] error getting VPCs: %w", err)
	}
	res := make([]*datamodel.VPC, 0)
	for _, vpc := range apiResult {
		res = append(res, datamodel.NewVPC(&vpc))
	}
	return res, nil
}

// Get (the first page of) Subnets
// Note: reserved IPs are collected through a second API call, also without paging
func getSubnets(vpcService *vpcv1.VpcV1) ([]*datamodel.Subnet, error) {
	subnetCollection, _, err := vpcService.ListSubnets(&vpcv1.ListSubnetsOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getSubnets] error getting Subnets: %w", err)
	}
	res := make([]*datamodel.Subnet, len(subnetCollection.Subnets))
	for i := range subnetCollection.Subnets {
		var reservedIPs *vpcv1.ReservedIPCollection
		res[i] = datamodel.NewSubnet(&subnetCollection.Subnets[i])

		// second API call to get the list of reserved IPs in this subnet
		subnetID := res[i].ID
		options := vpcService.NewListSubnetReservedIpsOptions(*subnetID)
		reservedIPs, _, err = vpcService.ListSubnetReservedIps(options)
		if err != nil {
			return nil, fmt.Errorf("[getSubnets]: error getting reserved IPs for %s",
				*res[i].Name)
		}
		res[i].ReservedIps = reservedIPs.ReservedIps
	}

	return res, nil
}

// Get (the first page of) Public Gateways
func getPublicGateways(vpcService *vpcv1.VpcV1) ([]*datamodel.PublicGateway, error) {
	publicGWCollection, _, err := vpcService.ListPublicGateways(&vpcv1.ListPublicGatewaysOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getPublicGateways] error getting Public Gateways: %w", err)
	}
	res := make([]*datamodel.PublicGateway, len(publicGWCollection.PublicGateways))
	for i := range publicGWCollection.PublicGateways {
		res[i] = datamodel.NewPublicGateway(&publicGWCollection.PublicGateways[i])
	}

	return res, nil
}

// Get (the first page of) Floating IPs
func getFloatingIPs(vpcService *vpcv1.VpcV1) ([]*datamodel.FloatingIP, error) {
	floatingIPCollection, _, err := vpcService.ListFloatingIps(&vpcv1.ListFloatingIpsOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getFloatingIPs] error getting Floating IPs: %w", err)
	}
	res := make([]*datamodel.FloatingIP, len(floatingIPCollection.FloatingIps))
	for i := range floatingIPCollection.FloatingIps {
		res[i] = datamodel.NewFloatingIP(&floatingIPCollection.FloatingIps[i])
	}

	return res, nil
}

// Get (the first page of) Network ACLs
func getNetworkACLs(vpcService *vpcv1.VpcV1) ([]*datamodel.NetworkACL, error) {
	networkACLsCollection, _, err := vpcService.ListNetworkAcls(&vpcv1.ListNetworkAclsOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getNetworkACLs] error getting Network ACLs: %w", err)
	}
	res := make([]*datamodel.NetworkACL, len(networkACLsCollection.NetworkAcls))
	for i := range networkACLsCollection.NetworkAcls {
		res[i] = datamodel.NewNetworkACL(&networkACLsCollection.NetworkAcls[i])
	}

	return res, nil
}

// Get (the first page of) Security Groups
func getSecurityGroups(vpcService *vpcv1.VpcV1) ([]*datamodel.SecurityGroup, error) {
	sgCollection, _, err := vpcService.ListSecurityGroups(&vpcv1.ListSecurityGroupsOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getSecurityGroups] error getting Security Groups: %w", err)
	}
	res := make([]*datamodel.SecurityGroup, len(sgCollection.SecurityGroups))
	for i := range sgCollection.SecurityGroups {
		res[i] = datamodel.NewSecurityGroup(&sgCollection.SecurityGroups[i])
	}

	return res, nil
}

// Get (the first page of) Endpoint Gateways (VPEs)
func getEndpointGateways(vpcService *vpcv1.VpcV1) ([]*datamodel.EndpointGateway, error) {
	vpeCollection, _, err := vpcService.ListEndpointGateways(&vpcv1.ListEndpointGatewaysOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getEndpointGateways] error getting Endpoint Gateways: %w", err)
	}
	res := make([]*datamodel.EndpointGateway, len(vpeCollection.EndpointGateways))
	for i := range vpeCollection.EndpointGateways {
		res[i] = datamodel.NewEndpointGateway(&vpeCollection.EndpointGateways[i])
	}

	return res, nil
}

// Get (the first page of) Instances
func getInstances(vpcService *vpcv1.VpcV1) ([]*datamodel.Instance, error) {
	instancesCollection, _, err := vpcService.ListInstances(&vpcv1.ListInstancesOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getInstances] error getting Instances: %w", err)
	}
	res := make([]*datamodel.Instance, len(instancesCollection.Instances))
	for i := range instancesCollection.Instances {
		var networkInterfaces *vpcv1.NetworkInterfaceUnpaginatedCollection
		res[i] = datamodel.NewInstance(&instancesCollection.Instances[i])

		// Second API call to get detailed network interfaces information
		options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
		options.SetInstanceID(*res[i].ID)
		networkInterfaces, _, err = vpcService.ListInstanceNetworkInterfaces(options)
		if err != nil {
			return nil, fmt.Errorf("[getInstances] error getting NW Interfaces for %s", *res[i].Name)
		}
		res[i].NetworkInterfaces = networkInterfaces.NetworkInterfaces
	}

	return res, nil
}

// Get (the first page of) Routing Tables
func getRoutingTables(vpcService *vpcv1.VpcV1, vpcList []*datamodel.VPC) ([]*datamodel.RoutingTable, error) {
	var res []*datamodel.RoutingTable

	for i := range vpcList {
		vpcID := vpcList[i].ID
		options := vpcService.NewListVPCRoutingTablesOptions(*vpcID)
		routingTableCollection, _, err := vpcService.ListVPCRoutingTables(options)
		if err != nil {
			return nil, fmt.Errorf("[getRoutingTables] error getting Routing Tables for %s: %w", *vpcID, err)
		}
		for j := range routingTableCollection.RoutingTables {
			rtID := routingTableCollection.RoutingTables[j].ID
			options := vpcService.NewListVPCRoutingTableRoutesOptions(*vpcID, *rtID)
			routeCollection, _, err := vpcService.ListVPCRoutingTableRoutes(options)
			if err != nil {
				return nil, fmt.Errorf("[getRoutingTables] error getting routes for %s: %w", *rtID, err)
			}
			res = append(res,
				datamodel.NewRoutingTable(&routingTableCollection.RoutingTables[i], routeCollection.Routes))
		}
	}

	return res, nil
}

// Get (the first page of) Load Balancers
func getLoadBalancers(vpcService *vpcv1.VpcV1) ([]*datamodel.LoadBalancer, error) {
	lbCollection, _, err := vpcService.ListLoadBalancers(&vpcv1.ListLoadBalancersOptions{})
	if err != nil {
		return nil, fmt.Errorf("[getLoadBalancers] error getting Load Balancer: %w", err)
	}
	res := make([]*datamodel.LoadBalancer, len(lbCollection.LoadBalancers))
	for i := range lbCollection.LoadBalancers {
		// get all the listeners
		lbID := *lbCollection.LoadBalancers[i].ID
		listenerOptions := &vpcv1.ListLoadBalancerListenersOptions{}
		listenerOptions.SetLoadBalancerID(lbID)
		listenersCollection, _, err := vpcService.ListLoadBalancerListeners(listenerOptions)
		if err != nil {
			return nil, fmt.Errorf("[getLoadBalancers] error getting listeners for %s: %w", lbID, err)
		}

		listeners := make([]datamodel.LoadBalancerListener, len(listenersCollection.Listeners))
		for j := range listenersCollection.Listeners {
			listenerID := listenersCollection.Listeners[j].ID
			policiesOptions := &vpcv1.ListLoadBalancerListenerPoliciesOptions{}
			policiesOptions.SetLoadBalancerID(lbID)
			policiesOptions.SetListenerID(*listenerID)
			policiesCollection, _, polErr := vpcService.ListLoadBalancerListenerPolicies(policiesOptions)
			if polErr != nil {
				return nil, fmt.Errorf("[getLoadBalancers] error getting policies for %s: %w", *listenerID, err)
			}

			policies := make([]datamodel.LoadBalancerListenerPolicy, len(policiesCollection.Policies))
			for k := range policiesCollection.Policies {
				policyID := policiesCollection.Policies[k].ID
				options := &vpcv1.ListLoadBalancerListenerPolicyRulesOptions{}
				options.SetLoadBalancerID(lbID)
				options.SetListenerID(*listenerID)
				options.SetPolicyID(*policyID)
				rules, _, ruleErr := vpcService.ListLoadBalancerListenerPolicyRules(options)
				if ruleErr != nil {
					return nil, fmt.Errorf("[getLoadBalancers] error getting policy rules for %s: %w", *policyID, err)
				}
				policies[k] = datamodel.NewLoadBalancerListenerPolicy(&policiesCollection.Policies[k], rules.Rules)
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
			poolID := poolsCollection.Pools[j].ID
			options := &vpcv1.ListLoadBalancerPoolMembersOptions{}
			options.SetLoadBalancerID(lbID)
			options.SetPoolID(*poolID)
			members, _, err := vpcService.ListLoadBalancerPoolMembers(options)
			if err != nil {
				return nil, fmt.Errorf("[getLoadBalancers] error getting pool members for %s: %w", *poolID, err)
			}
			pools[j] = datamodel.NewLoadBalancerPool(&poolsCollection.Pools[j], members.Members)
		}

		res[i] = datamodel.NewLoadBalancer(&lbCollection.LoadBalancers[i], listeners, pools)
	}

	return res, nil
}
