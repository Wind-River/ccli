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
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Find() handles the command for getting a part based on various aspects.
func Find(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for file
	findCmd := &cobra.Command{
		Use:   "find",
		Short: "Find a part from the Software Parts Catalog based on the find parameters like fvc, sha256, part query, part id.",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the find parameter. For more info run help")
		},
	}
	// add sub commands to find
	findCmd.AddCommand(FindPart(configFile, client, indent))
	findCmd.AddCommand(FindId(configFile, client, indent))
	findCmd.AddCommand(FindSha(configFile, client))
	findCmd.AddCommand(FindFvc(configFile, client))
	findCmd.AddCommand(FindProfile(configFile, client, indent))
	return findCmd
}

// FindPart() handles finding a part based on a search query/part name
func FindPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	findPartCmd := &cobra.Command{
		Use:   "part [search query]",
		Short: "Find a part using the name(i.e. search query)",
		// function to be run as a setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No part name or search query provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argSearchQuery := args[0]
			if argSearchQuery != "" {
				// search the catalog for the part using the search query
				slog.Debug("executing part search", slog.String("Query", argSearchQuery))
				response, err := graphql.Search(context.Background(), client, argSearchQuery)
				if err != nil {
					return errors.Wrapf(err, "error searching for part")
				}
				// marshal the response data into a json
				prettyJson, err := json.MarshalIndent(response, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("Result: %s\n", string(prettyJson))
			}
			return nil
		},
	}
	return findPartCmd
}

// FindId() handles finding a part based on part id
func FindId(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for find using id
	findIdCmd := &cobra.Command{
		Use:   "id [part id]",
		Short: "Find a part using the part id",
		// function to be run as a setup before command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No part id provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argPartID := args[0]
			if argPartID != "" {
				// find the part using the id
				slog.Debug("retrieving part by id", slog.String("ID", argPartID))
				response, err := graphql.GetPartByID(context.Background(), client, argPartID)
				if err != nil {
					return errors.Wrapf(err, "error getting part by id")
				}
				// marshal the response struct to a json
				prettyJson, err := json.MarshalIndent(&response, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("%s\n", string(prettyJson))
			}
			return nil
		},
	}
	return findIdCmd
}

// FindSha() handles finding a part id based on part sha256
func FindSha(configFile *config.ConfigData, client *graph.Client) *cobra.Command {
	// cobra command for find using sha
	findShaCmd := &cobra.Command{
		Use:   "sha256 [sha256]",
		Short: "Find a part using the Sha256",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No sha256 provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argSHA256 := args[0]
			if argSHA256 != "" {
				// get the part id using sha256
				slog.Debug("retrieving part id by sha256", slog.String("SHA256", argSHA256))
				partID, err := graphql.GetPartIDBySha256(context.Background(), client, argSHA256)
				if err != nil {
					return errors.Wrapf(err, "error retrieving part id")
				}
				fmt.Printf("Part ID: %s \n", partID.String())
			}
			return nil
		},
	}
	return findShaCmd
}

// FindFvc() handles finding a part id based on part file verification code
func FindFvc(configFile *config.ConfigData, client *graph.Client) *cobra.Command {
	// cobra command for find using fvc
	findFvcCmd := &cobra.Command{
		Use:   "fvc [fvc]",
		Short: "Find a part using the file verification code",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No file verification code provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argFVC := args[0]
			if argFVC != "" {
				// get the part id using the file verification code
				slog.Debug("retrieving part id by file verification code", slog.String("File Verification Code", argFVC))
				partID, err := graphql.GetPartIDByFVC(context.Background(), client, argFVC)
				if err != nil {
					return errors.Wrapf(err, "error retrieving part id")
				}
				fmt.Printf("Part ID: %s \n", partID.String())
			}
			return nil
		},
	}
	return findFvcCmd
}

// FindProfile() handles finding a specific type of part profile
// using its part id.
func FindProfile(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for finding profile
	findProfileCmd := &cobra.Command{
		Use:   "profile [profile type] [part id]",
		Short: "Find the specific profile of a given part using its profile type and part id",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("No profile type and part id provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argProfileType := args[0]
			argPartID := args[1]
			if argProfileType != "" {
				if argPartID == "" {
					return errors.New("error getting profile, missing part id")
				}
				// get the particular profile for a part based on profile type and part id
				slog.Debug("retrieving profile", slog.String("ID", argPartID), slog.String("Key", argProfileType))
				profile, err := graphql.GetProfile(context.Background(), client, argPartID, argProfileType)
				if err != nil {
					return errors.Wrapf(err, "error retrieving profile")
				}
				// check if any profiles were found
				if len(*profile) < 1 {
					fmt.Println("No documents found")
					os.Exit(0)
				}
				// marshal the profile data into a json
				prettyJson, err := json.MarshalIndent(&profile, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("%s\n", string(prettyJson))
				os.Exit(0)
			}
			return nil
		},
	}
	return findProfileCmd
}
