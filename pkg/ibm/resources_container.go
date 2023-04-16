package ibm

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"log"
	"os"
)

// tagsClient wraps the global search client and collects tags for all types of resources
type tagsClient struct {
	serviceClient   *globaltaggingv1.GlobalTaggingV1
	listTagsOptions *globaltaggingv1.ListTagsOptions
}
// Constructor for a tagsClient
func newTagsCollector() *tagsClient {

	serviceClientOptions := &globaltaggingv1.GlobalTaggingV1Options{}
	serviceClient, err := globaltaggingv1.NewGlobalTaggingV1UsingExternalConfig(serviceClientOptions)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create tagging service client (%w)", err))
	}

	listTagsOptions := serviceClient.NewListTagsOptions()

	return &tagsClient{serviceClient: serviceClient, listTagsOptions: listTagsOptions}
}
// collectTags gets the tags associated with a resource (based on its CRN)
func (tagsCollector *tagsClient) collectTags(resourceID string) []string {

	tagsCollector.listTagsOptions.SetAttachedTo(resourceID)
	tagList, _, err := tagsCollector.serviceClient.ListTags(tagsCollector.listTagsOptions)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to collect tags (%w)", err))
	}

	tags := make([]string, len(tagList.Items))
	for i := range tagList.Items {
		tags[i] = *tagList.Items[i].Name
	}

	return tags
}

// ResourcesContainer holds the results of collecting the configurations of all resources.
// This includes:
type ResourcesContainer struct {
	VpcList           []*datamodel.VPC             `json:"vpcs"`
	SubnetList        []*datamodel.Subnet          `json:"subnets"`
	PublicGWList      []*datamodel.PublicGateway   `json:"public_gateways"`
	FloatingIPList    []*datamodel.FloatingIP      `json:"floating_ips"`
	NetworkACLList    []*datamodel.NetworkACL      `json:"network_acls"`
	SecurityGroupList []*datamodel.SecurityGroup   `json:"security_groups"`
	EndpointGWList    []*datamodel.EndpointGateway `json:"endpoint_gateways"`
	InstanceList      []*datamodel.Instance        `json:"instances"`
}

// NewResourcesContainer creates an empty resources container
func NewResourcesContainer() *ResourcesContainer {
	return &ResourcesContainer{
		VpcList:           []*datamodel.VPC{},
		SubnetList:        []*datamodel.Subnet{},
		PublicGWList:      []*datamodel.PublicGateway{},
		FloatingIPList:    []*datamodel.FloatingIP{},
		NetworkACLList:    []*datamodel.NetworkACL{},
		SecurityGroupList: []*datamodel.SecurityGroup{},
		EndpointGWList:    []*datamodel.EndpointGateway{},
		InstanceList:      []*datamodel.Instance{},
	}
}

// PrintStats outputs the number of items of each type
func (resources *ResourcesContainer) PrintStats() {
	fmt.Printf("Found %d VPCs\n", len(resources.VpcList))
	fmt.Printf("Found %d subnets\n", len(resources.SubnetList))
	fmt.Printf("Found %d public gateways\n", len(resources.PublicGWList))
	fmt.Printf("Found %d floating IPs\n", len(resources.FloatingIPList))
	fmt.Printf("Found %d nACLs\n", len(resources.NetworkACLList))
	fmt.Printf("Found %d security groups\n", len(resources.SecurityGroupList))
	fmt.Printf("Found %d endpoint gateways (VPEs)\n", len(resources.EndpointGWList))
	fmt.Printf("Found %d instances\n", len(resources.InstanceList))
}

// ToJSONString converts a ResourcesContainer into a json-formatted-string
func (resources *ResourcesContainer) ToJSONString() (string, error) {
	toPrint, err := json.MarshalIndent(resources, "", "    ")
	return string(toPrint), err
}

