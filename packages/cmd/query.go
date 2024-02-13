package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Query(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Query the Software Parts Catalog",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No query provided.\nUsage: ccli query <graphql query>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argQuery := args[0]
			if argQuery == "" {
				log.Fatal().Msg("error executing user query, query subcommand usage: ccli query <GraphQL Query>")
			}
			if argQuery != "" {
				log.Debug().Msg("executing raw graphql query")
				response, err := graphql.Query(context.Background(), client, argQuery)
				if err != nil {
					log.Fatal().Err(err).Msg("error querying graphql")
				}

				// json result will be output in prettified format
				var data map[string]interface{}
				json.Unmarshal(response, &data)

				prettyJson, err := json.MarshalIndent(data, "", indent)
				if err != nil {
					log.Fatal().Err(err).Msg("error prettifying json")
				}
				fmt.Println(string(prettyJson))
			}
			return nil
		},
	}
}
