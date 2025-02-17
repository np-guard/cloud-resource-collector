/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/go-sdk-core/v5/core"
	tgw "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
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
	datamodel.ResourcesContainerModel
	regions         []string
	resourceGroupID string
}

// NewResourcesContainer creates an empty resources container
func NewResourcesContainer(regions []string, resourceGroupID string) *ResourcesContainer {
	if len(regions) == 0 {
		regions = allRegions()
	}

	return &ResourcesContainer{
		ResourcesContainerModel: *datamodel.NewResourcesContainerModel(),
		regions:                 regions,
		resourceGroupID:         resourceGroupID,
	}
}

func (resources *ResourcesContainer) GetResources() common.ResourcesModel {
	return &resources.ResourcesContainerModel
}

func (resources *ResourcesContainer) AllRegions() []string {
	return allRegions()
}

// collect the tags for all resources of all types
//
//nolint:gocyclo // because Golang forces me to replicate code per-resource-type
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

	for i := range resources.VirtualNIList {
		err := tagsCollector.setResourceTags(resources.VirtualNIList[i])
		if err != nil {
			return err
		}
	}

	for i := range resources.LBList {
		err := tagsCollector.setResourceTags(resources.LBList[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (resources *ResourcesContainer) verifyResourceGroupID(apiKey string) error {
	rm, err := resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating resource manager service: %w", err)
	}

	resourceGroup, _, err := rm.GetResourceGroup(rm.NewGetResourceGroupOptions(
		resources.resourceGroupID,
	))
	if err == nil && resourceGroup != nil {
		return nil // user provided us with a valid resource group ID
	}

	// check if the user provided us with the resource group name rather than its id
	listResourceGroupsOptions := rm.NewListResourceGroupsOptions()
	listResourceGroupsOptions.SetName(resources.resourceGroupID)

	resourceGroupList, _, err := rm.ListResourceGroups(&resourcemanagerv2.ListResourceGroupsOptions{Name: &resources.resourceGroupID})
	if err == nil && len(resourceGroupList.Resources) == 1 {
		resources.resourceGroupID = *resourceGroupList.Resources[0].ID
		return nil // user provided us with a valid resource group name
	}
	return fmt.Errorf("error getting Resource Group %s: no resource group found with this ID or name", resources.resourceGroupID)
}

// CollectResourcesFromAPI uses IBM APIs to collect resource configuration information
func (resources *ResourcesContainer) CollectResourcesFromAPI() error {
	//TODO: Enable supplying credentials through other means
	apiKey := os.Getenv("IBMCLOUD_API_KEY")
	if apiKey == "" {
		return errors.New("no API key set")
	}

	// Setup environment variables for Global Tagging Service
	err := os.Setenv("GLOBAL_TAGGING_APIKEY", apiKey)
	if err != nil {
		return errors.New("failed to set GLOBAL_TAGGING_APIKEY")
	}
	err = os.Setenv("GLOBAL_TAGGING_AUTHTYPE", "iam")
	if err != nil {
		return errors.New("failed to set GLOBAL_TAGGING_AUTHTYPE")
	}
	err = os.Setenv("GLOBAL_TAGGING_URL", "https://tags.global-search-tagging.cloud.ibm.com")
	if err != nil {
		return errors.New("failed to set GLOBAL_TAGGING_URL")
	}

	if resources.resourceGroupID != "" {
		err = resources.verifyResourceGroupID(apiKey)
		if err != nil {
			return err
		}
	}

	for _, region := range resources.regions {
		err = resources.collectRegionalResources(region, apiKey)
		if err != nil {
			return err
		}
	}

	err = resources.collectGlobalResources(apiKey)
	if err != nil {
		return err
	}

	return nil
}

//nolint:funlen // function is long because there are many types of resources we collect
func (resources *ResourcesContainer) collectRegionalResources(region, apiKey string) error {
	// check if region is valid
	if _, ok := vpcRegionURLs[region]; !ok {
		log.Printf("Unknown region %s. Available regions for provider ibm: %s\n", region, strings.Join(resources.AllRegions(), ", "))
		return nil
	}

	// Instantiate the VPC service with an API key based IAM authenticator
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
		URL: vpcRegionURLs[region].url,
	})
	if err != nil {
		return errors.New("error creating VPC Service")
	}

	log.Printf("Collecting resources from region %s\n", region)

	// VPCs
	vpcs, err := getVPCs(vpcService, region, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.VpcList = append(resources.VpcList, vpcs...)

	if len(vpcs) == 0 {
		return nil // no point in collecting other resources from this region if it has no VPCs
	}

	// Subnets
	subnets, err := getSubnets(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.SubnetList = append(resources.SubnetList, subnets...)

	// Public Gateways
	pgws, err := getPublicGateways(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.PublicGWList = append(resources.PublicGWList, pgws...)

	// Floating IPs
	fips, err := getFloatingIPs(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.FloatingIPList = append(resources.FloatingIPList, fips...)

	// Network ACLs
	nacls, err := getNetworkACLs(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.NetworkACLList = append(resources.NetworkACLList, nacls...)

	// Security Groups
	sgs, err := getSecurityGroups(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.SecurityGroupList = append(resources.SecurityGroupList, sgs...)

	// Endpoint Gateways (VPEs)
	vpes, err := getEndpointGateways(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.EndpointGWList = append(resources.EndpointGWList, vpes...)

	// Instances
	insts, err := getInstances(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.InstanceList = append(resources.InstanceList, insts...)

	vnis, err := getVirtualNIs(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.VirtualNIList = append(resources.VirtualNIList, vnis...)

	// Routing Tables
	rts, err := getRoutingTables(vpcService, vpcs)
	if err != nil {
		return err
	}
	resources.RoutingTableList = append(resources.RoutingTableList, rts...)

	// Load Balancers
	lbs, err := getLoadBalancers(vpcService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.LBList = append(resources.LBList, lbs...)
	return nil
}

func (resources *ResourcesContainer) collectGlobalResources(apiKey string) error {
	log.Println("Collecting global resources")

	// Transit Gateways
	// Instantiate the Networking service with an API key based IAM authenticator
	var tgServiceVersion = "2021-12-30"
	transitGWService, err := tgw.NewTransitGatewayApisV1(&tgw.TransitGatewayApisV1Options{
		Version: &tgServiceVersion,
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return errors.New("error creating Networking Service")
	}

	err = transitGWService.SetServiceURL("https://transit.cloud.ibm.com/v1")
	if err != nil {
		return errors.New("error setting Networking Service URL")
	}

	resources.TransitGatewayList, err = getTransitGateways(transitGWService, resources.resourceGroupID)
	if err != nil {
		return err
	}

	resources.TransitConnectionList, err = getTransitConnections(transitGWService, resources.TransitGatewayList)
	if err != nil {
		return err
	}

	// Instantiate the IKS service with an API key based IAM authenticator
	iksService, err := iksv1.NewKubernetesServiceApiV1(&iksv1.KubernetesServiceApiV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return errors.New("error creating IKS Service")
	}

	// Collect IKS Clusters
	clusters, err := getClusters(iksService, resources.resourceGroupID)
	if err != nil {
		return err
	}
	resources.IKSClusters = append(resources.IKSClusters, clusters...)

	// Add the tags to all (taggable) resources
	err = resources.collectTags()
	if err != nil {
		return err
	}

	return nil
}
