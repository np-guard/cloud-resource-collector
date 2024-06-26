/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"
)

func main() {
	err := newRootCommand().Execute()
	if err != nil {
		os.Exit(1) // error was already printed by Cobra
	}
}
