package main

import (
	"fmt"
	"os"

	"github.com/np-guard/cloud-resource-collector/pkg/awsCollector"
	"github.com/np-guard/cloud-resource-collector/pkg/common"
)

func main() {
	var inArgs InArgs
	err := ParseInArgs(&inArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing arguments: %v. exiting...\n", err)
		os.Exit(1)
	}

	// Initialize a collector for the requested provider
	var resources common.ResourcesContainerInf
	switch *inArgs.CollectFromProvider {
	case "aws":
		resources = awsCollector.NewResourcesContainer()
	}

	// Collect resources from the provider API and generate output
	resources.CollectResourcesFromAPI()
	OutputResources(resources, *inArgs.OutputFile)

	resources.PrintStats()
}
