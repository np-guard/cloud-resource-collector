/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

import (
	"fmt"

	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

const HTTPOK = 200

// Get (the first page of) IKS Clusters
func getClusters(iksService *iksv1.KubernetesServiceApiV1, resourceGroupID string) ([]*string, error) {
	clusterCollection, _, err := iksService.VpcGetClusters(&iksv1.VpcGetClustersOptions{XAuthResourceGroup: &resourceGroupID})
	if err != nil {
		return nil, fmt.Errorf("[getClusters] error getting Clusters: %w", err)
	}

	res := make([]*string, len(clusterCollection))
	for i := range clusterCollection {
		res[i] = clusterCollection[i].ID
	}
	return res, nil
}

// Get all worker nodes of a cluster
func getCLusterNodes(iksService *iksv1.KubernetesServiceApiV1, clusterID *string) ([]*datamodel.IKSWorkerNode, error) {
	workerResponse, _, err := iksService.VpcGetWorkers(&iksv1.VpcGetWorkersOptions{Cluster: clusterID})
	if err != nil {
		return nil, fmt.Errorf("[getClusterNodes] error getting workers for %s: %w", *clusterID, err)
	}

	res := make([]*datamodel.IKSWorkerNode, len(workerResponse))
	for i := range workerResponse {
		res[i] = datamodel.NewIKSWorkerNode(&workerResponse[i])
	}
	return res, nil
}
