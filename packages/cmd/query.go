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
	"log/slog"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Query() handles the execution of a given graphql query
func Query(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for graphql query
	return &cobra.Command{
		Use:   "query [graphql query]",
		Short: "Query the Software Parts Catalog",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No query provided.")
			}
			return nil
		},
		//function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argQuery := args[0]
			if argQuery == "" {
				return errors.New("error executing user query, query subcommand usage: ccli query <GraphQL Query>")
			}
			// check if the query is not nil and run the graphql query
			if argQuery != "" {
				slog.Debug("executing raw graphql query")
				response, err := graphql.Query(context.Background(), client, argQuery)
				if err != nil {
					return errors.Wrapf(err, "error querying graphql")
				}

				// json result will be output in prettified format
				var data map[string]interface{}
				json.Unmarshal(response, &data)
				// marshal the response into a json
				prettyJson, err := json.MarshalIndent(data, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Println(string(prettyJson))
			}
			return nil
		},
	}
}
