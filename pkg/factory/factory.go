/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package factory

import (
	"github.com/np-guard/cloud-resource-collector/pkg/aws"
	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/ibm"
)

func GetResourceContainer(provider common.Provider, regions []string, resourceGroup string) common.ResourcesContainerInf {
	var resources common.ResourcesContainerInf
	switch provider {
	case common.AWS:
		resources = aws.NewResourcesContainer(regions)
	case common.IBM:
		resources = ibm.NewResourcesContainer(regions, resourceGroup)
	}
	return resources
}
