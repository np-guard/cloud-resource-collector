package main

import (
	"flag"
	"fmt"
)

//TODO: create an enumerated type for supported providers

type InArgs struct {
	CollectFromProvider *string
	OutputFile          *string
}

func ParseInArgs(args *InArgs) error {
	args.CollectFromProvider = flag.String("provider", "", "cloud provider from which to collect resources")
	args.OutputFile = flag.String("out", "", "file path to store results")
	flag.Parse()

	if *args.CollectFromProvider != "aws" {
		flag.PrintDefaults()
		return fmt.Errorf("unsupported provider, currently supporting: aws")
	}

	return nil
}
