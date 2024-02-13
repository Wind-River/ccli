package cmd

import (
	"errors"
	"fmt"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Upload(configFile *config.ConfigData) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an archive to the Software Parts Catalog",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("No path provided.\nUsage: ccli upload <path to file>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			argPath := args[0]
			if argPath == "" {
				log.Fatal().Msg("error executing upload, upload subcommand usage: ccli upload <Path>")
			}
			if argPath != "" {
				log.Debug().Msg("uploading file to server")
				response, err := graphql.UploadFile(http.DefaultClient, configFile.ServerAddr, argPath, "")
				if err != nil {
					log.Fatal().Err(err).Msg("error uploading archive")
				}
				if response.Data != nil {
					fmt.Printf("Successfully uploaded package: %s\n", argPath)
				}
			}
			return nil
		},
	}
	return uploadCmd

}
