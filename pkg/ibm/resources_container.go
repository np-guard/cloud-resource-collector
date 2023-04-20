package ibm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

// tagsClient wraps the global search client and collects tags for all types of resources
type tagsClient struct {
	serviceClient   *globaltaggingv1.GlobalTaggingV1
	listTagsOptions *globaltaggingv1.ListTagsOptions
}

// Constructor for a tagsClient
func newTagsCollector() (*tagsClient, error) {
	serviceClientOptions := &globaltaggingv1.GlobalTaggingV1Options{}
	serviceClient, err := globaltaggingv1.NewGlobalTaggingV1UsingExternalConfig(serviceClientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create tagging service client (%w)", err)
	}

	listTagsOptions := serviceClient.NewListTagsOptions()

	return &tagsClient{serviceClient: serviceClient, listTagsOptions: listTagsOptions}, nil
}

// setResourceTags gets the tags associated with a resource (based on its CRN)
func (tagsCollector *tagsClient) setResourceTags(resource datamodel.TaggedResource) error {
	tagsCollector.listTagsOptions.SetAttachedTo(*resource.GetCRN())
	tagList, _, err := tagsCollector.serviceClient.ListTags(tagsCollector.listTagsOptions)
	if err != nil {
		return fmt.Errorf("failed to collect tags (%w)", err)
	}

	tags := make([]string, len(tagList.Items))
	for i := range tagList.Items {
		tags[i] = *tagList.Items[i].Name
	}
	resource.SetTags(tags)
	return nil
}

// ResourcesContainer holds the results of collecting the configurations of all resources.
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

// collect the tags for all resources of all types
//
//nolint:gocyclo,funlen // because Golang forces me to replicate code per-resource-type
func (resources *ResourcesContainer) collectTags() error {
	// Instantiate the tags collector
	tagsCollector, err := newTagsCollector()
	if err != nil {
		return err
	}

	for i := range resources.VpcList {
		err := tagsCollector.setResourceTags(resources.VpcList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.SubnetList {
		err := tagsCollector.setResourceTags(resources.SubnetList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.PublicGWList {
		err := tagsCollector.setResourceTags(resources.PublicGWList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.FloatingIPList {
		err := tagsCollector.setResourceTags(resources.FloatingIPList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.NetworkACLList {
		err := tagsCollector.setResourceTags(resources.NetworkACLList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.SecurityGroupList {
		err := tagsCollector.setResourceTags(resources.SecurityGroupList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.EndpointGWList {
		err := tagsCollector.setResourceTags(resources.EndpointGWList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.InstanceList {
		err := tagsCollector.setResourceTags(resources.InstanceList[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// CollectResourcesFromAPI uses IBM APIs to collect resource configuration information
//
//nolint:funlen,gocyclo // because Golang forces me to replicate code per-resource-type
func (resources *ResourcesContainer) CollectResourcesFromAPI() error {
	//TODO: Enable supplying credentials through other means
	apiKey := os.Getenv("IBMCLOUD_API_KEY")
	if apiKey == "" {
		return errors.New("no API key set")
	}

	// Instantiate the service with an API key based IAM authenticator
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return errors.New("error creating VPC Service")
	}

	// Get (the first page of) VPCs
	vpcCollection, _, err := vpcService.ListVpcs(&vpcv1.ListVpcsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting VPCs: %w", err)
	}
	resources.VpcList = make([]*datamodel.VPC, len(vpcCollection.Vpcs))
	for i := range vpcCollection.Vpcs {
		resources.VpcList[i] = datamodel.NewVPC(&vpcCollection.Vpcs[i])
	}

	// Get (the first page of) Subnets
	// Note: reserved IPs are collected through a second API call, also without paging
	subnetCollection, _, err := vpcService.ListSubnets(&vpcv1.ListSubnetsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Subnets: %w", err)
	}
	resources.SubnetList = make([]*datamodel.Subnet, len(subnetCollection.Subnets))
	for i := range subnetCollection.Subnets {
		var reservedIPs *vpcv1.ReservedIPCollection
		resources.SubnetList[i] = datamodel.NewSubnet(&subnetCollection.Subnets[i])

		// second API call to get the list of reserved IPs in this subnet
		subnetID := resources.SubnetList[i].ID
		options := vpcService.NewListSubnetReservedIpsOptions(*subnetID)
		reservedIPs, _, err = vpcService.ListSubnetReservedIps(options)
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
	}

	// Get (the first page of) Floating IPs
	floatingIPCollection, _, err := vpcService.ListFloatingIps(&vpcv1.ListFloatingIpsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Floating IPs: %w", err)
	}
	resources.FloatingIPList = make([]*datamodel.FloatingIP, len(floatingIPCollection.FloatingIps))
	for i := range floatingIPCollection.FloatingIps {
		resources.FloatingIPList[i] = datamodel.NewFloatingIP(&floatingIPCollection.FloatingIps[i])
	}

	// Get (the first page of) Network ACLs
	networkACLsCollection, _, err := vpcService.ListNetworkAcls(&vpcv1.ListNetworkAclsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Network ACLs: %w", err)
	}
	resources.NetworkACLList = make([]*datamodel.NetworkACL, len(networkACLsCollection.NetworkAcls))
	for i := range networkACLsCollection.NetworkAcls {
		resources.NetworkACLList[i] = datamodel.NewNetworkACL(&networkACLsCollection.NetworkAcls[i])
	}

	// Get (the first page of) Security Groups
	sgCollection, _, err := vpcService.ListSecurityGroups(&vpcv1.ListSecurityGroupsOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Security Groups: %w", err)
	}
	resources.SecurityGroupList = make([]*datamodel.SecurityGroup, len(sgCollection.SecurityGroups))
	for i := range sgCollection.SecurityGroups {
		resources.SecurityGroupList[i] = datamodel.NewSecurityGroup(&sgCollection.SecurityGroups[i])
	}

	// Get (the first page of) Endpoint Gateways (VPEs)
	vpeCollection, _, err := vpcService.ListEndpointGateways(&vpcv1.ListEndpointGatewaysOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Endpoint Gateways: %w", err)
	}
	resources.EndpointGWList = make([]*datamodel.EndpointGateway, len(vpeCollection.EndpointGateways))
	for i := range vpeCollection.EndpointGateways {
		resources.EndpointGWList[i] = datamodel.NewEndpointGateway(&vpeCollection.EndpointGateways[i])
	}

	// Get (the first page of) Instances
	instancesCollection, _, err := vpcService.ListInstances(&vpcv1.ListInstancesOptions{})
	if err != nil {
		return fmt.Errorf("CollectResourcesFromAPI error getting Instances: %w", err)
	}
	resources.InstanceList = make([]*datamodel.Instance, len(instancesCollection.Instances))
	for i := range instancesCollection.Instances {
		var networkInterfaces *vpcv1.NetworkInterfaceUnpaginatedCollection
		resources.InstanceList[i] = datamodel.NewInstance(&instancesCollection.Instances[i])

		// Second API call to get detailed network interfaces information
		options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
		options.SetInstanceID(*resources.InstanceList[i].ID)
		networkInterfaces, _, err = vpcService.ListInstanceNetworkInterfaces(options)
		if err != nil {
			return fmt.Errorf("CollectResourcesFromAPI error getting NW Interfaces for %s", *resources.InstanceList[i].Name)
		}
		resources.InstanceList[i].NetworkInterfaces = networkInterfaces.NetworkInterfaces
	}

	// Add the tags to all resources
	err = resources.collectTags()
	if err != nil {
		return err
	}

	return nil
}
