// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Example() displays a number of potential calls
// that can be made using the ccli.
func Example() *cobra.Command {
	// cobra command for examples
	exampleCmd := &cobra.Command{
		Use:   "examples",
		Short: "Ccli is used to interact with the Software Parts Catalog.",
		// function to be run on command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			// list of all the ccli examples that can be executed
			exampleString :=
				`	$ ccli add part openssl-1.1.1n.yml
	$ ccli add profile profile_openssl-1.1.1n.yml
	$ ccli query "{part(id:\"aR25sd-V8dDvs2-p3Gfae\"){file_verification_code}}"
	$ ccli export part id sdl3ga-naTs42g5-rbow2A -o file.yml
	$ ccli export template security -o file.yml
	$ ccli update openssl-1.1.1n.v4.yml
	$ ccli set openssl-1.1.1n.v4.yml
	$ ccli upload openssl-1.1.1n.tar.gz
	$ ccli find part busybox
	$ ccli find sha256 2493347f59c03...
	$ ccli find profile security werS12-da54FaSff-9U2aef
	$ ccli delete adjb23-A4D3faTa-d95Xufs
	$ ccli ping`
			fmt.Printf("%s\n", exampleString)
			return nil
		},
	}
	return exampleCmd
}
