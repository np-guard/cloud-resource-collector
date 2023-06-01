package ibm

import (
	"encoding/json"
	"fmt"

	iksv1 "github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

const HTTPOK = 200

// Get (the first page of) IKS Clusters
func getClusters(iksService *iksv1.KubernetesServiceApiV1) ([]*string, error) {
	clusterCollection, _, err := iksService.VpcGetClusters(&iksv1.VpcGetClustersOptions{})
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

// Get worker pools
func getWorkerPools(iksService *iksv1.KubernetesServiceApiV1, clusterID *string) ([]*datamodel.IKSWorkerPool, error) {
	var workerPools []*datamodel.IKSWorkerPool

	workerPoolResponse, response, err := iksService.VpcGetWorkerPools(&iksv1.VpcGetWorkerPoolsOptions{Cluster: clusterID})
	if err != nil {
		if response.StatusCode == HTTPOK {
			// This is a temporary workaround to handle a defect in the SDK in which it does not expect an array
			err = json.Unmarshal(response.RawResult, &workerPools)
			if err != nil {
				return nil, fmt.Errorf("[getWorkerPools] error unmarshelling VpcGetWorkerPools response for %s: %w",
					*clusterID, err)
			}
		} else {
			return nil, fmt.Errorf("[getWorkerPools] error getting worker pools for %s: %w", *clusterID, err)
		}
	} else {
		// SDK returns wrong type, aborting
		return nil, fmt.Errorf("cannot use SDK response %v", workerPoolResponse)
	}

	res := make([]*datamodel.IKSWorkerPool, len(workerPools))
	for i := range workerPools {
		options := iksv1.VpcGetWorkerPoolOptions{
			Cluster:    clusterID,
			Workerpool: workerPools[i].ID,
		}
		pool, _, err := iksService.VpcGetWorkerPool(&options)
		if err != nil {
			return nil, fmt.Errorf("[getWorkerPools] error getting worker pool %s: %w", *workerPools[i].ID, err)
		}
		res[i] = datamodel.NewIKSWorkerPool(pool)
	}

	return workerPools, nil
}
