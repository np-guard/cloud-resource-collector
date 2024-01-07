/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

func TestUnmarshal(t *testing.T) {
	unmarshalInputs := [4]string{
		"data/iks-on-goldeneye-vpc-config.json",
		"data/demo-with-instances-config.json",
		"data/transit-gateways.json",
		"data/experiments_env.json",
	}

	for i := range unmarshalInputs {
		byteSlice, err := os.ReadFile(unmarshalInputs[i])
		if err != nil {
			t.Errorf("couldn't read file: %s", unmarshalInputs[i])
		}
		config := datamodel.ResourcesContainerModel{}
		err = json.Unmarshal(byteSlice, &config)
		if err != nil {
			t.Errorf("Unmarshal failed with error message: %v", err)
		}
		toPrint, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			t.Errorf("MarshalIndent failed: %v", err)
		}

		if bytes.Equal(byteSlice, toPrint) {
			t.Logf("Unmarshaling successful for %s", unmarshalInputs[i])
		} else {
			t.Errorf("Unmarshaling failed for %s", unmarshalInputs[i])

			// Used for debugging test failures:

			file, err := os.Create("unmarshal_output.json")
			if err != nil {
				t.Errorf("failed with error %v", err)
			}

			_, err = file.WriteString(string(toPrint))
			if err != nil {
				t.Errorf("failed with error %v", err)
			}
		}
	}
}
