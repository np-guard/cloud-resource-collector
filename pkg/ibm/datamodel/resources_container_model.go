/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package datamodel

import (
	"encoding/json"
	"fmt"

	"github.com/np-guard/cloud-resource-collector/pkg/version"
)

// ResourcesContainerModel defines the model of a container for all resource types we can collect
type ResourcesContainerModel struct {
	VpcList               []*VPC               `json:"vpcs"`
	SubnetList            []*Subnet            `json:"subnets"`
	PublicGWList          []*PublicGateway     `json:"public_gateways"`
	FloatingIPList        []*FloatingIP        `json:"floating_ips"`
	NetworkACLList        []*NetworkACL        `json:"network_acls"`
	SecurityGroupList     []*SecurityGroup     `json:"security_groups"`
	EndpointGWList        []*EndpointGateway   `json:"endpoint_gateways"`
	InstanceList          []*Instance          `json:"instances"`
	RoutingTableList      []*RoutingTable      `json:"routing_tables"`
	LBList                []*LoadBalancer      `json:"load_balancers"`
	TransitConnectionList []*TransitConnection `json:"transit_connections"`
	TransitGatewayList    []*TransitGateway    `json:"transit_gateways"`
	IKSClusters           []*IKSCluster        `json:"iks_clusters"`
	Version               string               `json:"collector_version"`
}

// NewResourcesContainerModel creates an empty resources container
func NewResourcesContainerModel() *ResourcesContainerModel {
	return &ResourcesContainerModel{
		VpcList:               []*VPC{},
		SubnetList:            []*Subnet{},
		PublicGWList:          []*PublicGateway{},
		FloatingIPList:        []*FloatingIP{},
		NetworkACLList:        []*NetworkACL{},
		SecurityGroupList:     []*SecurityGroup{},
		EndpointGWList:        []*EndpointGateway{},
		InstanceList:          []*Instance{},
		RoutingTableList:      []*RoutingTable{},
		LBList:                []*LoadBalancer{},
		TransitConnectionList: []*TransitConnection{},
		TransitGatewayList:    []*TransitGateway{},
		IKSClusters:           []*IKSCluster{},
		Version:               version.VersionCore,
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
	fmt.Printf("Found %d transit connections\n", len(resources.TransitConnectionList))
	fmt.Printf("Found %d transit gateways\n", len(resources.TransitGatewayList))
	fmt.Printf("Found %d IKS clusters\n", len(resources.IKSClusters))
}

// ToJSONString converts a ResourcesContainerModel into a json-formatted-string
func (resources *ResourcesContainerModel) ToJSONString() (string, error) {
	toPrint, err := json.MarshalIndent(resources, "", "    ")
	return string(toPrint), err
}
