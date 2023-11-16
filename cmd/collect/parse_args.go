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
	OutputFile          *string
}

func ParseInArgs(args *InArgs) error {
	SupportedProviders := map[string]bool{
		AWS: true,
		IBM: true,
	}

	args.CollectFromProvider = flag.String("provider", "", "cloud provider from which to collect resources")
	args.OutputFile = flag.String("out", "", "file path to store results")
	flag.Var(&args.regions, "region", "cloud region from which to collect resources")
	flag.Parse()

	if !SupportedProviders[*args.CollectFromProvider] {
		flag.PrintDefaults()
		return fmt.Errorf("unsupported provider: %s", *args.CollectFromProvider)
	}

	return nil
}
