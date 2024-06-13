/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"flag"
	"fmt"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
)

type regionList []string

func (dp *regionList) String() string {
	return fmt.Sprintln(*dp)
}

func (dp *regionList) Set(region string) error {
	*dp = append(*dp, region)
	return nil
}

type InArgs struct {
	CollectFromProvider *string
	regions             regionList
	getRegions          *bool
	resourceGroupID     *string
	OutputFile          *string
	version             *bool
}

func ParseInArgs(args *InArgs) error {
	SupportedProviders := map[string]bool{
		common.AWS: true,
		common.IBM: true,
	}

	args.CollectFromProvider = flag.String("provider", "", "cloud provider from which to collect resources")
	flag.Var(&args.regions, "region", "cloud region from which to collect resources")
	args.getRegions = flag.Bool("get-regions", false, "just print the list of regions for the selected provider")
	args.resourceGroupID = flag.String("resource-group", "", "resource group id or name from which to collect resources")
	args.OutputFile = flag.String("out", "", "file path to store results")
	args.version = flag.Bool("version", false, "prints the release version number")
	flag.Parse()

	if !SupportedProviders[*args.CollectFromProvider] && !*args.version {
		flag.PrintDefaults()
		return fmt.Errorf("unsupported provider: %s", *args.CollectFromProvider)
	}

	if *args.CollectFromProvider != common.IBM {
		if len(args.regions) > 0 {
			return fmt.Errorf("setting regions from the command-line for provider %s is not yet supported. "+
				"Use environment variables or config files instead", *args.CollectFromProvider)
		}
		if *args.getRegions {
			return fmt.Errorf("getting the list of regions for provider %s is not yet supported", *args.CollectFromProvider)
		}
		if *args.resourceGroupID != "" {
			return fmt.Errorf("setting resourceGroup from the command-line for provider %s is not yet supported. ", *args.CollectFromProvider)
		}
	}

	return nil
}
