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
	"github.com/np-guard/models/pkg/ipblock"
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
	availableIPs, _ = ipblock.FromCidrList([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	nACLsOfVPC      = map[string][]*datamodel.NetworkACL{}

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

func chooseRandElem[T any](pool []T) *T {
	return &pool[rand.Intn(len(pool))] //nolint:gosec // weak random is ok here
}

func getRandomRegion() string {
	regionNum := rand.Intn(len(regionsAndZones)) //nolint:gosec // weak random is ok here
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
	prefix := defaultCidrPrefix - rand.Intn(2) //nolint:gosec // weak random is ok here
	baseIP := availableIPs.FirstIPAddress()
	cidr := fmt.Sprintf("%s/%d", baseIP, prefix)
	cidrIPBlock, _ := ipblock.FromCidr(cidr)
	availableIPs = availableIPs.Subtract(cidrIPBlock)
	return &cidr
}

func getRandomCidr() *string {
	var ipElem [4]int
	for i := 0; i < len(ipElem); i++ {
		ipElem[i] = rand.Intn(ipElementSize) //nolint:gosec // weak random is ok here
	}
	prefix := rand.Intn(2) //nolint:gosec // weak random is ok here
	cidr := fmt.Sprintf("%d.%d.%d.%d/%d", ipElem[0], ipElem[1], ipElem[2], ipElem[3], prefix)
	return &cidr
}

func makeNACLRules() []vpcv1.NetworkACLRuleItemIntf {
	res := []vpcv1.NetworkACLRuleItemIntf{}

	numRules := rand.Intn(maxNumRulesInNACL) //nolint:gosec // weak random is ok here
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
	numNacls := rand.Intn(maxNumACLsInVPC) + 1 //nolint:gosec // weak random is ok here
	res := []*datamodel.NetworkACL{}
	for i := 0; i < numNacls; i++ {
		naclID := getUID("nacl")
		sdkNACL := vpcv1.NetworkACL{ID: naclID, CRN: naclID, Name: naclID, VPC: getVPCRef(&vpcID)}
		sdkNACL.Rules = makeNACLRules()
		modelNacl := datamodel.NewNetworkACL(&sdkNACL)
		res = append(res, modelNacl)
		nACLsOfVPC[vpcID] = append(nACLsOfVPC[vpcID], modelNacl)
	}

	return res
}

func getNACLRef(nacl *datamodel.NetworkACL) *vpcv1.NetworkACLReference {
	return &vpcv1.NetworkACLReference{CRN: nacl.CRN, ID: nacl.ID, Name: nacl.Name}
}

func getSubnetRef(subnet *datamodel.Subnet) vpcv1.SubnetReference {
	return vpcv1.SubnetReference{CRN: subnet.CRN, ID: subnet.ID, Name: subnet.Name}
}

func (resources *ResourcesContainer) Fabricate(opts *common.FabricateOptions) {
	for i := 0; i < opts.NumVPCs; i++ {
		vpcID := getUID("vpc")
		vpcRegion := getRandomRegion()
		sdkVPC := vpcv1.VPC{ID: vpcID, Name: vpcID, CRN: vpcID}
		vpc := datamodel.NewVPC(&sdkVPC, vpcRegion, nil)
		resources.VpcList = append(resources.VpcList, vpc)

		resources.NetworkACLList = append(resources.NetworkACLList, makeNACLs(*vpcID)...)

		zone := vpcv1.ZoneReference{Name: chooseRandElem(regionsAndZones[vpcRegion])}
		for s := 0; s < opts.SubnetsPerVPC; s++ {
			subnetID := getUID("subnet")
			sdkSubnet := vpcv1.Subnet{ID: subnetID, Name: subnetID, CRN: subnetID, VPC: getVPCRef(vpcID), Zone: &zone}
			sdkSubnet.Ipv4CIDRBlock = getAvailableInternalCidrBlock()
			subnetNACL := *chooseRandElem(nACLsOfVPC[*vpcID])
			sdkSubnet.NetworkACL = getNACLRef(subnetNACL)
			subnet := datamodel.NewSubnet(&sdkSubnet, nil)
			subnetNACL.Subnets = append(subnetNACL.Subnets, getSubnetRef(subnet))
			resources.SubnetList = append(resources.SubnetList, subnet)
		}
	}
}
