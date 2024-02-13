package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Find(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	findCmd := &cobra.Command{
		Use:   "find",
		Short: "Find a part from the Software Parts Catalog based on the find parameters like fvc, sha256, part query, part id.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide the find parameter. For more info run help")
		},
	}
	findCmd.AddCommand(FindPart(configFile, client, indent))
	findCmd.AddCommand(FindId(configFile, client, indent))
	findCmd.AddCommand(FindSha(configFile, client))
	findCmd.AddCommand(FindFvc(configFile, client))
	findCmd.AddCommand(FindProfile(configFile, client, indent))
	return findCmd
}

func FindPart(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	findPartCmd := &cobra.Command{
		Use:   "part",
		Short: "Find a part using the name(i.e. search query)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No part name or search query provided.\nUsage: ccli find part <part name/search query>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argSearchQuery := args[0]
			if argSearchQuery != "" {
				log.Debug().Str("Query", argSearchQuery).Msg("executing part search")
				response, err := graphql.Search(context.Background(), client, argSearchQuery)
				if err != nil {
					log.Fatal().Err(err).Msg("error searching for part")
				}

				prettyJson, err := json.MarshalIndent(response, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Printf("Result: %s\n", string(prettyJson))
			}
			return nil
		},
	}
	return findPartCmd
}

func FindId(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	findIdCmd := &cobra.Command{
		Use:   "id",
		Short: "Find a part using the part id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No part id provided.\nUsage: ccli find id <part id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argPartID := args[0]
			if argPartID != "" {
				log.Debug().Str("ID", argPartID).Msg("retrieving part by id")
				response, err := graphql.GetPartByID(context.Background(), client, argPartID)
				if err != nil {
					log.Fatal().Err(err).Msg("error getting part by id")
				}
				prettyJson, err := json.MarshalIndent(&response, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Printf("%s\n", string(prettyJson))
			}
			return nil
		},
	}
	return findIdCmd
}
func FindSha(configFile *config.ConfigData, client *graph.Client) *cobra.Command {
	findShaCmd := &cobra.Command{
		Use:   "sha256",
		Short: "Find a part using the Sha256",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No sha256 provided.\nUsage: ccli find sha256 <sha256>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argSHA256 := args[0]
			if argSHA256 != "" {
				log.Debug().Str("SHA256", argSHA256).Msg("retrieving part id by sha256")
				partID, err := graphql.GetPartIDBySha256(context.Background(), client, argSHA256)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving part id")
				}
				fmt.Printf("Part ID: %s \n", partID.String())
			}
			return nil
		},
	}
	return findShaCmd
}
func FindFvc(configFile *config.ConfigData, client *graph.Client) *cobra.Command {
	findFvcCmd := &cobra.Command{
		Use:   "fvc",
		Short: "Find a part using the file verification code",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No file verification code provided.\nUsage: ccli find fvc <file verification code>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argFVC := args[0]
			if argFVC != "" {
				log.Debug().Str("File Verification Code", argFVC).Msg("retrieving part id by file verification code")
				partID, err := graphql.GetPartIDByFVC(context.Background(), client, argFVC)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving part id")
				}
				fmt.Printf("Part ID: %s \n", partID.String())
			}
			return nil
		},
	}
	return findFvcCmd
}
func FindProfile(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	findProfileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Find the specific profile of a given part using its profile type and part id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("No profile type and part id provided.\nUsage: ccli find profile <profile type> <part id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argProfileType := args[0]
			argPartID := args[1]
			if argProfileType != "" {
				if argPartID == "" {
					log.Fatal().Msg("error getting profile, missing part id")
				}
				log.Debug().Str("ID", argPartID).Str("Key", argProfileType).Msg("retrieving profile")
				profile, err := graphql.GetProfile(context.Background(), client, argPartID, argProfileType)
				if err != nil {
					log.Fatal().Err(err).Msg("error retrieving profile")
				}
				if len(*profile) < 1 {
					fmt.Println("No documents found")
					os.Exit(0)
				}
				prettyJson, err := json.MarshalIndent(&profile, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Printf("%s\n", string(prettyJson))
				os.Exit(0)
			}
			return nil
		},
	}
	return findProfileCmd
}
