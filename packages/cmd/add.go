package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/yaml"

	graph "github.com/hasura/go-graphql-client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Add(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a specific component(part or profile) to Software Parts Catalog.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the component type to be added(i.e. part or profile). For more info run help")
		},
	}
	addCmd.AddCommand(AddPart(configFile, client, indent))
	addCmd.AddCommand(AddProfile(configFile, client, indent))
	return addCmd
}
func AddPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	addPartCmd := &cobra.Command{
		Use:   "part",
		Short: "Add a part to Software Parts Catalog.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli add part <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argPartImportPath := args[0]
			if argPartImportPath != "" {
				if argPartImportPath[len(argPartImportPath)-5:] != ".yaml" && argPartImportPath[len(argPartImportPath)-4:] != ".yml" {
					log.Fatal().Msg("error importing part, import path not a yaml file")
				}
				f, err := os.Open(argPartImportPath)
				if err != nil {
					log.Fatal().Err(err).Msg("error opening file")
				}
				defer f.Close()
				log.Debug().Msgf("successfully opened file at: %s", argPartImportPath)
				data, err := io.ReadAll(f)
				if err != nil {
					log.Fatal().Err(err).Msg("error reading file")
				}
				var partData yaml.Part
				if err = yaml.Unmarshal(data, &partData); err != nil {
					log.Fatal().Err(err).Msg("error unmarshaling file contents")
				}
				log.Debug().Msg("adding part")
				createdPart, err := graphql.AddPart(context.Background(), client, partData)
				if err != nil {
					log.Fatal().Err(err).Msg("error adding part")
				}
				prettyPart, err := json.MarshalIndent(&createdPart, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Printf("Successfully added part from: %s\n", argPartImportPath)
				fmt.Printf("%s\n", string(prettyPart))
			}

			return nil
		},
	}
	return addPartCmd
}
func AddProfile(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	addProfileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Add a profile to a part in Software Parts Catalog.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli add profile <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argImportPath := args[0]
			if argImportPath != "" {
				if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
					log.Fatal().Msg("error importing profile, import path not a yaml file")
				}
				f, err := os.Open(argImportPath)
				if err != nil {
					log.Fatal().Err(err).Msg("error opening file")
				}
				defer f.Close()
				log.Debug().Msgf("successfully opened file at: %s", argImportPath)
				data, err := io.ReadAll(f)
				if err != nil {
					log.Fatal().Err(err).Msg("error reading file")
				}
				var profileData yaml.Profile
				if err = yaml.Unmarshal(data, &profileData); err != nil {
					log.Fatal().Err(err).Msg("error unmarshaling file contents")
				}
				log.Debug().Str("Key", profileData.Profile).Msg("adding profile")
				switch profileData.Profile {
				case "security":
					var securityProfile yaml.SecurityProfile
					if err = yaml.Unmarshal(data, &securityProfile); err != nil {
						log.Fatal().Err(err).Msg("error unmarshaling security profile")
					}
					jsonSecurityProfile, err := json.Marshal(securityProfile)
					if err != nil {
						log.Fatal().Err(err).Msg("error marshaling json")
					}
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						log.Fatal().Msg("error adding profile, no part identifier given")
					}
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonSecurityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					if profileData.FVC != "" {
						log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					if profileData.Sha256 != "" {
						log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
					}
				case "licensing":
					var licensingProfile yaml.LicensingProfile
					if err = yaml.Unmarshal(data, &licensingProfile); err != nil {
						log.Fatal().Err(err).Msg("error unmarshaling licensing profile")
					}
					jsonLicensingProfile, err := json.Marshal(licensingProfile)
					if err != nil {
						log.Fatal().Err(err).Msg("error marshaling json")
					}
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						log.Fatal().Msg("error adding profile, no part identifier given")
					}
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonLicensingProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					if profileData.FVC != "" {
						log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					if profileData.Sha256 != "" {
						log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
					}
				case "quality":
					var qualityProfile yaml.QualityProfile
					if err = yaml.Unmarshal(data, &qualityProfile); err != nil {
						log.Fatal().Err(err).Msg("error unmarshaling quality profile")
					}
					jsonQualityProfile, err := json.Marshal(qualityProfile)
					if err != nil {
						log.Fatal().Err(err).Msg("error marshaling json")
					}
					if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
						log.Fatal().Msg("error adding profile, no part identifier given")
					}
					if profileData.CatalogID != "" {
						if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonQualityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
					}
					if profileData.FVC != "" {
						log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
						uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by fvc")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
						}
						fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
						break
					}
					if profileData.Sha256 != "" {
						log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
						uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
						if err != nil {
							log.Fatal().Err(err).Msg("error retrieving part id by sha256")
						}
						if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
							log.Fatal().Err(err).Msg("error adding profile")
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
