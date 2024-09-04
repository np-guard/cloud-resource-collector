/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"fmt"
	"slices"
	"strings"
)

type Provider string

const (
	AWS Provider = "aws"
	IBM Provider = "ibm"
)

var AllProviders = []string{string(IBM), string(AWS)}
var AllProvidersStr = strings.Join(AllProviders, ", ")

func (p *Provider) String() string {
	return string(*p)
}

func (p *Provider) Set(v string) error {
	v = strings.ToLower(v)
	if slices.Contains(AllProviders, v) {
		*p = Provider(v)
		return nil
	}
	return fmt.Errorf("must be one of [%s]", AllProvidersStr)
}

func (p *Provider) Type() string {
	return "string"
}