// CollectResourcesFromAPI uses IBM APIs to collect resource configuration information
//
//nolint:all
func (resources *ResourcesContainer) CollectResourcesFromAPI() error {

	//TODO: Enable supplying credentials through other means
	apiKey := os.Getenv("IBMCLOUD_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key set")
	}

	// Instantiate the service with an API key based IAM authenticator
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		log.Fatal("Error creating VPC Service.")
	}

	// Instantiate the tags collector
	tagsCollector := newTagsCollector()

	// Get (the first page of) VPCs
	vpcCollection, _, err := vpcService.ListVpcs(&vpcv1.ListVpcsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting VPCs: %w", err)
	}
	resources.VpcList = make([]*datamodel.VPC, len(vpcCollection.Vpcs))
	for i := range vpcCollection.Vpcs {
		resources.VpcList[i] = datamodel.NewVPC(&vpcCollection.Vpcs[i])
		resources.VpcList[i].Tags = tagsCollector.collectTags(*resources.VpcList[i].CRN)
	}

	// Get (the first page of) Subnets
	// Note: reserved IPs are collected through a second API call, also without paging
	subnetCollection, _, err := vpcService.ListSubnets(&vpcv1.ListSubnetsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Subnets: %w", err)
	}
	resources.SubnetList = make([]*datamodel.Subnet, len(subnetCollection.Subnets))
	for i := range subnetCollection.Subnets {
		resources.SubnetList[i] = datamodel.NewSubnet(&subnetCollection.Subnets[i])
		resources.SubnetList[i].Tags = tagsCollector.collectTags(*resources.SubnetList[i].CRN)

		// second API call to get the list of reserved IPs in this subnet
		subnetID := resources.SubnetList[i].ID
		options := vpcService.NewListSubnetReservedIpsOptions(*subnetID)
		reservedIPs, _, err := vpcService.ListSubnetReservedIps(options)
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting reserved IPs for %s",
				*resources.SubnetList[i].Name)
		}
		resources.SubnetList[i].ReservedIps = reservedIPs.ReservedIps
	}

	// Get (the first page of) Public Gateways
	publicGWCollection, _, err := vpcService.ListPublicGateways(&vpcv1.ListPublicGatewaysOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Public Gateways: %w", err)
	}
	resources.PublicGWList = make([]*datamodel.PublicGateway, len(publicGWCollection.PublicGateways))
	for i := range publicGWCollection.PublicGateways {
		resources.PublicGWList[i] = datamodel.NewPublicGateway(&publicGWCollection.PublicGateways[i])
		resources.PublicGWList[i].Tags = tagsCollector.collectTags(*resources.PublicGWList[i].CRN)
	}

	// Get (the first page of) Floating IPs
	floatingIPCollection, _, err := vpcService.ListFloatingIps(&vpcv1.ListFloatingIpsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Floating IPs: %w", err)
	}
	resources.FloatingIPList = make([]*datamodel.FloatingIP, len(floatingIPCollection.FloatingIps))
	for i := range floatingIPCollection.FloatingIps {
		resources.FloatingIPList[i] = datamodel.NewFloatingIP(&floatingIPCollection.FloatingIps[i])
		resources.FloatingIPList[i].Tags = tagsCollector.collectTags(*resources.FloatingIPList[i].CRN)
	}

	// Get (the first page of) Network ACLs
	networkACLsCollection, _, err := vpcService.ListNetworkAcls(&vpcv1.ListNetworkAclsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Network ACLs: %w", err)
	}
	resources.NetworkACLList = make([]*datamodel.NetworkACL, len(networkACLsCollection.NetworkAcls))
	for i := range networkACLsCollection.NetworkAcls {
		resources.NetworkACLList[i] = datamodel.NewNetworkACL(&networkACLsCollection.NetworkAcls[i])
		resources.NetworkACLList[i].Tags = tagsCollector.collectTags(*resources.NetworkACLList[i].CRN)
	}

	// Get (the first page of) Security Groups
	sgCollection, _, err := vpcService.ListSecurityGroups(&vpcv1.ListSecurityGroupsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Security Groups: %w", err)
	}
	resources.SecurityGroupList = make([]*datamodel.SecurityGroup, len(sgCollection.SecurityGroups))
	for i := range sgCollection.SecurityGroups {
		resources.SecurityGroupList[i] = datamodel.NewSecurityGroup(&sgCollection.SecurityGroups[i])
		resources.SecurityGroupList[i].Tags = tagsCollector.collectTags(*resources.SecurityGroupList[i].CRN)
	}

	// Get (the first page of) Endpoint Gateways (VPEs)
	vpeCollection, _, err := vpcService.ListEndpointGateways(&vpcv1.ListEndpointGatewaysOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Endpoint Gateways: %w", err)
	}
	resources.EndpointGWList = make([]*datamodel.EndpointGateway, len(vpeCollection.EndpointGateways))
	for i := range vpeCollection.EndpointGateways {
		resources.EndpointGWList[i] = datamodel.NewEndpointGateway(&vpeCollection.EndpointGateways[i])
		resources.EndpointGWList[i].Tags = tagsCollector.collectTags(*resources.EndpointGWList[i].CRN)
	}

	// Get (the first page of) Instances
	instancesCollection, _, err := vpcService.ListInstances(&vpcv1.ListInstancesOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Endpoint Gateways: %w", err)
	}
	resources.InstanceList = make([]*datamodel.Instance, len(instancesCollection.Instances))
	for i := range instancesCollection.Instances {
		resources.InstanceList[i] = datamodel.NewInstance(&instancesCollection.Instances[i])
		resources.InstanceList[i].Tags = tagsCollector.collectTags(*resources.InstanceList[i].CRN)

		// Second API call to get detailed network interfaces information
		options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
		options.SetInstanceID(*resources.InstanceList[i].ID)
		networkInterfaces, _, err := vpcService.ListInstanceNetworkInterfaces(options)
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting NW Interfaces for %s", *resources.InstanceList[i].Name)
		}
		resources.InstanceList[i].NetworkInterfaces = networkInterfaces.NetworkInterfaces
	}

	return nil
}
