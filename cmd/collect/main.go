package main

import (
	"log"

	"github.com/np-guard/cloud-resource-collector/pkg/aws"
	"github.com/np-guard/cloud-resource-collector/pkg/common"
)

const (
	AWS string = "aws"
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
	}

	// Collect resources from the provider API and generate output
	err = resources.CollectResourcesFromAPI()
	if err != nil {
		log.Fatal(err)
	}
	OutputResources(resources, *inArgs.OutputFile)

	resources.PrintStats()
}
