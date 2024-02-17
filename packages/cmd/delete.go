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
	"context"
	"fmt"
	"log/slog"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Delete() removes a given part from the catalog using the
// part id and takes the flag for recursive and forced delete
// recursive and forced delete are currently disabled to
// reduce delete times and protect wrongfull deletion of files
func Delete(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	var argForcedMode bool
	var argRecursiveMode bool
	// cobra command for delete
	deleteCmd := &cobra.Command{
		Use:   "delete [part id]",
		Short: "Delete a part in the Software Parts Catalog.",
		// function to run as a setup on command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part id provided.")
			}
			return nil
		},
		// function to run on command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if the part id is provided as an argument
			argPartID := args[0]
			if argPartID == "" {
				return errors.New("error deleting part, delete subcommand usage: ./ccli delete <catalog_id>")
			}
			// delete the part if the part id is present
			if argPartID != "" {
				slog.Debug("deleting part", slog.String("ID", argPartID))
				if err := graphql.DeletePart(context.Background(), client, argPartID, argRecursiveMode, argForcedMode); err != nil {
					return errors.Wrapf(err, "error deleting part from catalog")
				}
				fmt.Printf("Successfully deleted id: %s from catalog\n", argPartID)
			}
			return nil
		},
	}
	// adding persistent flags for delete i.e. recursive and force
	deleteCmd.PersistentFlags().BoolVarP(&argRecursiveMode, "recursive", "r", false, "To delete parts recursively")
	deleteCmd.PersistentFlags().BoolVarP(&argForcedMode, "force", "f", false, "To delete parts forcefully")
	return deleteCmd
}
