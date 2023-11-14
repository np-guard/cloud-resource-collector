package ibm

import (
	"errors"
	"fmt"
	"os"

	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/go-sdk-core/v5/core"
	tgw "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
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
	datamodel.ResourcesContainerModel
}

// NewResourcesContainer creates an empty resources container
func NewResourcesContainer() *ResourcesContainer {
	return &ResourcesContainer{*datamodel.NewResourcesContainerModel()}
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

	for i := range resources.LBList {
		err := tagsCollector.setResourceTags(resources.LBList[i])
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

	// Instantiate the VPC service with an API key based IAM authenticator
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return errors.New("error creating VPC Service")
	}

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

	// Instantiate the IKS service with an API key based IAM authenticator
	iksService, err := iksv1.NewKubernetesServiceApiV1(&iksv1.KubernetesServiceApiV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return errors.New("error creating IKS Service")
	}

	// VPCs
	resources.VpcList, err = getVPCs(vpcService)
	if err != nil {
		return err
	}

	// Subnets
	resources.SubnetList, err = getSubnets(vpcService)
	if err != nil {
		return err
	}

	// Public Gateways
	resources.PublicGWList, err = getPublicGateways(vpcService)
	if err != nil {
		return err
	}

	// Floating IPs
	resources.FloatingIPList, err = getFloatingIPs(vpcService)
	if err != nil {
		return err
	}

	// Network ACLs
	resources.NetworkACLList, err = getNetworkACLs(vpcService)
	if err != nil {
		return err
	}

	// Security Groups
	resources.SecurityGroupList, err = getSecurityGroups(vpcService)
	if err != nil {
		return err
	}

	// Endpoint Gateways (VPEs)
	resources.EndpointGWList, err = getEndpointGateways(vpcService)
	if err != nil {
		return err
	}

	// Instances
	resources.InstanceList, err = getInstances(vpcService)
	if err != nil {
		return err
	}

	// Routing Tables
	resources.RoutingTableList, err = getRoutingTables(vpcService, resources.VpcList)
	if err != nil {
		return err
	}

	// Load Balancers
	resources.LBList, err = getLoadBalancers(vpcService)
	if err != nil {
		return err
	}

	// Transit Gateways
	resources.TransitConnectionList, err = getTransitConnections(transitGWService)
	if err != nil {
		return err
	}

	// IKS Clusters
	clusterIDs, err := getClusters(iksService)
	if err != nil {
		return err
	}

	// Collect from all clusters
	for i := range clusterIDs {
		// IKS Cluster Nodes
		iksWorkers, nodeErr := getCLusterNodes(iksService, clusterIDs[i])
		if nodeErr != nil {
			return nodeErr
		}
		resources.IKSWorkerNodes = append(resources.IKSWorkerNodes, iksWorkers...)
	}

	// Add the tags to all (taggable) resources
	err = resources.collectTags()
	if err != nil {
		return err
	}

	return nil
}
