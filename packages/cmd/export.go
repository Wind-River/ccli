package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/yaml"

	graph "github.com/hasura/go-graphql-client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Export(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	var output string
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export a component based on the subcommands to a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the export subcommand(part or template). For more info run help")
		},
	}
	exportCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Path to the output file")
	exportCmd.AddCommand(ExportPart(configFile, client, indent))
	exportCmd.AddCommand(ExportTemplate(configFile, client, indent))
	return exportCmd
}

func ExportPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportPartCmd := &cobra.Command{
		Use:   "part",
		Short: "Export a part from the Software Parts Catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the export part subcommand(i.e. id, sha256, fvc). For more info run help")
		},
	}
	exportPartCmd.AddCommand(ExportPartId(configFile, client, indent))
	exportPartCmd.AddCommand(ExportPartSha(configFile, client, indent))
	exportPartCmd.AddCommand(ExportPartFvc(configFile, client, indent))

	return exportPartCmd
}

func ExportPartId(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportPartIdCmd := &cobra.Command{
		Use:   "id",
		Short: "Export a part using part id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part id provided.\nUsage: ccli export part id <part id> -o <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argPartID := args[0]
			if argPartID != "" {
				log.Debug().Str("ID", argPartID).Msg("retrieving part by id")
				part, err := graphql.GetPartByID(context.Background(), client, argPartID)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving part")
				}
				ExportHelper(part, argExportPath)
			}
			return nil
		},
	}

	return exportPartIdCmd
}

func ExportPartSha(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportPartShaCmd := &cobra.Command{
		Use:   "sha256",
		Short: "Export a part using the Sha256",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part sha256 provided.\nUsage: ccli export part sha256 <sha256> -o <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argSHA256 := args[0]
			if argSHA256 != "" {
				log.Debug().Str("SHA256", argSHA256).Msg("retrieving part by sha256")
				part, err := graphql.GetPartBySHA256(context.Background(), client, argSHA256)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving part")
				}
				ExportHelper(part, argExportPath)
			}

			return nil
		},
	}
	return exportPartShaCmd
}

func ExportPartFvc(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportPartFvcCmd := &cobra.Command{
		Use:   "fvc",
		Short: "Export a part using fvc",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part file verification code provided.\nUsage: ccli export part fvc <fvc> -o <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			argFVC := args[0]
			if argFVC != "" {
				log.Debug().Str("File Verification Code", argFVC).Msg("retrieving part by file verification code")
				part, err := graphql.GetPartByFVC(context.Background(), client, argFVC)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving part")
				}
				ExportHelper(part, argExportPath)
			}
			return nil
		},
	}
	return exportPartFvcCmd
}

func ExportTemplate(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateCmd := &cobra.Command{
		Use:   "template",
		Short: "Export a template to a given file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide a the find parameter. For more info run help")
		},
	}
	exportTemplateCmd.AddCommand(ExportTemplatePart(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateSecurity(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateQuality(configFile, client, indent))
	exportTemplateCmd.AddCommand(ExportTemplateLicense(configFile, client, indent))

	return exportTemplateCmd
}

func ExportTemplatePart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplatePartCmd := &cobra.Command{
		Use:   "part",
		Short: "Export a part template",
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			yamlPart := new(yaml.Part)
			yamlPart.Format = 1.0
			f, err := os.Create(argExportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating template file")
			}
			defer f.Close()
			yamlPartTemplate, err := yaml.Marshal(&yamlPart)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling part template")
			}
			_, err = f.Write(yamlPartTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			fmt.Printf("Part template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplatePartCmd
}

func ExportTemplateSecurity(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateSecurityCmd := &cobra.Command{
		Use:   "security",
		Short: "Export a security template",
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			yamlProfile := new(yaml.Profile)
			yamlSecurityProfile := new(yaml.SecurityProfile)
			yamlCVE := new(yaml.CVE)
			yamlProfile.Format = 1.0
			yamlSecurityProfile.CVEList = append(yamlSecurityProfile.CVEList, *yamlCVE)
			f, err := os.Create(argExportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating template file")
			}
			defer f.Close()
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling profile template")
			}
			yamlSecurityProfileTemplate, err := yaml.Marshal(&yamlSecurityProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling security profile template")
			}
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			_, err = f.Write(yamlSecurityProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateSecurityCmd
}

func ExportTemplateQuality(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateQualityCmd := &cobra.Command{
		Use:   "quality",
		Short: "Export a quality template",
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			yamlProfile := new(yaml.Profile)
			yamlQualityProfile := new(yaml.QualityProfile)
			yamlBug := new(yaml.Bug)
			yamlProfile.Format = 1.0
			yamlQualityProfile.BugList = append(yamlQualityProfile.BugList, *yamlBug)
			f, err := os.Create(argExportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating template file")
			}
			defer f.Close()
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling profile template")
			}
			yamlQualityProfileTemplate, err := yaml.Marshal(&yamlQualityProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling quality profile template")
			}
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			_, err = f.Write(yamlQualityProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateQualityCmd
}

func ExportTemplateLicense(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	exportTemplateLicenseCmd := &cobra.Command{
		Use:   "license",
		Short: "Export a license template",
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argExportPath, _ := cmd.Flags().GetString("output")
			if argExportPath == "" {
				return errors.New("Output path for exporting is not provided")
			}
			yamlProfile := new(yaml.Profile)
			yamlLicensingProfile := new(yaml.LicensingProfile)
			yamlLicense := new(yaml.License)
			yamlProfile.Format = 1.0
			yamlLicensingProfile.LicenseAnalysis = append(yamlLicensingProfile.LicenseAnalysis, *yamlLicense)
			f, err := os.Create(argExportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating template file")
			}
			defer f.Close()
			yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling profile template")
			}
			yamlLicensingProfileTemplate, err := yaml.Marshal(&yamlLicensingProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("error marshaling licensing profile template")
			}
			_, err = f.Write(yamlProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			_, err = f.Write(yamlLicensingProfileTemplate)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing template to file")
			}
			fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
			return nil
		},
	}
	return exportTemplateLicenseCmd
}

func ExportHelper(part *graphql.Part, argExportPath string) {
	var yamlPart yaml.Part
	if err := graphql.UnmarshalPart(part, &yamlPart); err != nil {
		log.Fatal().Err(err).Msg("error parsing part into yaml")
	}

	ret, err := yaml.Marshal(yamlPart)
	if err != nil {
		log.Fatal().Err(err).Msg("error marshalling yaml")
	}

	yamlFile, err := os.Create(argExportPath)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating yaml file")
	}
	defer yamlFile.Close()

	_, err = yamlFile.Write(ret)
	if err != nil {
		log.Fatal().Err(err).Msg("error writing part to yaml file")
	}
	fmt.Printf("Part successfully exported to path: %s\n", argExportPath)
}
