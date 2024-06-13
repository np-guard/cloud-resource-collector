/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

const (
	AWS string = "aws"
	IBM string = "ibm"
)

// ResourcesContainerInf is the interface common to all resources containers
type ResourcesContainerInf interface {
	CollectResourcesFromAPI() error
	PrintStats()
	ToJSONString() (string, error)
	AllRegions() []string
	GetResources() ResourcesModel
}

type ResourcesModel interface {
}

type ResourceModelMetadata struct {
	Version  string `json:"collector_version"`
	Provider string `json:"provider"`
}
