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
	rootCmd := &cobra.Command{
		Use:   "ccli",
		Short: "Ccli is used to interact with the Software Parts Catalog.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verboseFlag {
				slogOptions := new(slog.HandlerOptions)
				slogOptions.Level = slog.LevelDebug
				if configFile.LogLevel == 2 {
					slogOptions.AddSource = true
				}
				slog.SetDefault(slog.New(slog.NewJSONHandler(logWriter, slogOptions)))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Please provide a sub-command to be executed. Refer to the examples by running ccli examples or use help for more information.")

		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "To Execute commands in verbose mode")
	return rootCmd
}
