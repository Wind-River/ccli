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
var indent string
var argPartIdentifier string
var argExportPath string

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

	//subcommand flag sets
	addSubcommand := flag.NewFlagSet("add", flag.ExitOnError)
	addSubcommand.StringVar(&argPartIdentifier, "part", "", "Add part test flag")

	exportSubcommand := flag.NewFlagSet("export", flag.ExitOnError)
	exportSubcommand.StringVar(&argExportPath, "o", "", "output path for export subcommand")
	exportSubcommand.StringVar(&argPartIdentifier, "part", "", "part id for export subcommand")

	querySubcommand := flag.NewFlagSet("query", flag.ExitOnError)

	findSubcommand := flag.NewFlagSet("find", flag.ExitOnError)
	findSubcommand.StringVar(&argPartIdentifier, "part", "", "part id for find subcommand")

	client := graphql.GetNewClient(configData.ServerAddr, http.DefaultClient)

	// route based on subcommand
	subcommand := os.Args[1]
	switch subcommand {
	case "export":
		if err := exportSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("error exporting data")
		}
		if argPartIdentifier != "" && argExportPath != "" {
			fmt.Printf("Now exporting part: %s to path:%s\n", argPartIdentifier, argExportPath)
		}
		if argPartIdentifier == "" {
			fmt.Println("Part ID required to export data")
		}
		if argExportPath == "" {
			fmt.Println("Path required to export data")
		}
	case "add":
		if err := addSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("error adding part")
		}
		if argPartIdentifier != "" {
			fmt.Printf("Now adding part: %s\n", argPartIdentifier)
		}
		if argPartIdentifier == "" {
			fmt.Println("Part data required to add")
		}
	case "find":
		if err := findSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("error finding part")
		}
		if argPartIdentifier != "" {
			partID, err := graphql.GetPartID(context.Background(), client, argPartIdentifier)
			if err != nil {
				fmt.Printf("Error retrieving part from SHA256: %s\n", argPartIdentifier)
				logger.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argPartIdentifier == "" {
			fmt.Println("error finding part, find part usage: ccli find -part <SHA256>")
		}
	case "query":
		if err := querySubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("error executing query, query subcommand usage: ccli query <GraphQL Query>")
		}
		argQuery := querySubcommand.Arg(0)
		if argQuery != "" {
			response, err := graphql.Query(context.Background(), client, argQuery)
			if err != nil {
				fmt.Printf("Error executing query, more information available in logs.")
				logger.Fatal().Err(err).Msg("error querying graphql")
			}

			// json result will be output in prettified format
			var data map[string]interface{}
			json.Unmarshal(response, &data)

			prettyJson, err := json.MarshalIndent(data, "", indent)
			if err != nil {
				fmt.Println("error prettifying json response, more information available in logs.")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Println(string(prettyJson))
		}
		if argQuery == "" {
			fmt.Println("error executing query, query subcommand usage: ccli query <GraphQL Query>")
		}

	default:
		fmt.Println(flag.ErrHelp)
	}
}
