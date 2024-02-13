package cmd

import (
	"errors"
	"fmt"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Ping(configFile *config.ConfigData) *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "Ping the Catalog server for the current time",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			argVerboseMode, _ := cmd.Flags().GetBool("verbose")
			if argVerboseMode {
				zerolog.SetGlobalLevel(0)
			}
			if configFile.ServerAddr == "" {
				return errors.New("invalid configuration file, no server address located")
			}
			log.Debug().Str("Address", configFile.ServerAddr).Msg("pinging server")
			resp, err := http.DefaultClient.Get(configFile.ServerAddr)
			if err != nil {
				return errors.New("error contacting server")
			}
			resp.Body.Close()
			if resp.StatusCode != 200 && resp.StatusCode != 422 {
				return errors.New("error reaching server, status code:" + fmt.Sprint(resp.StatusCode))
			} else {
				fmt.Println("Ping Result: Success")
			}
			return nil
		},
	}
}
