/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ibm

import (
	"fmt"
	"math/rand"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/np-guard/cloud-resource-collector/pkg/common"
	"github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
	"github.com/np-guard/models/pkg/netset"
)

const (
	ipElementSize     = 256
	defaultCidrPrefix = 24
	maxNumACLsInVPC   = 10
	maxNumRulesInNACL = 10
)

var (
	regionsAndZones = map[string][]string{
		"us-south": {"us-south1", "us-south2", "us-south3"},
		"us-east":  {"us-east1", "us-east2", "us-east3"},
	}
	uid             = map[string]int{}
	availableIPs, _ = netset.IPBlockFromCidrList([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})

	vpcType      = vpcv1.VPCReferenceResourceTypeVPCConst
	ipv4         = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllIPVersionIpv4Const
	allProtocols = vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllProtocolAllConst
)

func getUID(resource string) *string {
	res := fmt.Sprintf("%s-%d", resource, uid[resource])
	uid[resource]++
	return &res
}

func getVPCRef(vpcID *string) *vpcv1.VPCReference {
	return &vpcv1.VPCReference{CRN: vpcID, ID: vpcID, Name: vpcID, ResourceType: &vpcType}
}

func weakRand(upperBound int) int {
	return rand.Intn(upperBound) //nolint:gosec // weak random is ok here
}

func chooseRandElem[T any](pool []T) *T {
	return &pool[weakRand(len(pool))]
}

func getRandomRegion() string {
	regionNum := weakRand(len(regionsAndZones))
	i := 0
	for k := range regionsAndZones {
		if i == regionNum {
			return k
		}
		i++
	}
	return ""
}

func getAvailableInternalCidrBlock() *string {
	prefix := defaultCidrPrefix - weakRand(2)
	baseIP := availableIPs.FirstIPAddress()
	cidr := fmt.Sprintf("%s/%d", baseIP, prefix)
	cidrIPBlock, _ := netset.IPBlockFromCidr(cidr)
	availableIPs = availableIPs.Subtract(cidrIPBlock)
	return &cidr
}

func getRandomCidr() *string {
	var ipElem [4]int
	for i := 0; i < len(ipElem); i++ {
		ipElem[i] = weakRand(ipElementSize)
	}
	prefix := weakRand(2)
	cidr := fmt.Sprintf("%d.%d.%d.%d/%d", ipElem[0], ipElem[1], ipElem[2], ipElem[3], prefix)
	return &cidr
}

func makeNACLRules() []vpcv1.NetworkACLRuleItemIntf {
	res := []vpcv1.NetworkACLRuleItemIntf{}

	numRules := weakRand(maxNumRulesInNACL)
	for i := 0; i < numRules; i++ {
		ruleID := getUID("aclRule")
		rule := vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll{
			ID:   ruleID,
			Name: ruleID,
			Action: chooseRandElem([]string{
				vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionAllowConst,
				vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllActionDenyConst}),
			Direction: chooseRandElem([]string{
				vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllDirectionInboundConst,
				vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAllDirectionOutboundConst,
			}),
			Source:      getRandomCidr(),
			Destination: getRandomCidr(),
			Protocol:    &allProtocols,
			IPVersion:   &ipv4,
		}
		res = append(res, &rule)
	}

	return res
}

func makeNACLs(vpcID string) []*datamodel.NetworkACL {
	numNacls := weakRand(maxNumACLsInVPC) + 1
	res := []*datamodel.NetworkACL{}
	for i := 0; i < numNacls; i++ {
		naclID := getUID("nacl")
		sdkNACL := vpcv1.NetworkACL{ID: naclID, CRN: naclID, Name: naclID, VPC: getVPCRef(&vpcID)}
		sdkNACL.Rules = makeNACLRules()
		modelNacl := datamodel.NewNetworkACL(&sdkNACL)
		res = append(res, modelNacl)
	}

	return res
}

func makePublicGateways(vpcID, vpcRegion string) map[string]*datamodel.PublicGateway { // map from zone to pgw in this zone
	res := map[string]*datamodel.PublicGateway{}
	zones := regionsAndZones[vpcRegion]
	for i := range zones {
		if weakRand(2) == 0 {
			continue // not all zones should get a PGW
		}
		zone := zones[i]
		pgwID := getUID("pgw")
		sdkPGW := vpcv1.PublicGateway{ID: pgwID, CRN: pgwID, Name: pgwID, VPC: getVPCRef(&vpcID), Zone: getZoneRef(zone)}
		res[zone] = datamodel.NewPublicGateway(&sdkPGW)
	}
	return res
}

func getNACLRef(nacl *datamodel.NetworkACL) *vpcv1.NetworkACLReference {
	return &vpcv1.NetworkACLReference{CRN: nacl.CRN, ID: nacl.ID, Name: nacl.Name}
}

func getPGWRef(nacl *datamodel.PublicGateway) *vpcv1.PublicGatewayReference {
	return &vpcv1.PublicGatewayReference{CRN: nacl.CRN, ID: nacl.ID, Name: nacl.Name}
}

func getSubnetRef(subnet *datamodel.Subnet) vpcv1.SubnetReference {
	return vpcv1.SubnetReference{CRN: subnet.CRN, ID: subnet.ID, Name: subnet.Name}
}

func getZoneRef(zone string) *vpcv1.ZoneReference {
	return &vpcv1.ZoneReference{Name: &zone}
}

func (resources *ResourcesContainer) Fabricate(opts *common.FabricateOptions) {
	for i := 0; i < opts.NumVPCs; i++ {
		vpcID := getUID("vpc")
		vpcRegion := getRandomRegion()
		sdkVPC := vpcv1.VPC{ID: vpcID, Name: vpcID, CRN: vpcID}
		vpc := datamodel.NewVPC(&sdkVPC, vpcRegion, nil)
		resources.VpcList = append(resources.VpcList, vpc)

		vpcNACLs := makeNACLs(*vpcID)
		resources.NetworkACLList = append(resources.NetworkACLList, vpcNACLs...)
		vpcPGWs := makePublicGateways(*vpcID, vpcRegion)
		for _, pgw := range vpcPGWs {
			resources.PublicGWList = append(resources.PublicGWList, pgw)
		}

		for s := 0; s < opts.SubnetsPerVPC; s++ {
			subnetID := getUID("subnet")
			subnetZone := chooseRandElem(regionsAndZones[vpcRegion])
			sdkSubnet := vpcv1.Subnet{ID: subnetID, Name: subnetID, CRN: subnetID, VPC: getVPCRef(vpcID), Zone: getZoneRef(*subnetZone)}
			sdkSubnet.Ipv4CIDRBlock = getAvailableInternalCidrBlock()

			subnetNACL := *chooseRandElem(vpcNACLs)
			sdkSubnet.NetworkACL = getNACLRef(subnetNACL)

			subnet := datamodel.NewSubnet(&sdkSubnet, nil)
			subnetNACL.Subnets = append(subnetNACL.Subnets, getSubnetRef(subnet))

			if pgw, ok := vpcPGWs[*subnetZone]; ok && weakRand(2) > 0 { // randomly decide if subnet is connected to a PGW
				subnet.PublicGateway = getPGWRef(pgw)
			}

			resources.SubnetList = append(resources.SubnetList, subnet)
		}
	}
}
