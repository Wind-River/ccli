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
	"log/slog"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Upload() uses the graphql upload library to upload
// an archive present at the given path.
func Upload(configFile *config.ConfigData) *cobra.Command {
	// cobra command for upload
	uploadCmd := &cobra.Command{
		Use:   "upload [path]",
		Short: "Upload an archive to the Software Parts Catalog",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argPath := args[0]
			if argPath == "" {
				return errors.New("error executing upload, upload subcommand usage: ccli upload <Path>")
			}
			// check if the file path is present and upload it to the catalog
			if argPath != "" {
				slog.Debug("uploading file to server")
				response, err := graphql.UploadFile(http.DefaultClient, configFile.ServerAddr, argPath, "")
				if err != nil {
					return errors.Wrapf(err, "error uploading archive")
				}
				// check if the response is present
				if response != nil {
					if len(response.Errors) > 0 {
						return errors.New(fmt.Sprintf("error uploading archive: %v", response.Errors))
					}
					if response.Data != nil {
						fmt.Printf("Successfully uploaded package: %s\n", argPath)
					}
				}
			}
			return nil
		},
	}
	return uploadCmd

}
