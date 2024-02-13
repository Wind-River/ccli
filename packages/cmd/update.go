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

func Update(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a part in the Software Parts Catalog",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli update <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argImportPath := args[0]
			if argImportPath == "" {
				log.Fatal().Msg("error updating part, update subcommand usage: ./ccli update <Path>")
			}
			if argImportPath != "" {
				if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
					log.Fatal().Msg("error importing part, import path not a yaml file")
				}
				f, err := os.Open(argImportPath)
				if err != nil {
					log.Fatal().Err(err).Msg("error opening file")
				}
				defer f.Close()
				data, err := io.ReadAll(f)
				if err != nil {
					log.Fatal().Err(err).Msg("error reading file")
				}
				var partData yaml.Part
				if err = yaml.Unmarshal(data, &partData); err != nil {
					log.Fatal().Err(err).Msg("error decoding file contents")
				}
				log.Debug().Msg("updating part")
				returnPart, err := graphql.UpdatePart(context.Background(), client, &partData)
				if err != nil {
					log.Fatal().Err(err).Msg("error updating part")
				}
				prettyJson, err := json.MarshalIndent(&returnPart, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Printf("Part successfully updated\n%s\n", string(prettyJson))
			}
			return nil
		},
	}
	return uploadCmd

}
