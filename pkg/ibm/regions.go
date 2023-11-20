/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

var vpcRegionURLs = map[string]string{
	"us-east":  "https://us-east.iaas.cloud.ibm.com/v1",
	"us-south": "https://us-south.iaas.cloud.ibm.com/v1",
	"ca-tor":   "https://ca-tor.iaas.cloud.ibm.com/v1",
	"br-sao":   "https://br-sao.iaas.cloud.ibm.com/v1",
	"eu-de":    "https://eu-de.iaas.cloud.ibm.com/v1",
	"eu-es":    "https://eu-es.iaas.cloud.ibm.com/v1",
	"eu-gb":    "https://eu-gb.iaas.cloud.ibm.com/v1",
	"au-syd":   "https://au-syd.iaas.cloud.ibm.com/v1",
	"jp-osa":   "https://jp-osa.iaas.cloud.ibm.com/v1",
	"jp-tok":   "https://jp-tok.iaas.cloud.ibm.com/v1",
}

func allRegions() []string {
	regions := make([]string, len(vpcRegionURLs))
	i := 0
	for region := range vpcRegionURLs {
		regions[i] = region
		i++
	}
	return regions
}
