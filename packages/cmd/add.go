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

	"github.com/pkg/errors"

	graph "github.com/hasura/go-graphql-client"

	"github.com/spf13/cobra"
)

// Add() handles uploading a logical part or a part profile
// to the catalog from a yml file
func Add(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// create a cobra command
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a specific component(part or profile) to Software Parts Catalog.",
		// the function to be executed when the command is ran
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the component type to be added(i.e. part or profile). For more info run help")
		},
	}
	// attach sub commands
	addCmd.AddCommand(AddPart(configFile, client, indent))
	addCmd.AddCommand(AddProfile(configFile, client, indent))
	return addCmd
}

// AddPart() handles the sub command for uploading a logical
// part using the path to a yml file.
func AddPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	addPartCmd := &cobra.Command{
		Use:   "part [path]",
		Short: "Add a part to Software Parts Catalog.",
		// the function to be executed as a setup to the command being ran
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.")
			}
			return nil
		},
		// the function be executed when the command is ran
		RunE: func(cmd *cobra.Command, args []string) error {
			// scan the arguments to the command
			argPartImportPath := args[0]
			if argPartImportPath != "" {
				// check if the file is of yaml/yml format
				if argPartImportPath[len(argPartImportPath)-5:] != ".yaml" && argPartImportPath[len(argPartImportPath)-4:] != ".yml" {
					return errors.New("error importing part, import path not a yaml file")
				}
				// open the file
				f, err := os.Open(argPartImportPath)
				if err != nil {
					return errors.Wrap(err, "error opening file")
				}
				defer f.Close()
				slog.Debug("successfully opened file", slog.String("file:", argPartImportPath))
				// read the contents of the file
				data, err := io.ReadAll(f)
				if err != nil {
					return errors.Wrapf(err, "error reading file")
				}
				var partData yaml.Part
				// unmarshal all the data of the file into a struct
				if err = yaml.Unmarshal(data, &partData); err != nil {
					return errors.Wrapf(err, "error unmarshaling file contents")
				}
				slog.Debug("adding part")
				// call the graphql helper for adding a new part
				createdPart, err := graphql.AddPart(context.Background(), client, partData)
				if err != nil {
					return errors.Wrapf(err, "error adding part")
				}
				// create a indented json output for the path added
				prettyPart, err := json.MarshalIndent(&createdPart, "", indent)
				if err != nil {
					return errors.Wrapf(err, "error prettifying json")
				}
				fmt.Printf("Successfully added part from: %s\n", argPartImportPath)
				fmt.Printf("%s\n", string(prettyPart))
			}

			return nil
		},
	}
	return addPartCmd
}

// AddProfile() handles the upload of a part's profile
// like license, security and quality using a yml file
func AddProfile(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// create a cobra command
	addProfileCmd := &cobra.Command{
		Use:   "profile [path]",
		Short: "Add a profile to a part in Software Parts Catalog.",
		// function to be run as a setup for the command
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.")
			}
			return nil
		},
		// function to be run on command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argImportPath := args[0]
			if argImportPath != "" {
				// check if the file is of yaml/yml format
				if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
					return errors.New("error importing profile, import path not a yaml file")
				}
				// open the file
				f, err := os.Open(argImportPath)
				if err != nil {
					return errors.Wrapf(err, "error opening file")
				}
				defer f.Close()
				slog.Debug("Successfully opened file", slog.String("file:", argImportPath))
				// read the file data
				data, err := io.ReadAll(f)
				if err != nil {
					return errors.Wrapf(err, "error reading file")
				}
				var profileData yaml.Profile
				// unmarshal the file data into a struct
				if err = yaml.Unmarshal(data, &profileData); err != nil {
					return errors.Wrapf(err, "error unmarshaling file contents")
				}
				slog.Debug("adding profile", slog.String("Key", profileData.Profile))
				// switch case for various profile types
				switch profileData.Profile {
				case "security":
					var securityProfile yaml.SecurityProfile
					// unmarshal the data into a security profile struct
					if err = yaml.Unmarshal(data, &securityProfile); err != nil {
						return errors.Wrapf(err, "error unmarshaling security profile")
					}
					// marshal the security profile struct into a json
					jsonSecurityProfile, err := json.Marshal(securityProfile)
					if err != nil {
						return errors.Wrapf(err, "error marshaling json")
					}
					// check if the part identifier is present
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						return errors.New("error adding profile, no part identifier given")
					}
					// add the profile if the part id is given
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonSecurityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					// add the profile by first getting the part id using the fvc
					if profileData.FVC != "" {
						slog.Debug("retrieving part id by file verification code", slog.String("File Verification Code", profileData.FVC))
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					// add the profile by first getting the part id using the sha256
					if profileData.Sha256 != "" {
						slog.Debug("retrieving part id by sha256", slog.String("SHA256", profileData.Sha256))
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
					}
				case "licensing":
					var licensingProfile yaml.LicensingProfile
					// unmarshal the data into a licensing profile struct
					if err = yaml.Unmarshal(data, &licensingProfile); err != nil {
						return errors.Wrapf(err, "error unmarshaling licensing profile")
					}
					// marshal the licensing profile struct to json
					jsonLicensingProfile, err := json.Marshal(licensingProfile)
					if err != nil {
						return errors.Wrapf(err, "error marshaling json")
					}
					// check if the part identifier is given
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						return errors.New("error adding profile, no part identifier given")
					}
					// add the profile using the part id
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonLicensingProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					// add the profile by first getting the part id using the fvc
					if profileData.FVC != "" {
						slog.Debug("retrieving part id by file verification code", slog.String("File Verification Code", profileData.FVC))
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					// add the profile by first getting the part id using the sha256
					if profileData.Sha256 != "" {
						slog.Debug("retrieving part id by sha256", slog.String("SHA256", profileData.Sha256))
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
					}
				case "quality":
					var qualityProfile yaml.QualityProfile
					// unmarshal the data into a quality profile struct
					if err = yaml.Unmarshal(data, &qualityProfile); err != nil {
						return errors.Wrapf(err, "error unmarshaling quality profile")
					}
					// marshal the quality profile struct to a json
					jsonQualityProfile, err := json.Marshal(qualityProfile)
					if err != nil {
						return errors.Wrapf(err, "error marshaling json")
					}
					// cehck if the part identifier is given
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						return errors.Wrapf(err, "error adding profile, no part identifier given")
					}
					// add the profile by using the part id
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonQualityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					// add the profile by first getting the part id using the fvc
					if profileData.FVC != "" {
						slog.Debug("retrieving part id by file verification code", slog.String("File Verification Code", profileData.FVC))
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					// add the profile by first getting the part id using the sha256
					if profileData.Sha256 != "" {
						slog.Debug("retrieving part id by sha256", slog.String("SHA256", profileData.Sha256))
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							return errors.Wrapf(err, "error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
							return errors.Wrapf(err, "error adding profile")
						}
						fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
					}
				}
			}

			return nil
		},
	}
	return addProfileCmd
}
