package main

import (
	"flag"
	"fmt"
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
	OutputFile          *string
}

func ParseInArgs(args *InArgs) error {
	SupportedProviders := map[string]bool{
		AWS: true,
		IBM: true,
	}

	args.CollectFromProvider = flag.String("provider", "", "cloud provider from which to collect resources")
	flag.Var(&args.regions, "region", "cloud region from which to collect resources")
	args.getRegions = flag.Bool("get-regions", false, "just print the list of regions for the selected provider")
	args.OutputFile = flag.String("out", "", "file path to store results")
	flag.Parse()

	if !SupportedProviders[*args.CollectFromProvider] {
		flag.PrintDefaults()
		return fmt.Errorf("unsupported provider: %s", *args.CollectFromProvider)
	}

	if *args.CollectFromProvider != IBM {
		if len(args.regions) > 0 {
			return fmt.Errorf("setting regions from the command-line for provider %s is not yet supported. "+
				"Use environment variables or config files instead", *args.CollectFromProvider)
		}
		if *args.getRegions {
			return fmt.Errorf("getting the list of regions for provider %s is not yet supported", *args.CollectFromProvider)
		}
	}

	return nil
}
