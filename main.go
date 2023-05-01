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
var argPartImportPath string
var argFVC string
var argSearchQuery string
var argTemplate string
var argProfileType string

var addSubcommand *flag.FlagSet
var exportSubcommand *flag.FlagSet
var querySubcommand *flag.FlagSet
var findSubcommand *flag.FlagSet
var uploadSubcommand *flag.FlagSet
var updateSubcommand *flag.FlagSet

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
	fmt.Println("upload")
	fmt.Println("\tUsed to upload packages - ccli upload <Path>")
	fmt.Println("update")
	fmt.Println("\tUsed to update part data - ccli update <Path>")
	os.Exit(1)
}

func main() {
	// set global log level to value found in configuration file
	zerolog.SetGlobalLevel(zerolog.Level(configData.LogLevel))

	if configData.LogFile == "" || configData.LogFile[len(configData.LogFile)-4:] != ".txt" {
		log.Fatal().Msg("*** ERROR - Error reading config file, log file must be a .txt file")
	}
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
	addSubcommand.StringVar(&argPartImportPath, "part", "", "add part import path")

	exportSubcommand = flag.NewFlagSet("export", flag.ExitOnError)
	exportSubcommand.StringVar(&argExportPath, "o", "", "output path for export subcommand")
	exportSubcommand.StringVar(&argPartID, "id", "", "part id for export subcommand")
	exportSubcommand.StringVar(&argSHA256, "sha256", "", "sha256 for export subcommand")
	exportSubcommand.StringVar(&argFVC, "fvc", "", "file verification code for export subcommand")
	exportSubcommand.StringVar(&argTemplate, "template", "", "used to export profile/part template to output: ./ccli export -template <Type>")

	querySubcommand = flag.NewFlagSet("query", flag.ExitOnError)

	findSubcommand = flag.NewFlagSet("find", flag.ExitOnError)
	findSubcommand.StringVar(&argProfileType, "profile", "", "retrieve profile using key")
	findSubcommand.StringVar(&argSHA256, "sha256", "", "retrieve part id using sha256")
	findSubcommand.StringVar(&argSearchQuery, "part", "", "retrieve part data using search query")
	findSubcommand.StringVar(&argFVC, "fvc", "", "retrieve part id using file verification code")
	findSubcommand.StringVar(&argPartID, "id", "", "retrieve part data using part id")

	updateSubcommand = flag.NewFlagSet("update", flag.ExitOnError)
	updateSubcommand.StringVar(&argImportPath, "part", "", "update part import path")

	uploadSubcommand = flag.NewFlagSet("upload", flag.ExitOnError)

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
		if argPartID == "" && argFVC == "" && argSHA256 == "" && argTemplate == "" {
			fmt.Println("*** ERROR - Part identifier required to export data")
			logger.Fatal().Msg("error exporting part, no part identifier given")
		}
		if argTemplate != "" {
			switch argTemplate {
			case "part":
				yamlPart := new(yaml.Part)
				yamlPart.Format = 1.0
				f, err := os.Create(argExportPath)
				if err != nil {
					fmt.Println("*** ERROR - Error creating template file")
					logger.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlPartTemplate, err := yaml.Marshal(&yamlPart)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling part template")
					logger.Fatal().Err(err).Msg("error marshaling part template")
				}
				_, err = f.Write(yamlPartTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				fmt.Printf("Part template successfully output to: %s\n", argExportPath)
				os.Exit(0)
			case "security":
				yamlProfile := new(yaml.Profile)
				yamlSecurityProfile := new(yaml.SecurityProfile)
				yamlCVE := new(yaml.CVE)
				yamlProfile.Format = 1.0
				yamlSecurityProfile.CVEList = append(yamlSecurityProfile.CVEList, *yamlCVE)
				f, err := os.Create(argExportPath)
				if err != nil {
					fmt.Println("*** ERROR - Error creating template file")
					logger.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling profile template")
					logger.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlSecurityProfileTemplate, err := yaml.Marshal(&yamlSecurityProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling security profile template")
					logger.Fatal().Err(err).Msg("error marshaling security profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlSecurityProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
				os.Exit(0)
			case "quality":
				yamlProfile := new(yaml.Profile)
				yamlQualityProfile := new(yaml.QualityProfile)
				yamlBug := new(yaml.Bug)
				yamlProfile.Format = 1.0
				yamlQualityProfile.BugList = append(yamlQualityProfile.BugList, *yamlBug)
				f, err := os.Create(argExportPath)
				if err != nil {
					fmt.Println("*** ERROR - Error creating template file")
					logger.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling profile template")
					logger.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlQualityProfileTemplate, err := yaml.Marshal(&yamlQualityProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling quality profile template")
					logger.Fatal().Err(err).Msg("error marshaling quality profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlQualityProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
				os.Exit(0)
			case "licensing":
				yamlProfile := new(yaml.Profile)
				yamlLicensingProfile := new(yaml.LicensingProfile)
				yamlLicense := new(yaml.License)
				yamlProfile.Format = 1.0
				yamlLicensingProfile.LicenseAnalysis = append(yamlLicensingProfile.LicenseAnalysis, *yamlLicense)
				f, err := os.Create(argExportPath)
				if err != nil {
					fmt.Println("*** ERROR - Error creating template file")
					logger.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling profile template")
					logger.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlLicensingProfileTemplate, err := yaml.Marshal(&yamlLicensingProfile)
				if err != nil {
					fmt.Println("*** ERROR - Error marshaling licensing profile template")
					logger.Fatal().Err(err).Msg("error marshaling licensing profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlLicensingProfileTemplate)
				if err != nil {
					fmt.Println("*** ERROR - Error writing template to file")
					logger.Fatal().Err(err).Msg("error writing template to file")
				}
				fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
				os.Exit(0)
			}
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
		var yamlPart yaml.Part
		if err := graphql.UnmarshalPart(part, &yamlPart); err != nil {
			fmt.Println("*** ERROR parsing part into yaml")
			logger.Fatal().Err(err).Msg("error parsing part into yaml")
		}

		ret, err := yaml.Marshal(yamlPart)
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

		_, err = yamlFile.Write(ret)
		if err != nil {
			fmt.Println("*** ERROR - Error writing yaml file")
			logger.Fatal().Err(err).Msg("error writing part to yaml file")
		}
		fmt.Printf("Part successfully exported to path: %s\n", argExportPath)

	case "add":
		if err := addSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error parsing add subcommand flags")
			logger.Fatal().Err(err).Msg("error parsing add subcommand flags")
		}
		if argImportPath == "" && argPartImportPath == "" {
			fmt.Println("*** ERROR - Data required to add a part/profile, usage: ./ccli add -profile|-part <Path>")
			logger.Fatal().Msg("error adding data, no import path given")
		}
		if argImportPath != "" {
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
				fmt.Println("*** ERROR - Error unmarshaling file contents")
				logger.Fatal().Err(err).Msg("error unmarshaling file contents")
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
				fmt.Printf("%s\n", string(jsonSecurityProfile))
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
		}
		if argPartImportPath != "" {
			if argPartImportPath[len(argPartImportPath)-5:] != ".yaml" && argPartImportPath[len(argPartImportPath)-4:] != ".yml" {
				fmt.Println("*** ERROR - Import path must be a .yaml or .yml file")
				logger.Fatal().Msg("error importing part, import path not a yaml file")
			}
			f, err := os.Open(argPartImportPath)
			if err != nil {
				fmt.Println("*** ERROR - Error opening part file, check logs for more info")
				logger.Fatal().Err(err).Msg("error opening file")
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				fmt.Println("*** ERROR - Error reading file")
				logger.Fatal().Err(err).Msg("error reading file")
			}
			var partData yaml.Part
			if err = yaml.Unmarshal(data, &partData); err != nil {
				fmt.Println("*** ERROR - Error unmarshaling file contents")
				logger.Fatal().Err(err).Msg("error unmarshaling file contents")
			}
			createdPart, err := graphql.AddPart(context.Background(), client, partData)
			if err != nil {
				fmt.Println("*** ERROR - Error adding part, check logs for more info")
				logger.Fatal().Err(err).Msg("error adding part")
			}
			prettyPart, err := json.MarshalIndent(&createdPart, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - Error prettifying json response")
			}
			fmt.Printf("Successfully added part from: %s\n", argPartImportPath)
			fmt.Printf("%s\n", string(prettyPart))

		}
	case "find":
		if err := findSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error finding part")
			logger.Fatal().Err(err).Msg("error parsing find subcommand flags")
		}
		if argSHA256 == "" && argPartID == "" && argFVC == "" && argSearchQuery == "" {
			fmt.Println("*** ERROR - error finding part, find part usage:")
			findSubcommand.PrintDefaults()
			logger.Fatal().Msg("error finding part, no data provided")
		}
		if argProfileType != "" {
			if argPartID == "" {
				fmt.Println("*** ERROR - Error retrieving profile, part id needed: ./ccli find -profile <type> -id <part_id>")
				logger.Fatal().Msg("error getting profile, missing part id")
			}
			profile, err := graphql.GetProfile(context.Background(), client, argPartID, argProfileType)
			if err != nil {
				fmt.Println("*** ERROR - Error retrieving profile, check logs for more info")
				logger.Fatal().Err(err).Msg("error retrieving profile")
			}
			prettyJson, err := json.MarshalIndent(&profile, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - Error prettifying json")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("%s\n", string(prettyJson))
			os.Exit(0)
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
		if argPartID != "" {
			response, err := graphql.GetPartByID(context.Background(), client, argPartID)
			if err != nil {
				fmt.Println("*** ERROR - Error getting part by id")
				logger.Fatal().Err(err).Msg("error getting part by id")
			}
			prettyJson, err := json.MarshalIndent(&response, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - Error prettifying response")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("%s\n", string(prettyJson))
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
	case "upload":
		if err := uploadSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error executing upload, upload subcommand usage: ccli upload <Path>")
			logger.Fatal().Err(err).Msg("error parsing upload subcommand")
		}
		argPath := uploadSubcommand.Arg(0)
		if argPath == "" {
			fmt.Println("*** ERROR - Error executing upload, upload subcommand usage: ccli upload <Path>")
			logger.Fatal().Msg("error executing upload, no path given")
		}
		if argPath != "" {
			response, err := graphql.UploadFile(http.DefaultClient, configData.ServerAddr, argPath, "")
			if err != nil {
				fmt.Println("*** ERROR - Error executing upload, check logs for more info")
				logger.Fatal().Err(err).Msg("error uploading archive")
			}
			if response.StatusCode == 200 {
				fmt.Printf("Successfully uploaded package: %s\n", argPath)
			}
		}
	case "update":
		if err := updateSubcommand.Parse(os.Args[2:]); err != nil {
			fmt.Println("*** ERROR - Error updating part, update subcommand usage: ./ccli update -part <Path>")
			logger.Fatal().Err(err).Msg("error updating part")
		}
		if argImportPath == "" {
			fmt.Println("*** ERROR - Error updating part, update subcommand usage: ./ccli update -part <Path>")
			logger.Fatal().Msg("error updating part, no file given")
		}
		if argImportPath != "" {
			if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
				fmt.Println("*** ERROR - Import path must be a .yaml or .yml file")
				logger.Fatal().Msg("error importing part, import path not a yaml file")
			}
			f, err := os.Open(argImportPath)
			if err != nil {
				fmt.Println("*** ERROR - Error opening part file, check logs for more info")
				logger.Fatal().Err(err).Msg("error opening file")
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				fmt.Println("*** ERROR - Error reading file")
				logger.Fatal().Err(err).Msg("error reading file")
			}
			var partData yaml.Part
			if err = yaml.Unmarshal(data, &partData); err != nil {
				fmt.Println("*** ERROR - Error decoding file contents")
				logger.Fatal().Err(err).Msg("error decoding file contents")
			}
			returnPart, err := graphql.UpdatePart(context.Background(), client, &partData)
			if err != nil {
				fmt.Println("*** ERROR - Error updating part, check logs for more info")
				logger.Fatal().Err(err).Msg("error updating part")
			}
			prettyJson, err := json.MarshalIndent(&returnPart, "", indent)
			if err != nil {
				fmt.Println("*** ERROR - Error prettifying returned part")
				logger.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("Part successfully updated\n%s\n", string(prettyJson))
		}
	default:
		printHelp()
	}
}
