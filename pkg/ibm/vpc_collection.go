package ibm

import (
	"fmt"

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

//nolint:dupl // The duplication is essentially creating the adapter
func getVPCs(vpcService *vpcv1.VpcV1) ([]*datamodel.VPC, error) {
	APIFunc := func(pageSize int64, next *string) (*vpcv1.VPCCollection, any, error) {
		return vpcService.ListVpcs(&vpcv1.ListVpcsOptions{Limit: &pageSize, Start: next})
	}
	getArray := func(collection *vpcv1.VPCCollection) []vpcv1.VPC {
		return collection.Vpcs
	}
	vpcs, err := iteratePagedAPI(APIFunc, getArray)
	if err != nil {
		return nil, fmt.Errorf("[getVPCs] error getting VPCs: %w", err)
	}
	res := make([]*datamodel.VPC, len(vpcs))
	for i := range vpcs {
		res[i] = datamodel.NewVPC(&vpcs[i])
	}
	return res, nil
}

func getSubnets(vpcService *vpcv1.VpcV1) ([]*datamodel.Subnet, error) {
	subnetAPIFunc := func(pageSize int64, next *string) (*vpcv1.SubnetCollection, any, error) {
		return vpcService.ListSubnets(&vpcv1.ListSubnetsOptions{Limit: &pageSize, Start: next})
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

//nolint:dupl // See getVPCs
func getPublicGateways(vpcService *vpcv1.VpcV1) ([]*datamodel.PublicGateway, error) {
	gatewayAPIFunc := func(pageSize int64, next *string) (*vpcv1.PublicGatewayCollection, any, error) {
		return vpcService.ListPublicGateways(&vpcv1.ListPublicGatewaysOptions{Limit: &pageSize, Start: next})
	}
	gatewayGetArray := func(collection *vpcv1.PublicGatewayCollection) []vpcv1.PublicGateway {
		return collection.PublicGateways
	}
	gateways, err := iteratePagedAPI(gatewayAPIFunc, gatewayGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getPublicGateways] error getting Public Gateways: %w", err)
	}
	res := make([]*datamodel.PublicGateway, len(gateways))
	for i := range gateways {
		res[i] = datamodel.NewPublicGateway(&gateways[i])
	}

	return res, nil
}

//nolint:dupl // See getVPCs
func getFloatingIPs(vpcService *vpcv1.VpcV1) ([]*datamodel.FloatingIP, error) {
	floatingIPAPIFunc := func(pageSize int64, next *string) (*vpcv1.FloatingIPCollection, any, error) {
		return vpcService.ListFloatingIps(&vpcv1.ListFloatingIpsOptions{Limit: &pageSize, Start: next})
	}
	floatingIPGetArray := func(collection *vpcv1.FloatingIPCollection) []vpcv1.FloatingIP {
		return collection.FloatingIps
	}
	floatingIps, err := iteratePagedAPI(floatingIPAPIFunc, floatingIPGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getFloatingIPs] error getting Floating IPs: %w", err)
	}
	res := make([]*datamodel.FloatingIP, len(floatingIps))
	for i := range floatingIps {
		res[i] = datamodel.NewFloatingIP(&floatingIps[i])
	}

	return res, nil
}

//nolint:dupl // See getVPCs
func getNetworkACLs(vpcService *vpcv1.VpcV1) ([]*datamodel.NetworkACL, error) {
	networkACLAPIFunc := func(pageSize int64, next *string) (*vpcv1.NetworkACLCollection, any, error) {
		return vpcService.ListNetworkAcls(&vpcv1.ListNetworkAclsOptions{Limit: &pageSize, Start: next})
	}
	networkACLGetArray := func(collection *vpcv1.NetworkACLCollection) []vpcv1.NetworkACL {
		return collection.NetworkAcls
	}
	networkAcls, err := iteratePagedAPI(networkACLAPIFunc, networkACLGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getNetworkACLs] error getting Network ACLs: %w", err)
	}
	res := make([]*datamodel.NetworkACL, len(networkAcls))
	for i := range networkAcls {
		res[i] = datamodel.NewNetworkACL(&networkAcls[i])
	}

	return res, nil
}

//nolint:dupl // See getVPCs
func getSecurityGroups(vpcService *vpcv1.VpcV1) ([]*datamodel.SecurityGroup, error) {
	securityGroupAPIFunc := func(pageSize int64, next *string) (*vpcv1.SecurityGroupCollection, any, error) {
		return vpcService.ListSecurityGroups(&vpcv1.ListSecurityGroupsOptions{Limit: &pageSize, Start: next})
	}
	securityGroupGetArray := func(collection *vpcv1.SecurityGroupCollection) []vpcv1.SecurityGroup {
		return collection.SecurityGroups
	}
	securityGroups, err := iteratePagedAPI(securityGroupAPIFunc, securityGroupGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getSecurityGroups] error getting Security Groups: %w", err)
	}
	res := make([]*datamodel.SecurityGroup, len(securityGroups))
	for i := range securityGroups {
		res[i] = datamodel.NewSecurityGroup(&securityGroups[i])
	}

	return res, nil
}

// Get all Endpoint Gateways (VPEs)
//
//nolint:dupl // See getVPCs
func getEndpointGateways(vpcService *vpcv1.VpcV1) ([]*datamodel.EndpointGateway, error) {
	endpointGatewayAPIFunc := func(pageSize int64, next *string) (*vpcv1.EndpointGatewayCollection, any, error) {
		return vpcService.ListEndpointGateways(&vpcv1.ListEndpointGatewaysOptions{Limit: &pageSize, Start: next})
	}
	endpointGatewayGetArray := func(collection *vpcv1.EndpointGatewayCollection) []vpcv1.EndpointGateway {
		return collection.EndpointGateways
	}
	endpointGateways, err := iteratePagedAPI(endpointGatewayAPIFunc, endpointGatewayGetArray)
	if err != nil {
		return nil, fmt.Errorf("[getEndpointGateways] error getting Endpoint Gateways: %w", err)
	}
	res := make([]*datamodel.EndpointGateway, len(endpointGateways))
	for i := range endpointGateways {
		res[i] = datamodel.NewEndpointGateway(&endpointGateways[i])
	}

	return res, nil
}

func getInstances(vpcService *vpcv1.VpcV1) ([]*datamodel.Instance, error) {
	instanceAPIFunc := func(pageSize int64, next *string) (*vpcv1.InstanceCollection, any, error) {
		return vpcService.ListInstances(&vpcv1.ListInstancesOptions{Limit: &pageSize, Start: next})
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

	for i := range vpcList {
		vpcID := *vpcList[i].ID

		routingTables, err := iteratePagedAPI(routingTableAPIFunc(vpcID), routingTableGetArray)
		if err != nil {
			return nil, fmt.Errorf("[getRoutingTables] error getting Routing Tables for %s: %w", vpcID, err)
		}
		for j := range routingTables {
			routes, err := getRoutes(vpcService, vpcID, *routingTables[j].ID)
			if err != nil {
				return nil, err
			}
			res = append(res, datamodel.NewRoutingTable(&routingTables[j], routes))
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
func getLoadBalancers(vpcService *vpcv1.VpcV1) ([]*datamodel.LoadBalancer, error) {
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
	res := make([]*datamodel.LoadBalancer, len(loadBalancers))
	for i := range loadBalancers {
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

		res[i] = datamodel.NewLoadBalancer(&loadBalancers[i], listeners, pools)
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
