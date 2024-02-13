package cmd

import (
	"context"
	"errors"
	"fmt"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"

	graph "github.com/hasura/go-graphql-client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete(configFile *config.ConfigData, client *graph.Client, indent string) *cobra.Command {
	var argForcedMode bool
	var argRecursiveMode bool
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a part in the Software Parts Catalog.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No part id provided.\nUsage: ccli delete <part id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			argPartID := args[0]
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			if argPartID == "" {
				log.Fatal().Msg("error deleting part, delete subcommand usage: ./ccli delete <catalog_id>")
			}
			if argPartID != "" {
				log.Debug().Str("ID", argPartID).Msg("deleting part")
				if err := graphql.DeletePart(context.Background(), client, argPartID, argRecursiveMode, argForcedMode); err != nil {
					log.Fatal().Err(err).Msg("error deleting part from catalog")
				}
				fmt.Printf("Successfully deleted id: %s from catalog\n", argPartID)
			}
			return nil
		},
	}
	deleteCmd.PersistentFlags().BoolVarP(&argRecursiveMode, "recursive", "r", false, "To delete parts recursively")
	deleteCmd.PersistentFlags().BoolVarP(&argForcedMode, "force", "f", false, "To delete parts forcefully")
	return deleteCmd
}
