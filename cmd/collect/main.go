/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/np-guard/cloud-resource-collector/pkg/aws"
	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/ibm"
)

const (
	AWS string = "aws"
	IBM string = "ibm"
)

func main() {
	var inArgs InArgs
	err := ParseInArgs(&inArgs)
	if err != nil {
		log.Fatalf("error parsing arguments: %v. exiting...\n", err)
	}

	// Initialize a collector for the requested provider
	var resources common.ResourcesContainerInf
	switch *inArgs.CollectFromProvider {
	case AWS:
		resources = aws.NewResourcesContainer()
	case IBM:
		resources = ibm.NewResourcesContainer(inArgs.regions, *inArgs.resourceGroupID)
	}

	if *inArgs.getRegions {
		providerRegions := strings.Join(resources.AllRegions(), ", ")
		fmt.Printf("Available regions for provider %s: %s\n", *inArgs.CollectFromProvider, providerRegions)
		return
	}

	// Collect resources from the provider API and generate output
	err = resources.CollectResourcesFromAPI()
	if err != nil {
		log.Fatal(err)
	}
	OutputResources(resources, *inArgs.OutputFile)

	resources.PrintStats()
}
