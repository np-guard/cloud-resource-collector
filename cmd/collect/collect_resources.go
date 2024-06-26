/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/factory"
)

func collectResources(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true // if we got this far, flags are syntactically correct, so no need to print usage

	if provider != common.IBM {
		if len(regions) > 0 {
			return fmt.Errorf("setting regions from the command-line for provider %s is not yet supported. "+
				"Use environment variables or config files instead", provider)
		}
		if resourceGroupID != "" {
			return fmt.Errorf("setting resource-group from the command-line for provider %s is not yet supported. ", provider)
		}
	}

	resources := factory.GetResourceContainer(provider, regions, resourceGroupID)
	// Collect resources from the provider API and generate output
	err := resources.CollectResourcesFromAPI()
	if err != nil {
		return err
	}
	OutputResources(resources, outputFile)
	resources.PrintStats()
	return nil
}
