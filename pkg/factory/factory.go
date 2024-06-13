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

func GetResourceContainer(provider string, regions []string, resourceGroup string) common.ResourcesContainerInf {
	var resources common.ResourcesContainerInf
	switch provider {
	case common.AWS:
		resources = aws.NewResourcesContainer()
	case common.IBM:
		resources = ibm.NewResourcesContainer(regions, resourceGroup)
	}
	return resources
}
