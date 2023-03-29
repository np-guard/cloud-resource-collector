package awsCollector

import (
	"encoding/json"
	"fmt"
	"os"
)

func Collect(outputFileName string) *ResourcesContainer {
	resources := CollectResourcesFromAPI()
	if outputFileName != "" {
		jsonString, _ := json.MarshalIndent(resources, "", "    ")
		err := os.WriteFile(outputFileName, jsonString, os.ModePerm)
		if err != nil {
			fmt.Print("Something went wrong!")
		}
	}

	return resources
}
