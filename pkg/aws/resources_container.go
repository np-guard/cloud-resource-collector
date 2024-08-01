/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	aws2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/version"
)

type VPC struct {
	aws2.Vpc
	Region string
}

// ResourcesContainer holds the results of collecting the configurations of all resources.
// This includes: instances, internet gateways, network ACLs, security groups, subnets, and VPCs
type ResourcesContainer struct {
	common.ResourceModelMetadata
	InstancesList      []*aws2.Instance        `json:"instances"`
	InternetGWList     []*aws2.InternetGateway `json:"internet_gateways"`
	NetworkACLsList    []*aws2.NetworkAcl      `json:"network_acls"`
	SecurityGroupsList []*aws2.SecurityGroup   `json:"security_groups"`
	SubnetsList        []*aws2.Subnet          `json:"subnets"`
	VpcsList           []*VPC                  `json:"vpcs"`
	regions            []string
}

// NewResourcesContainer creates an empty resources container
func NewResourcesContainer(regions []string) *ResourcesContainer {
	if len(regions) == 0 {
		regions = awsRegions
	}
	return &ResourcesContainer{
		InstancesList:         []*aws2.Instance{},
		InternetGWList:        []*aws2.InternetGateway{},
		NetworkACLsList:       []*aws2.NetworkAcl{},
		SecurityGroupsList:    []*aws2.SecurityGroup{},
		SubnetsList:           []*aws2.Subnet{},
		VpcsList:              []*VPC{},
		ResourceModelMetadata: common.ResourceModelMetadata{Version: version.VersionCore, Provider: string(common.AWS)},
		regions:               regions,
	}
}

// PrintStats outputs the number of items of each type
func (resources *ResourcesContainer) PrintStats() {
	fmt.Printf("Found %d instances\n", len(resources.InstancesList))
	fmt.Printf("Found %d internet gateways\n", len(resources.InternetGWList))
	fmt.Printf("Found %d nACLs\n", len(resources.NetworkACLsList))
	fmt.Printf("Found %d security groups\n", len(resources.SecurityGroupsList))
	fmt.Printf("Found %d subnets\n", len(resources.SubnetsList))
	fmt.Printf("Found %d VPCs\n", len(resources.VpcsList))
}

// ToJSONString converts a ResourcesContainer into a json-formatted-string
func (resources *ResourcesContainer) ToJSONString() (string, error) {
	toPrint, err := json.MarshalIndent(resources, "", "    ")
	return string(toPrint), err
}

func (resources *ResourcesContainer) AllRegions() []string {
	return awsRegions
}

func (resources *ResourcesContainer) GetResources() common.ResourcesModel {
	return resources
}

func (resources *ResourcesContainer) Fabricate(opts *common.FabricateOptions) { // TODO: implement
}

// CollectResourcesFromAPI uses AWS APIs to collect resource configuration information
func (resources *ResourcesContainer) CollectResourcesFromAPI() error { //nolint:gocyclo // due to many API calls
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI encountered an error loading AWS the configuration: %w", err)
	}

	for _, region := range resources.regions {
		if !slices.Contains(awsRegions, region) {
			log.Printf("Unknown region %s. Available regions for provider aws: %s\n", region, strings.Join(awsRegions, ", "))
			continue
		}

		log.Printf("Collecting resources from region %s\n", region)
		cfg.Region = region
		client := ec2.NewFromConfig(cfg) // Create an Amazon ec2 service client

		// Get (the first page of) VPCs
		vpcsFromAPI, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting VPCs: %w", err)
		}
		for i := range vpcsFromAPI.Vpcs {
			vpc := VPC{Region: region, Vpc: vpcsFromAPI.Vpcs[i]}
			resources.VpcsList = append(resources.VpcsList, &vpc)
		}

		// Get (the first page of) Internet Gateways
		intGWFromAPI, err := client.DescribeInternetGateways(context.TODO(), &ec2.DescribeInternetGatewaysInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting internet gateways: %w", err)
		}
		for i := range intGWFromAPI.InternetGateways {
			resources.InternetGWList = append(resources.InternetGWList, &intGWFromAPI.InternetGateways[i])
		}

		// Get (the first page of) Subnets
		subnetsFromAPI, err := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting subnets: %w", err)
		}
		for i := range subnetsFromAPI.Subnets {
			resources.SubnetsList = append(resources.SubnetsList, &subnetsFromAPI.Subnets[i])
		}

		// Get (the first page of) Network ACLs
		nACLsFromAPI, err := client.DescribeNetworkAcls(context.TODO(), &ec2.DescribeNetworkAclsInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting nACLs: %w", err)
		}
		for i := range nACLsFromAPI.NetworkAcls {
			resources.NetworkACLsList = append(resources.NetworkACLsList, &nACLsFromAPI.NetworkAcls[i])
		}

		// Get (the first page of) Security Groups
		secGroupsFromAPI, err := client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting security groups: %w", err)
		}
		for i := range secGroupsFromAPI.SecurityGroups {
			resources.SecurityGroupsList = append(resources.SecurityGroupsList, &secGroupsFromAPI.SecurityGroups[i])
		}

		// Get (the first page of) Instances
		instancesFromAPI, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting instances: %w", err)
		}
		for i := range instancesFromAPI.Reservations {
			for j := range instancesFromAPI.Reservations[i].Instances {
				resources.InstancesList = append(resources.InstancesList, &instancesFromAPI.Reservations[i].Instances[j])
			}
		}
	}

	return nil
}
