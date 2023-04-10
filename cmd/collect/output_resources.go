package main

import (
	"fmt"
	"log"
	"os"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
)

func OutputResources(rc common.ResourcesContainerInf, outputFileName string) {
	jsonString, err := rc.ToJSONString()
	if err != nil {
		log.Fatal(fmt.Errorf("error converting resources to string: %w", err))
	}

	if outputFileName == "" {
		fmt.Print(jsonString)
	} else {
		log.Printf("Writing to file: %s", outputFileName)

		file, err := os.Create(outputFileName)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(jsonString)
		if err != nil {
			log.Fatal(err)
		}
	}
}
