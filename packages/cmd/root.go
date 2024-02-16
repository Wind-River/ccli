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
	"wrs/catalog/ccli/packages/config"

	"github.com/pkg/errors"

	"log/slog"

	"github.com/spf13/cobra"
)

// RootCmd() is the root command which results in an error and a usage
// message advising the user to add sub commands
func RootCmd(configFile *config.ConfigData, logWriter *config.LogWriter) *cobra.Command {
	var verboseFlag bool
	// cobra command for root ccli
	rootCmd := &cobra.Command{
		Use:   "ccli",
		Short: "Ccli is used to interact with the Software Parts Catalog.",
		// function which is always to be reun before command execution
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// check if the verbose flag is on
			if verboseFlag {
				// create ne logger options and level
				slogOptions := new(slog.HandlerOptions)
				slogOptions.Level = slog.LevelDebug
				// if the log level is 2, add he source information to the options
				if configFile.LogLevel == 2 {
					slogOptions.AddSource = true
				}
				// set a new default for logging
				slog.SetDefault(slog.New(slog.NewJSONHandler(logWriter, slogOptions)))
			}
			return nil
		},
		// function to be run during command execution
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide a sub-command to be executed. Refer to the examples by running ccli examples or use help for more information.")

		},
	}
	// add a flag to the root command for having a verbose value
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "To Execute commands in verbose mode")
	return rootCmd
}
