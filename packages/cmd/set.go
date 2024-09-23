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

// Set() is a sub command responsible for setting part information including zero values
// based on a given yml file.
func Set(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for set part
	setCmd := &cobra.Command{
		Use:   "set [path]",
		Short: "Set the fields of a part in the Software Parts Catalog including empty values",
		// function to be run as setup for the command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// check if exactly 1 argument is present
			if len(args) < 1 {
				return errors.New("No path provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argImportPath := args[0]
			if argImportPath == "" {
				return errors.New("error setting part fields, set subcommand usage: ./ccli set <Path>")
			}
			// check if the file is of yaml/yml format
			if argImportPath != "" {
				if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
					return errors.New("error importing part, import path not a yaml file")
				}
				// open the file
				f, err := os.Open(argImportPath)
				if err != nil {
					return errors.Wrapf(err, "error opening file")
				}
				defer f.Close()
				// read all the data from the file
				data, err := io.ReadAll(f)
				if err != nil {
					return errors.Wrapf(err, "error reading file")
				}
				// unmarshal the data of the file into a struct
				var partData yaml.Part
				if err = yaml.Unmarshal(data, &partData); err != nil {
					return errors.Wrapf(err, "error decoding file contents")
				}
				slog.Debug("setting part fields")
				// set the part fields with the given part data
				returnPart, err := graphql.SetPart(context.Background(), client, &partData)
				if err != nil {
					return errors.Wrapf(err, "error setting part fields")
				}
				// marshal the struct into a json
				prettyJson, err := json.MarshalIndent(&returnPart, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("Part fields successfully set\n%s\n", string(prettyJson))
			}
			return nil
		},
	}
	return setCmd

}
