package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"
	"wrs/catalog/ccli/packages/json"
	"wrs/catalog/ccli/packages/yaml"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// variable assignment for configuration file and command line flags
var configData config.ConfigData
var indent string
var argPartID string
var argSHA256 string
var argExportPath string
var argImportPath string
var argFVC string
var argSearchQuery string

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
	addSubcommand.StringVar(&argImportPath, "profile", "", "add profile import path")

	exportSubcommand = flag.NewFlagSet("export", flag.ExitOnError)
	exportSubcommand.StringVar(&argExportPath, "o", "", "output path for export subcommand")
	exportSubcommand.StringVar(&argPartID, "id", "", "part id for export subcommand")
	exportSubcommand.StringVar(&argSHA256, "sha256", "", "sha256 for export subcommand")
	exportSubcommand.StringVar(&argFVC, "fvc", "", "file verification code for export subcommand")

	querySubcommand = flag.NewFlagSet("query", flag.ExitOnError)

	findSubcommand = flag.NewFlagSet("find", flag.ExitOnError)
	findSubcommand.StringVar(&argSHA256, "sha256", "", "retrieve part id using sha256")
	findSubcommand.StringVar(&argSearchQuery, "part", "", "retrieve part data using search query")
	findSubcommand.StringVar(&argFVC, "fvc", "", "retrieve part id using file verification code")

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
			logger.Fatal().Err(err).Msg("error parsing export subcommand flags")
		}
		if argExportPath == "" {
			fmt.Println("*** ERROR - Path required to export data")
			logger.Fatal().Msg("error exporting part, no path given")
		}
		if argExportPath[len(argExportPath)-5:] != ".yaml" && argExportPath[len(argExportPath)-4:] != ".yml" {
			fmt.Println("*** ERROR - Export path must be a .yaml or .yml file")
			logger.Fatal().Msg("error exporting part, export path not a yaml file")
		}
		if argPartID == "" && argFVC == "" && argSHA256 == "" {
			fmt.Println("*** ERROR - Part ID required to export data")
			logger.Fatal().Msg("error exporting part, no part identifier given")
		}
		var part *graphql.Part
		if argPartID != "" {
			part, err = graphql.GetPartByID(context.Background(), client, argPartID)
			if err != nil {
				fmt.Println("*** ERROR - Error retrieving part by catalog id, check logs for more info")
				logger.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		if argSHA256 != "" {
			part, err = graphql.GetPartBySHA256(context.Background(), client, argSHA256)
			if err != nil {
				fmt.Println("*** ERROR - Error retrieving part by sha256, check logs for more info")
				logger.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		if argFVC != "" {
			part, err = graphql.GetPartByFVC(context.Background(), client, argFVC)
			if err != nil {
				fmt.Println("*** ERROR - Error retrieving part by file verification code, check logs for more info")
				logger.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		yamlPart, err := yaml.Marshal(part)
		if err != nil {
			fmt.Println("*** ERROR - Error marshalling part into yaml")
			logger.Fatal().Err(err).Msg("error marshalling yaml")
		}

		yamlFile, err := os.Create(argExportPath)
		if err != nil {
			fmt.Println("*** ERROR - Error creating yaml file")
			logger.Fatal().Err(err).Msg("error creating yaml file")
		}
		defer yamlFile.Close()

		_, err = yamlFile.Write(yamlPart)
		if err != nil {
			fmt.Println("*** ERROR - Error writing yaml file")
			logger.Fatal().Err(err).Msg("error writing part to yaml file")
		}
		fmt.Printf("Part successfully exported to path: %s\n", argExportPath)

	case "add":
		if err := addSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error adding part")
			logger.Fatal().Err(err).Msg("error parsing add subcommand flags")
		}
		if argImportPath == "" {
			fmt.Println("*** ERROR - Profile data required to add a profile")
			logger.Fatal().Msg("error adding profile, no import path given")
		}
		if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
			fmt.Println("*** ERROR - Import path must be a .yaml or .yml file")
			logger.Fatal().Msg("error importing profile, import path not a yaml file")
		}
		f, err := os.Open(argImportPath)
		if err != nil {
			fmt.Println("*** ERROR - Error opening profile file, check logs for more info")
			logger.Fatal().Err(err).Msg("error opening file")
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("*** ERROR - Error reading file")
			logger.Fatal().Err(err).Msg("error reading file")
		}
		var profileData yaml.Profile
		if err = yaml.Unmarshal(data, &profileData); err != nil {
			fmt.Println("*** ERROR - Error decoding file contents")
			logger.Fatal().Err(err).Msg("error decoding file contents")
		}
		switch profileData.Profile {
		case "security":
			var securityProfile yaml.SecurityProfile
			if err = yaml.Unmarshal(data, &securityProfile); err != nil {
				fmt.Println("*** ERROR - Error unmarshaling security profile")
				logger.Fatal().Err(err).Msg("error unmarshaling security profile")
			}
			jsonSecurityProfile, err := json.Marshal(securityProfile)
			if err != nil {
				fmt.Println("*** ERROR - Error marshaling json")
				logger.Fatal().Err(err).Msg("error marshaling json")
			}
			if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
				fmt.Println("*** ERROR - Error adding profile, no part identifier given")
				logger.Fatal().Msg("error adding profile, no part identifier given")
			}
			if profileData.CatalogID != "" {
				if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonSecurityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
			}
			if profileData.FVC != "" {
				uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by fvc")
					logger.Fatal().Err(err).Msg("error retrieving part id by fvc")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
				break
			}
			if profileData.Sha256 != "" {
				uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by sha256")
					logger.Fatal().Err(err).Msg("error retrieving part id by sha256")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
			}
		case "licensing":
			var licensingProfile yaml.LicensingProfile
			if err = yaml.Unmarshal(data, &licensingProfile); err != nil {
				fmt.Println("*** ERROR - Error unmarshaling licensing profile")
				logger.Fatal().Err(err).Msg("error unmarshaling licensing profile")
			}
			jsonLicensingProfile, err := json.Marshal(licensingProfile)
			if err != nil {
				fmt.Println("*** ERROR - Error marshaling json")
				logger.Fatal().Err(err).Msg("error marshaling json")
			}
			if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
				fmt.Println("*** ERROR - Error adding profile, no part identifier given")
				logger.Fatal().Msg("error adding profile, no part identifier given")
			}
			if profileData.CatalogID != "" {
				if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonLicensingProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
			}
			if profileData.FVC != "" {
				uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by fvc")
					logger.Fatal().Err(err).Msg("error retrieving part id by fvc")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
				break
			}
			if profileData.Sha256 != "" {
				uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by sha256")
					logger.Fatal().Err(err).Msg("error retrieving part id by sha256")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
			}
		case "quality":
			var qualityProfile yaml.QualityProfile
			if err = yaml.Unmarshal(data, &qualityProfile); err != nil {
				fmt.Println("*** ERROR - Error unmarshaling quality profile")
				logger.Fatal().Err(err).Msg("error unmarshaling quality profile")
			}
			jsonQualityProfile, err := json.Marshal(qualityProfile)
			if err != nil {
				fmt.Println("*** ERROR - Error marshaling json")
				logger.Fatal().Err(err).Msg("error marshaling json")
			}
			if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
				fmt.Println("*** ERROR - Error adding profile, no part identifier given")
				logger.Fatal().Msg("error adding profile, no part identifier given")
			}
			if profileData.CatalogID != "" {
				if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonQualityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
			}
			if profileData.FVC != "" {
				uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by fvc")
					logger.Fatal().Err(err).Msg("error retrieving part id by fvc")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
				break
			}
			if profileData.Sha256 != "" {
				uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
				if err != nil {
					fmt.Println("*** ERROR - Error retrieving part id by sha256")
					logger.Fatal().Err(err).Msg("error retrieving part id by sha256")
				}
				if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
					fmt.Println("*** ERROR - Error adding profile, check logs for more info")
					logger.Fatal().Err(err).Msg("error adding profile")
				}
				fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
			}
		}
	case "find":
		if err := findSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error finding part")
			logger.Fatal().Err(err).Msg("error parsing find subcommand flags")
		}
		if argSHA256 == "" && argSearchQuery == "" && argFVC == "" {
			fmt.Println("*** ERROR - error finding part, find part usage: ccli find -part <SHA256>")
			logger.Fatal().Msg("error finding part, no data provided")
		}
		if argSHA256 != "" {
			partID, err := graphql.GetPartIDBySha256(context.Background(), client, argSHA256)
			if err != nil {
				fmt.Printf("*** ERROR - Error retrieving part from SHA256: %s\n", argSHA256)
				logger.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argFVC != "" {
			partID, err := graphql.GetPartIDByFVC(context.Background(), client, argFVC)
			if err != nil {
				fmt.Printf("*** ERROR - Error retrieving part from file verification code: %s\n", argFVC)
				logger.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argSearchQuery != "" {
			response, err := graphql.Search(context.Background(), client, argSearchQuery)
			if err != nil {
				fmt.Println("*** ERROR - Error searching for part")
				logger.Fatal().Err(err).Msg("error searching for part")
			}

			prettyJson, err := json.MarshalIndent(response, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - Error prettifying json")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("Result: %s\n", string(prettyJson))
		}
	case "query":
		if err := querySubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - error executing query, query subcommand usage: ccli query <GraphQL Query>")
			logger.Fatal().Err(err).Msg("error parsing query subcommand")
		}
		argQuery := querySubcommand.Arg(0)
		if argQuery == "" {
			fmt.Println("*** ERROR - error executing query, query subcommand usage: ccli query <GraphQL Query>")
			logger.Fatal().Msg("error executing user query, no data provided")
		}
		if argQuery != "" {
			response, err := graphql.Query(context.Background(), client, argQuery)
			if err != nil {
				fmt.Printf("***ERROR - Error executing query, check logs for more info")
				logger.Fatal().Err(err).Msg("error querying graphql")
			}

			// json result will be output in prettified format
			var data map[string]interface{}
			json.Unmarshal(response, &data)

			prettyJson, err := json.MarshalIndent(data, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - error prettifying json response, check logs for more info")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Println(string(prettyJson))
		}
	default:
		printHelp()
	}
}
