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

	fabricatesOpts common.FabricateOptions
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "collector",
		Short:   "cloud-resource-collector is a CLI for collecting VPC-related cloud resources",
		Long:    `cloud-resource-collector uses cloud-provider SDK to gather VPC-related resources defining network connectivity`,
		Version: version.VersionCore,
	}

	providerHelp := fmt.Sprintf("collect resources from an account in this cloud provider. Supported providers: %s", common.AllProvidersStr)
	rootCmd.PersistentFlags().VarP(&provider, providerFlag, "p", providerHelp)
	_ = rootCmd.MarkPersistentFlagRequired(providerFlag)

	rootCmd.PersistentFlags().StringVar(&outputFile, "out", "", "file path to store results")

	rootCmd.AddCommand(newCollectCommand())
	rootCmd.AddCommand(newGetRegionsCommand())
	rootCmd.AddCommand(newFabricateCommand())

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

	return collectCmd
}

func newGetRegionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-regions",
		Short: "List available regions for a given provider",
		Long:  `List all regions that can be used with the --region flag`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			resources := factory.GetResourceContainer(provider, nil, "")
			providerRegions := strings.Join(resources.AllRegions(), ", ")
			fmt.Printf("Available regions for provider %s: %s\n", provider, providerRegions)
			return nil
		},
	}
}

func newFabricateCommand() *cobra.Command {
	fabricateCmd := &cobra.Command{
		Use:   "fabricate",
		Short: "Fabricate synthetic data",
		Long:  `Generates synthetic data with a given number of VPCs, Subnets, VSIs, ...`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			resources := factory.GetResourceContainer(provider, regions, "")
			resources.Fabricate(&fabricatesOpts)
			OutputResources(resources, outputFile)
			return nil
		},
	}
	fabricateCmd.Flags().IntVar(&fabricatesOpts.NumVPCs, "num-vpcs", 1, "Number of VPCs to generate")
	fabricateCmd.Flags().IntVar(&fabricatesOpts.SubnetsPerVPC, "subnets-per-vpc", 1, "Number of subnets to generate in each VPC")

	return fabricateCmd
}
