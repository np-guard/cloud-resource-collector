package datamodel

import (
	"encoding/json"
	"fmt"
)

// ResourcesContainerModel defines the model of a container for all resource types we can collect
type ResourcesContainerModel struct {
	VpcList           []*VPC             `json:"vpcs"`
	SubnetList        []*Subnet          `json:"subnets"`
	PublicGWList      []*PublicGateway   `json:"public_gateways"`
	FloatingIPList    []*FloatingIP      `json:"floating_ips"`
	NetworkACLList    []*NetworkACL      `json:"network_acls"`
	SecurityGroupList []*SecurityGroup   `json:"security_groups"`
	EndpointGWList    []*EndpointGateway `json:"endpoint_gateways"`
	InstanceList      []*Instance        `json:"instances"`
	RoutingTableList  []*RoutingTable    `json:"routing_tables"`
	LBList            []*LoadBalancer    `json:"load_balancers"`
	IKSWorkerNodes    []*IKSWorkerNode   `json:"iks_worker_nodes"`
}

// NewResourcesContainerModel creates an empty resources container
func NewResourcesContainerModel() *ResourcesContainerModel {
	return &ResourcesContainerModel{
		VpcList:           []*VPC{},
		SubnetList:        []*Subnet{},
		PublicGWList:      []*PublicGateway{},
		FloatingIPList:    []*FloatingIP{},
		NetworkACLList:    []*NetworkACL{},
		SecurityGroupList: []*SecurityGroup{},
		EndpointGWList:    []*EndpointGateway{},
		InstanceList:      []*Instance{},
		RoutingTableList:  []*RoutingTable{},
		LBList:            []*LoadBalancer{},
		IKSWorkerNodes:    []*IKSWorkerNode{},
	}
}

// PrintStats outputs the number of items of each type
func (resources *ResourcesContainerModel) PrintStats() {
	fmt.Printf("Found %d VPCs\n", len(resources.VpcList))
	fmt.Printf("Found %d subnets\n", len(resources.SubnetList))
	fmt.Printf("Found %d public gateways\n", len(resources.PublicGWList))
	fmt.Printf("Found %d floating IPs\n", len(resources.FloatingIPList))
	fmt.Printf("Found %d nACLs\n", len(resources.NetworkACLList))
	fmt.Printf("Found %d security groups\n", len(resources.SecurityGroupList))
	fmt.Printf("Found %d endpoint gateways (VPEs)\n", len(resources.EndpointGWList))
	fmt.Printf("Found %d instances\n", len(resources.InstanceList))
	fmt.Printf("Found %d routing tables\n", len(resources.RoutingTableList))
	fmt.Printf("Found %d load balancers\n", len(resources.LBList))
	fmt.Printf("Found %d IKS worker nodes\n", len(resources.IKSWorkerNodes))
}

// ToJSONString converts a ResourcesContainerModel into a json-formatted-string
func (resources *ResourcesContainerModel) ToJSONString() (string, error) {
	toPrint, err := json.MarshalIndent(resources, "", "    ")
	return string(toPrint), err
}