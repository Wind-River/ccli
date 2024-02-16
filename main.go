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
package main

import (
	_ "embed"
	"log/slog"
	"os"

	"wrs/catalog/ccli/packages/cmd"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configFile config.ConfigData
var indent string
var NewLogWriter config.LogWriter

// initialize configuration file and flag values
func init() {
	// set the config file and its path
	viper.SetConfigFile("ccli_config.yml")
	viper.AddConfigPath(".")
	// create a default slog logger which logs to stdout
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
	// read the config file
	if err := viper.ReadInConfig(); err != nil {
		if err != nil {
			if err != errors.New("open ccli_config.yml: no such file or directory") {
				slog.Error("User configuration file not found. Please create ccli_config.yml and copy the contents of ccli_config.DEFAULT.yml.")
			} else {
				slog.Error("Error reading in config file", slog.Any("error", err))
			}
			os.Exit(1)
		}
	}
	// unmarshal the config file parameters to a struct
	if err := viper.Unmarshal(&configFile); err != nil {
		slog.Error("Could not unmarshal config file parameters")
		os.Exit(1)
	}
	// check if the log file is present and has the correct extension
	if configFile.LogFile == "" || configFile.LogFile[len(configFile.LogFile)-4:] != ".txt" {
		slog.Error("*** ERROR - Error reading config file, log file must be a .txt file")
		os.Exit(1)
	}
	// check if the log level is accurate
	if configFile.LogLevel > 2 || configFile.LogLevel < 1 {
		slog.Error("*** ERROR - Error reading log level, log level must be either 1 or 2")
		os.Exit(1)
	}
	indentString := ""
	for i := 0; i < int(configFile.JsonIndent); i++ {
		indentString += " "
	}
	indent = indentString
}

func main() {
	// check if the server address is provided
	if configFile.ServerAddr == "" {
		slog.Error("invalid configuration file, no server address located")
		os.Exit(1)
	}
	// contact the given server
	resp, err := http.DefaultClient.Get(configFile.ServerAddr)
	if err != nil {
		slog.Error("error contacting server", slog.Any("error:", err))
		os.Exit(1)
	}
	resp.Body.Close()
	// check if the response suggets a successful connection to the server
	if resp.StatusCode != 200 && resp.StatusCode != 422 {
		slog.Error("server connection error, check config file and network configuration", slog.Int("Status Code:", resp.StatusCode))
		os.Exit(1)
	}
	// create the log file or truncate it if already present
	logFile, err := os.Create(configFile.LogFile)
	if err != nil {
		slog.Error("*** ERROR - Error opening log file:", slog.Any("error", err))
		os.Exit(1)
	}
	// create a new log writer for writing to the log file and stdout simultaneously
	NewLogWriter = config.LogWriter{Stdout: os.Stderr, File: logFile}
	slogOptions := new(slog.HandlerOptions)
	slogOptions.Level = slog.LevelDebug
	if configFile.LogLevel == 2 {
		slogOptions.AddSource = true
	}
	// set the default slog logger to the log file with the given log level
	slog.SetDefault(slog.New(slog.NewJSONHandler(NewLogWriter.File, slogOptions)))
	slog.Debug("slog.SetDefault JSONHandler", slog.Group("HandlerOptions", slog.Bool("AddSource", slogOptions.AddSource), slog.Any("Level", slogOptions.Level)))
	client := graphql.GetNewClient(configFile.ServerAddr, http.DefaultClient)
	slog.Debug("successfully connected to server")
	// add all the sub commands to the root command
	rootCmd := cmd.RootCmd(&configFile, &NewLogWriter)
	rootCmd.AddCommand(cmd.Example())
	rootCmd.AddCommand(cmd.Ping(&configFile))
	rootCmd.AddCommand(cmd.Upload(&configFile))
	rootCmd.AddCommand(cmd.Update(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Query(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Find(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Export(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Add(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Delete(&configFile, client, indent))
	// bind and execute the root command and the sub commands
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", slog.Any("error", err))
		os.Exit(1)
	}

}
