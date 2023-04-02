package main

import (
	"fmt"
	"log"
	"os"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
)

func OutputResources(rc common.ResourcesContainerInf, outputFileName string) {
	jsonString, err := rc.ToString()
	if err != nil {
		log.Fatal(err)
	}

	if outputFileName == "" {
		fmt.Print(jsonString)
	} else {

		file, err := os.Create(outputFileName)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(jsonString)
		if err != nil {
			fmt.Print(err)
		}
	}

}
