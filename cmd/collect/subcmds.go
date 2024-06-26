/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/factory"
	"github.com/np-guard/cloud-resource-collector/pkg/version"
)

const (
	providerFlag = "provider"
)

var (
	provider        common.Provider
	regions         []string
	resourceGroupID string
	outputFile      string
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "collector",
		Short:   "cloud-resource-collector is a CLI for collecting VPC-related cloud resources",
		Long:    `cloud-resource-collector uses cloud-provider SDK to gather VPC-related resources defining network connectivity`,
		Version: version.VersionCore,
	}

	rootCmd.PersistentFlags().VarP(&provider, providerFlag, "p", "collect resources from an account in this cloud provider")
	_ = rootCmd.MarkPersistentFlagRequired(providerFlag)

	rootCmd.AddCommand(newCollectCommand())
	rootCmd.AddCommand(newGetRegionsCommand())

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) // disable help command. should use --help flag instead

	return rootCmd
}

func newCollectCommand() *cobra.Command {
	collectCmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect VPC-related cloud resources",
		Long:  `Use cloud-provider SDK to gather VPC-related resources defining network connectivity`,
		Args:  cobra.NoArgs,
		RunE:  collectResources,
	}

	collectCmd.Flags().StringArrayVarP(&regions, "region", "r", nil, "cloud region from which to collect resources")
	collectCmd.Flags().StringVar(&resourceGroupID, "resource-group", "", "resource group id or name from which to collect resources")
	collectCmd.Flags().StringVar(&outputFile, "out", "", "file path to store results")

	return collectCmd
}

func newGetRegionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-regions",
		Short: "List available regions for a given provider",
		Long:  `List all regions that can be used with the --region flag`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if provider != common.IBM {
				return fmt.Errorf("command not supported for provider %s", provider)
			}
			resources := factory.GetResourceContainer(provider, nil, "")
			providerRegions := strings.Join(resources.AllRegions(), ", ")
			fmt.Printf("Available regions for provider %s: %s\n", provider, providerRegions)
			return nil
		},
	}
}
