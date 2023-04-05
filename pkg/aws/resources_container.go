package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	aws2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ResourcesContainer holds the results of collecting the configurations of all resources.
// This includes: instances, internet gateways, network ACLs, security groups, subnets, and VPCs
type ResourcesContainer struct {
	InstancesList      []*aws2.Instance        `json:"instances"`
	InternetGWList     []*aws2.InternetGateway `json:"internet_gateways"`
	NetworkACLsList    []*aws2.NetworkAcl      `json:"network_acls"`
	SecurityGroupsList []*aws2.SecurityGroup   `json:"security_groups"`
	SubnetsList        []*aws2.Subnet          `json:"subnets"`
	VpcsList           []*aws2.Vpc             `json:"vpcs"`
}

// NewResourcesContainer creates an empty resources container
func NewResourcesContainer() *ResourcesContainer {
	return &ResourcesContainer{
		InstancesList:      []*aws2.Instance{},
		InternetGWList:     []*aws2.InternetGateway{},
		NetworkACLsList:    []*aws2.NetworkAcl{},
		SecurityGroupsList: []*aws2.SecurityGroup{},
		SubnetsList:        []*aws2.Subnet{},
		VpcsList:           []*aws2.Vpc{},
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

// ToString converts a ResourcesContainer into a json-formatted-string
func (resources *ResourcesContainer) ToString() (string, error) {
	toPrint, err := json.MarshalIndent(resources, "", "    ")
	return string(toPrint), err
}

// CollectResourcesFromAPI uses AWS APIs to collect resource configuration information
func (resources *ResourcesContainer) CollectResourcesFromAPI() error {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI encountered an error loading AWS the configuration: %w\n", err)
	}

	// Create an Amazon ec2 service client
	client := ec2.NewFromConfig(cfg)

	// Get (the first page of) VPCs
	vpcsFromAPI, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting VPCs: %w", err)
	}
	for i := range vpcsFromAPI.Vpcs {
		resources.VpcsList = append(resources.VpcsList, &vpcsFromAPI.Vpcs[i])
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

	return nil
}
