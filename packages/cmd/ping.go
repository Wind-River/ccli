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
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/http"

	"log/slog"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Ping() makes a call to the catalog server and checks if the catalog server
// is responding and ready for further api calls.
func Ping(configFile *config.ConfigData) *cobra.Command {
	// cobra command for pinging the server
	return &cobra.Command{
		Use:   "ping",
		Short: "Ping the Catalog server for the current time",
		Args:  cobra.MinimumNArgs(0),
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if the server address is nil
			if configFile.ServerAddr == "" {
				return errors.New("invalid configuration file, no server address located")
			}
			slog.Debug("Pinging server", slog.String("Address", configFile.ServerAddr))
			// ping the server
			resp, err := http.DefaultClient.Get(configFile.ServerAddr)
			if err != nil {
				return errors.New("error contacting server")
			}
			resp.Body.Close()
			// check if the response's status code is valid
			if resp.StatusCode != 200 && resp.StatusCode != 422 {
				return errors.New("error reaching server, status code:" + fmt.Sprint(resp.StatusCode))
			} else {
				fmt.Println("Ping Result: Success")
			}
			return nil
		},
	}
}
