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
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/yaml"

	graph "github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Export() handles getting a part or a template and
// saving it out to a file on the given path
func Export(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	var output string
	// cobra command for export
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export a component based on the subcommands to a file",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the export subcommand(part or template). For more info run help")
		},
	}
	// add a persistent flag for output file
	exportCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Path to the output file")
	// add subcommands for export
	exportCmd.AddCommand(ExportPart(configFile, client, indent))
	exportCmd.AddCommand(ExportTemplate(configFile, client, indent))
	return exportCmd
}

// ExportPart() is a sub command and handles the download
// of part data and writing to a file on a given path
func ExportPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for exporting part
	exportPartCmd := &cobra.Command{
		Use:   "part",
		Short: "Export a part from the Software Parts Catalog",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the export part subcommand(i.e. id, sha256, fvc). For more info run help")
		},
	}
	// add sub commands for part export based on search parameter
	exportPartCmd.AddCommand(ExportPartId(configFile, client, indent))
	exportPartCmd.AddCommand(ExportPartSha(configFile, client, indent))
	exportPartCmd.AddCommand(ExportPartFvc(configFile, client, indent))

	return exportPartCmd
}

// ExportPartId() gets the part based on part id
func ExportPartId(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for exporting part based on id
	exportPartIdCmd := &cobra.Command{
		Use:   "id [part id] [-o] [export path]",
		Short: "Export a part using part id",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part id provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argPartID := args[0]
			// get the part data using the part id
			if argPartID != "" {
				slog.Debug("retrieving part by id", slog.String("ID", argPartID))
				part, err := graphql.GetPartByID(context.Background(), client, argPartID)
				if err != nil {
					return errors.Wrapf(err, "error retrieving part")
				}
				// export the part data into a file on the given path
				err = ExportHelper(part, argExportPath)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	return exportPartIdCmd
}

// ExportPartSha() gets the part based on part sha256
func ExportPartSha(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for export using sha256
	exportPartShaCmd := &cobra.Command{
		Use:   "sha256 [sha256] [-o] [export path]",
		Short: "Export a part using the Sha256",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part sha256 provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argSHA256 := args[0]
			if argSHA256 != "" {
				// get the part data using sha256
				slog.Debug("retrieving part by sha256", slog.String("SHA256", argSHA256))
				part, err := graphql.GetPartBySHA256(context.Background(), client, argSHA256)
				if err != nil {
					return errors.Wrapf(err, "error retrieving part")
				}
				// export the part data into a file at the given path
				err = ExportHelper(part, argExportPath)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	return exportPartShaCmd
}

// ExportPartFvc() gets the part based on part file verification code
func ExportPartFvc(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for export using fvc
	exportPartFvcCmd := &cobra.Command{
		Use:   "fvc [fvc] [-o] [export path]",
		Short: "Export a part using fvc",
		// function to be run as setup for command execution
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part file verification code provided.")
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argFVC := args[0]
			if argFVC != "" {
				// get the part data using file verification code
				slog.Debug("retrieving part by file verification code", slog.String("File Verification Code", argFVC))
				part, err := graphql.GetPartByFVC(context.Background(), client, argFVC)
				if err != nil {
					return errors.Wrapf(err, "error retrieving part")
				}
				// export the part data into a file at the given path
				err = ExportHelper(part, argExportPath)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	return exportPartFvcCmd
}

// ExportTemplate() handles getting out a template for various part/profile
// data into a file on the given path
func ExportTemplate(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for exporting template
	exportTemplateCmd := &cobra.Command{
		Use:   "template [-o] [export path]",
		Short: "Export a template to a given file",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide a the find parameter. For more info run help")
		},
	}
	// add sub commands for template export
	exportTemplateCmd.AddCommand(ExportTemplatePart(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateSecurity(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateQuality(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateLicense(configFile, client, indent))

	return exportTemplateCmd
}

// ExportTemplatePart() handles the template for a part
func ExportTemplatePart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	// cobra command for exporting part template
	exportTemplatePartCmd := &cobra.Command{
		Use:   "part [-o] [export path]",
		Short: "Export a part template",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			// create a new part template
			yamlPart := new(yaml.Part)
			yamlPart.Format = 1.0
			// create a new file at the given path
			f, err := os.Create(argExportPath)
			if err != nil {
				return errors.Wrapf(err, "error creating template file")
			}
			defer f.Close()
			// marshal the part template into yaml
			yamlPartTemplate, err := yaml.Marshal(&yamlPart)
			if err != nil {
				return errors.Wrapf(err, "error marshaling part template")
			}
			// write the part template into the file
			_, err = f.Write(yamlPartTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			fmt.Printf("Part template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplatePartCmd
}

// ExportTemplateSecurity() handles the template for a security profile
func ExportTemplateSecurity(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	//cobra command for exporting security template
	exportTemplateSecurityCmd := &cobra.Command{
		Use:   "security [-o] [export path]",
		Short: "Export a security template",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			// create a new profile template
			yamlProfile := new(yaml.Profile)
			yamlSecurityProfile := new(yaml.SecurityProfile)
			yamlCVE := new(yaml.CVE)
			yamlProfile.Format = 1.0
			// create a new security profile template
			yamlSecurityProfile.CVEList = append(yamlSecurityProfile.CVEList, *yamlCVE)
			// create the file at the given path
			f, err := os.Create(argExportPath)
			if err != nil {
				return errors.Wrapf(err, "error creating template file")
			}
			defer f.Close()
			// marshal the profile template to yaml
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling profile template")
			}
			//marshal the security profile to yaml
			yamlSecurityProfileTemplate, err := yaml.Marshal(&yamlSecurityProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling security profile template")
			}
			// write the profile template to the file
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			// write the security profile template to the file
			_, err = f.Write(yamlSecurityProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateSecurityCmd
}

// ExportTemplateQuality() handles the template for a quality profile
func ExportTemplateQuality(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateQualityCmd := &cobra.Command{
		Use:   "quality [-o] [export path]",
		Short: "Export a quality template",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			// create a new profile template
			yamlProfile := new(yaml.Profile)
			yamlQualityProfile := new(yaml.QualityProfile)
			yamlBug := new(yaml.Bug)
			yamlProfile.Format = 1.0
			// create a new quality profile template
			yamlQualityProfile.BugList = append(yamlQualityProfile.BugList, *yamlBug)
			// create the file at the given path
			f, err := os.Create(argExportPath)
			if err != nil {
				return errors.Wrapf(err, "error creating template file")
			}
			defer f.Close()
			//marshal the  profile to yaml
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling profile template")
			}
			//marshal the quality profile to yaml
			yamlQualityProfileTemplate, err := yaml.Marshal(&yamlQualityProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling quality profile template")
			}
			// write the profile template to the file
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			// write the quality profile template to the file
			_, err = f.Write(yamlQualityProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateQualityCmd
}

// ExportTemplateLicense() handles the template for a license profile
func ExportTemplateLicense(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateLicenseCmd := &cobra.Command{
		Use:   "license [-o] [export path]",
		Short: "Export a license template",
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			// create a new profile template
			yamlProfile := new(yaml.Profile)
			yamlLicensingProfile := new(yaml.LicensingProfile)
			yamlLicense := new(yaml.License)
			yamlProfile.Format = 1.0
			// create a new license profile template
			yamlLicensingProfile.LicenseAnalysis = append(yamlLicensingProfile.LicenseAnalysis, *yamlLicense)
			// create the file at the given path
			f, err := os.Create(argExportPath)
			if err != nil {
				return errors.Wrapf(err, "error creating template file")
			}
			defer f.Close()
			//marshal the profile to yaml
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling profile template")
			}
			//marshal the license profile to yaml
			yamlLicensingProfileTemplate, err := yaml.Marshal(&yamlLicensingProfile)
			if err != nil {
				return errors.Wrapf(err, "error marshaling licensing profile template")
			}
			// write the profile template to the file
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			// write the license profile template to the file
			_, err = f.Write(yamlLicensingProfileTemplate)
			if err != nil {
				return errors.Wrapf(err, "error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateLicenseCmd
}

// ExportHelper() is a helper function for storing the data
// into a file on a given path.
func ExportHelper(part *graphql.Part, argExportPath string) error {
	var yamlPart yaml.Part
	// unmarshal the part data into yaml
	if err := graphql.UnmarshalPart(part, &yamlPart); err != nil {
		return errors.Wrapf(err, "error parsing part into yaml")
	}
	// marshal the yaml
	ret, err := yaml.Marshal(yamlPart)
	if err != nil {
		return errors.Wrapf(err, "error marshalling yaml")
	}
	// create the file on the given path
	yamlFile, err := os.Create(argExportPath)
	if err != nil {
		return errors.Wrapf(err, "error creating yaml file")
	}
	defer yamlFile.Close()
	// write the data to the file
	_, err = yamlFile.Write(ret)
	if err != nil {
		return errors.Wrapf(err, "error writing part to yaml file")
	}
	fmt.Printf("Part successfully exported to path: %s\n", argExportPath)
	return nil
}
