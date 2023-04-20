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
var argPartID string
var argSHA256 string
var argExportPath string

var addSubcommand *flag.FlagSet
var exportSubcommand *flag.FlagSet
var querySubcommand *flag.FlagSet
var findSubcommand *flag.FlagSet

// initialize configuration file and flag values
func init() {
	configFile, err := os.Open("ccli_config.yml")
	if err != nil {
		fmt.Println("User configuration file not found, using default.")
	}
	if configFile == nil {
		configFile, err = os.Open("ccli_config.DEFAULT.yml")
		if err != nil {
			fmt.Println("*** ERROR - Default configuration file not found")
		}
	}
	defer configFile.Close()
	data, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("*** ERROR - Error reading configuration file")
	}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		log.Fatal().Err(err).Msg("*** ERROR - Error parsing config data")
	}
	indentString := ""
	for i := 0; i < int(configData.JsonIndent); i++ {
		indentString += " "
	}
	indent = indentString

}

func printHelp() {
	fmt.Println("Please enter a command:")
	fmt.Println("add")
	addSubcommand.PrintDefaults()
	fmt.Println("export")
	exportSubcommand.PrintDefaults()
	fmt.Println("query")
	fmt.Println("\tUsed to input custom GraphQL query - ccli query <Graphql Query>")
	fmt.Println("find")
	findSubcommand.PrintDefaults()
	os.Exit(1)
}

func main() {
	// set global log level to value found in configuration file
	zerolog.SetGlobalLevel(zerolog.Level(configData.LogLevel))

	// open log file and set logging output
	logFile, err := os.Create(configData.LogFile)
	if err != nil {
		fmt.Println("*** ERROR - Error opening log file")
	}
	defer logFile.Close()
	logger := zerolog.New(logFile)

	//subcommand flag sets
	addSubcommand = flag.NewFlagSet("add", flag.ExitOnError)
	addSubcommand.StringVar(&argPartID, "part", "", "add part test flag")

	exportSubcommand = flag.NewFlagSet("export", flag.ExitOnError)
	exportSubcommand.StringVar(&argExportPath, "o", "", "output path for export subcommand")
	exportSubcommand.StringVar(&argPartID, "part", "", "part id for export subcommand")

	querySubcommand = flag.NewFlagSet("query", flag.ExitOnError)

	findSubcommand = flag.NewFlagSet("find", flag.ExitOnError)
	findSubcommand.StringVar(&argSHA256, "sha256", "", "retrieve part id using sha256 for find subcommand")

	if len(os.Args) < 2 {
		printHelp()
	}

	client := graphql.GetNewClient(configData.ServerAddr, http.DefaultClient)

	// route based on subcommand
	subcommand := os.Args[1]
	switch subcommand {
	case "export":
		if err := exportSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error exporting data")
		}
		if argExportPath[:4] != ".yaml" {
			fmt.Println("*** ERROR - Export path must be a .yaml file")
		}
		if argPartID != "" && argExportPath != "" {
			fmt.Printf("Now exporting part: %s to path:%s\n", argPartID, argExportPath)
		}
		if argPartID == "" {
			fmt.Println("*** ERROR - Part ID required to export data")
		}
		if argExportPath == "" {
			fmt.Println("*** ERROR - Path required to export data")
		}
	case "add":
		if err := addSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error adding part")
		}
		if argPartID != "" {
			fmt.Printf("Now adding part: %s\n", argPartID)
		}
		if argPartID == "" {
			fmt.Println("*** ERROR - Part data required to add")
		}
	case "find":
		if err := findSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - error finding part")
		}
		if argSHA256 != "" {
			partID, err := graphql.GetPartID(context.Background(), client, argSHA256)
			if err != nil {
				fmt.Printf("*** ERROR - Error retrieving part from SHA256: %s\n", argSHA256)
				logger.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argSHA256 == "" {
			fmt.Println("*** ERROR - error finding part, find part usage: ccli find -part <SHA256>")
		}
	case "query":
		if err := querySubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - error executing query, query subcommand usage: ccli query <GraphQL Query>")
		}
		argQuery := querySubcommand.Arg(0)
		if argQuery != "" {
			response, err := graphql.Query(context.Background(), client, argQuery)
			if err != nil {
				fmt.Printf("***ERROR - Error executing query, more information available in logs.")
				logger.Fatal().Err(err).Msg("error querying graphql")
			}

			// json result will be output in prettified format
			var data map[string]interface{}
			json.Unmarshal(response, &data)

			prettyJson, err := json.MarshalIndent(data, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - error prettifying json response, more information available in logs.")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Println(string(prettyJson))
		}
		if argQuery == "" {
			fmt.Println("*** ERROR - error executing query, query subcommand usage: ccli query <GraphQL Query>")
		}

	default:
		printHelp()
	}
}
