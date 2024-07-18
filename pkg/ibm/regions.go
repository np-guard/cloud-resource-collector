/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

type ibmRegion struct {
	url       string
	isPrivate bool
}

var vpcRegionURLs = map[string]ibmRegion{
	"us-east":  {"https://us-east.iaas.cloud.ibm.com/v1", false},
	"us-south": {"https://us-south.iaas.cloud.ibm.com/v1", false},
	"ca-tor":   {"https://ca-tor.iaas.cloud.ibm.com/v1", false},
	"br-sao":   {"https://br-sao.iaas.cloud.ibm.com/v1", false},
	"eu-de":    {"https://eu-de.iaas.cloud.ibm.com/v1", false},
	"eu-es":    {"https://eu-es.iaas.cloud.ibm.com/v1", false},
	"eu-gb":    {"https://eu-gb.iaas.cloud.ibm.com/v1", false},
	"eu-fr2":   {"https://eu-fr2.iaas.cloud.ibm.com/v1", true},
	"au-syd":   {"https://au-syd.iaas.cloud.ibm.com/v1", false},
	"jp-osa":   {"https://jp-osa.iaas.cloud.ibm.com/v1", false},
	"jp-tok":   {"https://jp-tok.iaas.cloud.ibm.com/v1", false},
}

// returns a list of shorthand names of all (public) regions
func allRegions() []string {
	regions := []string{}
	for regionName, regionDetails := range vpcRegionURLs {
		if !regionDetails.isPrivate {
			regions = append(regions, regionName)
		}
	}
	return regions
}
