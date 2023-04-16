package main

import (
	"flag"
	"fmt"
)

type InArgs struct {
	CollectFromProvider *string
	OutputFile          *string
}

func ParseInArgs(args *InArgs) error {
	SupportedProviders := map[string]bool{
		AWS: true,
		IBM: true,
	}

	args.CollectFromProvider = flag.String("provider", "", "cloud provider from which to collect resources")
	args.OutputFile = flag.String("out", "", "file path to store results")
	flag.Parse()

	if !SupportedProviders[*args.CollectFromProvider] {
		flag.PrintDefaults()
		return fmt.Errorf("unsupported provider: %s", *args.CollectFromProvider)
	}

	return nil
}
