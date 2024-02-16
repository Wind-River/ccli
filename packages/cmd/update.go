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
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/yaml"

	graph "github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Update() is a sub command responsible for updating part information
// based on a given yml file.
func Update(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "update [path]",
		Short: "Update a part in the Software Parts Catalog",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli update <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argImportPath := args[0]
			if argImportPath == "" {
				return errors.New("error updating part, update subcommand usage: ./ccli update <Path>")
			}
			if argImportPath != "" {
				if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
					return errors.New("error importing part, import path not a yaml file")
				}
				f, err := os.Open(argImportPath)
				if err != nil {
					return errors.Wrapf(err, "error opening file")
				}
				defer f.Close()
				data, err := io.ReadAll(f)
				if err != nil {
					return errors.Wrapf(err, "error reading file")
				}
				var partData yaml.Part
				if err = yaml.Unmarshal(data, &partData); err != nil {
					return errors.Wrapf(err, "error decoding file contents")
				}
				slog.Debug("updating part")
				returnPart, err := graphql.UpdatePart(context.Background(), client, &partData)
				if err != nil {
					return errors.Wrapf(err, "error updating part")
				}
				prettyJson, err := json.MarshalIndent(&returnPart, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("Part successfully updated\n%s\n", string(prettyJson))
			}
			return nil
		},
	}
	return uploadCmd

}
