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
	uploadCmd := &cobra.Command{
		Use:   "upload [path]",
		Short: "Upload an archive to the Software Parts Catalog",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli upload <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argPath := args[0]
			if argPath == "" {
				return errors.New("error executing upload, upload subcommand usage: ccli upload <Path>")
			}
			if argPath != "" {
				slog.Debug("uploading file to server")
				response, err := graphql.UploadFile(http.DefaultClient, configFile.ServerAddr, argPath, "")
				if err != nil {
					return errors.Wrapf(err, "error uploading archive")
				}
				if response.Data != nil {
					fmt.Printf("Successfully uploaded package: %s\n", argPath)
				}
			}
			return nil
		},
	}
	return uploadCmd

}
