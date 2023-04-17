package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// variable assignment for configuration file and command line flags
var configData config.ConfigData
var argSha256 string
var argQuery string
var indent string

// initialize configuration file and flag values
func init() {
	configFile, err := os.Open("ccli_config.DEFAULT.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("configuration file not found")
	}
	defer configFile.Close()
	data, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading configuration file")
	}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		log.Fatal().Err(err).Msg("error parsing config data")
	}
	indentString := ""
	for i := 0; i < int(configData.JsonIndent); i++ {
		indentString += " "
	}
	indent = indentString

	//TODO: possibly remove any base level flags and require subcommand
	// flag.StringVar(&argEndpoint, "endpoint", "https://aws-ip-services-tk-dev.wrs.com/api/graphql", "GraphQL Endpoint")
	// flag.IntVar(&argLogLevel, "log", int(zerolog.InfoLevel), "Log Level Minimum")
	flag.StringVar(&argSha256, "sha256", "", "SHA256 Part Query")
	flag.StringVar(&argQuery, "query", "", "Graphql Query")
}

func main() {
	// set global log level to value found in configuration file
	zerolog.SetGlobalLevel(zerolog.Level(configData.LogLevel))

	// open log file and set logging output
	logFile, err := os.Create(configData.LogFile)
	if err != nil {
		log.Fatal().Err(err).Msg("error opening log file")
	}
	defer logFile.Close()
	logger := zerolog.New(logFile)

	//parse command line flag values
	flag.Parse()

	client := graphql.GetNewClient(configData.ServerAddr, http.DefaultClient)

	// route based on subcommand
	subcommand := os.Args[1]
	switch subcommand {
	case "export":
		fmt.Println("Exporting")
	case "add":
		fmt.Println("Adding")
	default:
		if argSha256 != "" {
			partID, err := graphql.GetPartID(context.Background(), client, argSha256)
			if err != nil {
				fmt.Printf("Error retrieving part from SHA256: %s\n", argSha256)
				logger.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}

		if argQuery != "" {
			response, err := graphql.Query(context.Background(), client, argQuery)
			if err != nil {
				fmt.Printf("Error executing query, more information available in logs.")
				logger.Fatal().Err(err).Msg("error querying graphql")
			}

			var data map[string]interface{}
			json.Unmarshal(response, &data)

			prettyJson, err := json.MarshalIndent(data, "", indent)
			if err != nil {
				fmt.Println("error prettifying json response, more information available in logs.")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Println(string(prettyJson))
		}
	}
}
