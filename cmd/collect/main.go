package main

import (
	"github.com/np-guard/cloud-resource-collector/pkg/awsCollector"
)

func main() {
	resources := awsCollector.Collect("config_object.json")
	//fmt.Printf("the result: \n %v \n", resources)
	resources.PrintStats()
}
