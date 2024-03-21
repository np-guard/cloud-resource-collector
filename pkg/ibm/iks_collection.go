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

// Get (the first page of) IKS Clusters and all of it's worker nodes
func getClusters(iksService *iksv1.KubernetesServiceApiV1, resourceGroupID string) ([]*datamodel.IKSCluster, error) {
	clusterCollection, _, err := iksService.VpcGetClusters(&iksv1.VpcGetClustersOptions{XAuthResourceGroup: &resourceGroupID})
	if err != nil {
		return nil, fmt.Errorf("[getClusters] error getting Clusters: %w", err)
	}

	res := make([]*datamodel.IKSCluster, len(clusterCollection))
	for i := range clusterCollection {
		workerResponse, _, err := iksService.VpcGetWorkers(&iksv1.VpcGetWorkersOptions{Cluster: clusterCollection[i].ID})
		if err != nil {
			return nil, fmt.Errorf("[getClusterNodes] error getting workers for %s: %w", *clusterCollection[i].ID, err)
		}
		res[i] = datamodel.NewCluster(&clusterCollection[i], workerResponse)
	}
	return res, nil
}
