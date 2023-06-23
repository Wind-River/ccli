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
var argExampleMode bool
var argVerboseMode bool
var argHelp bool
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
var deleteSubcommand *flag.FlagSet
var pingSubcommand *flag.FlagSet

// initialize configuration file and flag values
func init() {
	configFile, err := os.Open("ccli_config.yml")
	if err != nil {
		fmt.Println("User configuration file not found. Please create ccli_config.yml and copy the contents of ccli_config.DEFAULT.yml.")
		os.Exit(1)
	}
	defer configFile.Close()
	data, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Printf("*** ERROR - Error reading configuration file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		fmt.Printf("*** ERROR - Error parsing config data: %s\n", err.Error())
		os.Exit(1)
	}
	indentString := ""
	for i := 0; i < int(configData.JsonIndent); i++ {
		indentString += " "
	}
	indent = indentString
	flag.BoolVar(&argExampleMode, "e", false, "used to print subcommand usage examples")
	flag.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")
	flag.BoolVar(&argHelp, "h", false, "used to print default subcommand usage")
	flag.BoolVar(&argHelp, "help", false, "used to print default subcommand usage")

}

// Prints default subcommand usage for help flag and subcommand errors
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
	updateSubcommand.PrintDefaults()
	fmt.Println("delete")
	deleteSubcommand.PrintDefaults()
	fmt.Println("ping")
	fmt.Println("\tUsed to ping server address")
	fmt.Println("-e - Print a list of example subcommand usage")
	os.Exit(0)
}

func main() {
	flag.Parse()
	//Used to print examples of ccli usage before checking server connection
	if argExampleMode {
		exampleString :=
			`	$ ccli add --part openssl-1.1.1n.yml
	$ ccli add --profile profile_openssl-1.1.1n.yml
	$ ccli query "{part(id:\"aR25sd-V8dDvs2-p3Gfae\"){file_verification_code}}"
	$ ccli export --id sdl3ga-naTs42g5-rbow2A -o file.yml
	$ ccli export --template security -o file.yml
	$ ccli update --part openssl-1.1.1n.v4.yml
	$ ccli upload openssl-1.1.1n.tar.gz
	$ ccli find -part busybox
	$ ccli find -sha256 2493347f59c03...
	$ ccli find -profile security -id werS12-da54FaSff-9U2aef
	$ ccli delete -id adjb23-A4D3faTa-d95Xufs
	$ ccli ping
	$ ccli -e`
		fmt.Printf("%s\n", exampleString)
		os.Exit(0)
	}

	// set global log level to value found in configuration file
	zerolog.SetGlobalLevel(zerolog.Level(configData.LogLevel))
	if argVerboseMode {
		zerolog.SetGlobalLevel(0)
	}

	if configData.LogFile == "" || configData.LogFile[len(configData.LogFile)-4:] != ".txt" {
		fmt.Println("*** ERROR - Error reading config file, log file must be a .txt file")
		os.Exit(1)
	}
	// open log file and set logging output
	logFile, err := os.Create(configData.LogFile)
	if err != nil {
		fmt.Printf("*** ERROR - Error opening log file: %s\n", err.Error())
	}
	defer logFile.Close()
	multiLogger := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	log.Logger = log.Output(multiLogger)
	log.Debug().Msgf("log file opened at: %s", configData.LogFile)

	//subcommand flag sets
	addSubcommand = flag.NewFlagSet("add", flag.ExitOnError)
	addSubcommand.StringVar(&argImportPath, "profile", "", "add profile import path, ccli usage: add --profile <file.yaml>")
	addSubcommand.StringVar(&argPartImportPath, "part", "", "add part import path, ccli usage: add --part <file.yaml>")
	addSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	exportSubcommand = flag.NewFlagSet("export", flag.ExitOnError)
	exportSubcommand.StringVar(&argExportPath, "o", "", "output path for export subcommand, ccli usage: -o <file.yaml>")
	exportSubcommand.StringVar(&argPartID, "id", "", "part id for export subcommand, ccli usage: --id <catalog_id>")
	exportSubcommand.StringVar(&argSHA256, "sha256", "", "sha256 for export subcommand, ccli usage: --sha256 <SHA256>")
	exportSubcommand.StringVar(&argFVC, "fvc", "", "file verification code for export subcommand, ccli usage: --fvc <file_verification_code>")
	exportSubcommand.StringVar(&argTemplate, "template", "", "used to export profile/part template to output, ccli usage: export --template <Type>")
	exportSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	querySubcommand = flag.NewFlagSet("query", flag.ExitOnError)
	querySubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	findSubcommand = flag.NewFlagSet("find", flag.ExitOnError)
	findSubcommand.StringVar(&argProfileType, "profile", "", "retrieve profile using key, ccli usage: --profile <type> --id <catalog_id>")
	findSubcommand.StringVar(&argSHA256, "sha256", "", "retrieve part id using sha256, ccli usage: --sha256 <SHA256>")
	findSubcommand.StringVar(&argSearchQuery, "part", "", "retrieve part data using search query, ccli usage: --part <name>")
	findSubcommand.StringVar(&argFVC, "fvc", "", "retrieve part id using file verification code, ccli usage: --fvc <file_verification_code>")
	findSubcommand.StringVar(&argPartID, "id", "", "retrieve part data using part id, ccli usage: --id <catalog_id>")
	findSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	updateSubcommand = flag.NewFlagSet("update", flag.ExitOnError)
	updateSubcommand.StringVar(&argImportPath, "part", "", "update part import path, ccli usage: --part <file.yaml>")
	updateSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	uploadSubcommand = flag.NewFlagSet("upload", flag.ExitOnError)
	uploadSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	deleteSubcommand = flag.NewFlagSet("delete", flag.ExitOnError)
	deleteSubcommand.StringVar(&argPartID, "id", "", "delete part from catalog using catalog id, ccli usage: --id <catalog_id>")
	deleteSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	pingSubcommand = flag.NewFlagSet("ping", flag.ExitOnError)
	pingSubcommand.BoolVar(&argVerboseMode, "v", false, "used to run ccli in verbose mode")

	if len(os.Args) < 2 {
		printHelp()
	}

	if argHelp {
		printHelp()
	}

	if configData.ServerAddr == "" {
		log.Fatal().Msg("invalid configuration file, no server address located")
	}
	resp, err := http.DefaultClient.Get(configData.ServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("error contacting server")
	}
	resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 422 {
		log.Fatal().Msgf("server connection error, check config file and network configuration: Status Code (%d)\n", resp.StatusCode)
	}
	log.Debug().Msgf("successfully connected to server")

	client := graphql.GetNewClient(configData.ServerAddr, http.DefaultClient)

	// route based on subcommand
	subcommand := os.Args[1]
	switch subcommand {
	case "export":
		if err := exportSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error parsing export subcommand flags")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if argPartID == "" && argFVC == "" && argSHA256 == "" && argTemplate == "" {
			log.Fatal().Msg("error exporting part, no part identifier given")
		}
		if argExportPath == "" {
			log.Fatal().Msg("error exporting part, no path given")
		}
		if argExportPath[len(argExportPath)-5:] != ".yaml" && argExportPath[len(argExportPath)-4:] != ".yml" {
			log.Fatal().Msg("error exporting part, export path not a yaml file")
		}
		if argTemplate != "" {
			switch argTemplate {
			case "part":
				yamlPart := new(yaml.Part)
				yamlPart.Format = 1.0
				f, err := os.Create(argExportPath)
				if err != nil {
					log.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlPartTemplate, err := yaml.Marshal(&yamlPart)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling part template")
				}
				_, err = f.Write(yamlPartTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
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
					log.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlSecurityProfileTemplate, err := yaml.Marshal(&yamlSecurityProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling security profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlSecurityProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
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
					log.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlQualityProfileTemplate, err := yaml.Marshal(&yamlQualityProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling quality profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlQualityProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
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
					log.Fatal().Err(err).Msg("error creating template file")
				}
				defer f.Close()
				yamlProfileTemplate, err := yaml.Marshal(&yamlProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling profile template")
				}
				yamlLicensingProfileTemplate, err := yaml.Marshal(&yamlLicensingProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling licensing profile template")
				}
				_, err = f.Write(yamlProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
				}
				_, err = f.Write(yamlLicensingProfileTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("error writing template to file")
				}
				fmt.Printf("Profile template successfully output to: %s\n", argExportPath)
				os.Exit(0)
			}
		}
		var part *graphql.Part
		if argPartID != "" {
			log.Debug().Str("ID", argPartID).Msg("retrieving part by id")
			part, err = graphql.GetPartByID(context.Background(), client, argPartID)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		if argSHA256 != "" {
			log.Debug().Str("SHA256", argSHA256).Msg("retrieving part by sha256")
			part, err = graphql.GetPartBySHA256(context.Background(), client, argSHA256)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		if argFVC != "" {
			log.Debug().Str("File Verification Code", argFVC).Msg("retrieving part by file verification code")
			part, err = graphql.GetPartByFVC(context.Background(), client, argFVC)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part")
			}
		}
		var yamlPart yaml.Part
		if err := graphql.UnmarshalPart(part, &yamlPart); err != nil {
			log.Fatal().Err(err).Msg("error parsing part into yaml")
		}

		ret, err := yaml.Marshal(yamlPart)
		if err != nil {
			log.Fatal().Err(err).Msg("error marshalling yaml")
		}

		yamlFile, err := os.Create(argExportPath)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating yaml file")
		}
		defer yamlFile.Close()

		_, err = yamlFile.Write(ret)
		if err != nil {
			log.Fatal().Err(err).Msg("error writing part to yaml file")
		}
		fmt.Printf("Part successfully exported to path: %s\n", argExportPath)

	case "add":
		if err := addSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error parsing add subcommand flags")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if argImportPath == "" && argPartImportPath == "" {
			log.Fatal().Msg("error adding data, no import path given")
		}
		if argImportPath != "" {
			if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
				log.Fatal().Msg("error importing profile, import path not a yaml file")
			}
			f, err := os.Open(argImportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error opening file")
			}
			defer f.Close()
			log.Debug().Msgf("successfully opened file at: %s", argImportPath)
			data, err := io.ReadAll(f)
			if err != nil {
				log.Fatal().Err(err).Msg("error reading file")
			}
			var profileData yaml.Profile
			if err = yaml.Unmarshal(data, &profileData); err != nil {
				log.Fatal().Err(err).Msg("error unmarshaling file contents")
			}
			log.Debug().Str("Key", profileData.Profile).Msg("adding profile")
			switch profileData.Profile {
			case "security":
				var securityProfile yaml.SecurityProfile
				if err = yaml.Unmarshal(data, &securityProfile); err != nil {
					log.Fatal().Err(err).Msg("error unmarshaling security profile")
				}
				jsonSecurityProfile, err := json.Marshal(securityProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling json")
				}
				if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
					log.Fatal().Msg("error adding profile, no part identifier given")
				}
				if profileData.CatalogID != "" {
					if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonSecurityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
				}
				if profileData.FVC != "" {
					log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
					uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by fvc")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
					break
				}
				if profileData.Sha256 != "" {
					log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
					uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by sha256")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonSecurityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added security profile to %s-%s\n", profileData.Name, profileData.Version)
				}
			case "licensing":
				var licensingProfile yaml.LicensingProfile
				if err = yaml.Unmarshal(data, &licensingProfile); err != nil {
					log.Fatal().Err(err).Msg("error unmarshaling licensing profile")
				}
				jsonLicensingProfile, err := json.Marshal(licensingProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling json")
				}
				if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
					log.Fatal().Msg("error adding profile, no part identifier given")
				}
				if profileData.CatalogID != "" {
					if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonLicensingProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
				}
				if profileData.FVC != "" {
					log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
					uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by fvc")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
					break
				}
				if profileData.Sha256 != "" {
					log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
					uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by sha256")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonLicensingProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added licensing profile to %s-%s\n", profileData.Name, profileData.Version)
				}
			case "quality":
				var qualityProfile yaml.QualityProfile
				if err = yaml.Unmarshal(data, &qualityProfile); err != nil {
					log.Fatal().Err(err).Msg("error unmarshaling quality profile")
				}
				jsonQualityProfile, err := json.Marshal(qualityProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("error marshaling json")
				}
				if profileData.CatalogID == "" && profileData.FVC == "" && profileData.Sha256 == "" {
					log.Fatal().Msg("error adding profile, no part identifier given")
				}
				if profileData.CatalogID != "" {
					if err = graphql.AddProfile(context.Background(), client, profileData.CatalogID, profileData.Profile, jsonQualityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
				}
				if profileData.FVC != "" {
					log.Debug().Str("File Verification Code", profileData.FVC).Msg("retrieving part id by file verification code")
					uuid, err := graphql.GetPartIDByFVC(context.Background(), client, profileData.FVC)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by fvc")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
					break
				}
				if profileData.Sha256 != "" {
					log.Debug().Str("SHA256", profileData.Sha256).Msg("retrieving part id by sha256")
					uuid, err := graphql.GetPartIDBySha256(context.Background(), client, profileData.Sha256)
					if err != nil {
						log.Fatal().Err(err).Msg("error retrieving part id by sha256")
					}
					if err = graphql.AddProfile(context.Background(), client, uuid.String(), profileData.Profile, jsonQualityProfile); err != nil {
						log.Fatal().Err(err).Msg("error adding profile")
					}
					fmt.Printf("Successfully added quality profile to %s-%s\n", profileData.Name, profileData.Version)
				}
			}
		}
		if argPartImportPath != "" {
			if argPartImportPath[len(argPartImportPath)-5:] != ".yaml" && argPartImportPath[len(argPartImportPath)-4:] != ".yml" {
				log.Fatal().Msg("error importing part, import path not a yaml file")
			}
			f, err := os.Open(argPartImportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error opening file")
			}
			defer f.Close()
			log.Debug().Msgf("successfully opened file at: %s", argPartImportPath)
			data, err := io.ReadAll(f)
			if err != nil {
				log.Fatal().Err(err).Msg("error reading file")
			}
			var partData yaml.Part
			if err = yaml.Unmarshal(data, &partData); err != nil {
				log.Fatal().Err(err).Msg("error unmarshaling file contents")
			}
			log.Debug().Msg("adding part")
			createdPart, err := graphql.AddPart(context.Background(), client, partData)
			if err != nil {
				log.Fatal().Err(err).Msg("error adding part")
			}
			prettyPart, err := json.MarshalIndent(&createdPart, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("Successfully added part from: %s\n", argPartImportPath)
			fmt.Printf("%s\n", string(prettyPart))

		}
	case "find":
		if err := findSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error parsing find subcommand flags")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if argSHA256 == "" && argPartID == "" && argFVC == "" && argSearchQuery == "" {
			findSubcommand.PrintDefaults()
			log.Fatal().Msg("error finding part, no data provided")
		}
		if argProfileType != "" {
			if argPartID == "" {
				log.Fatal().Msg("error getting profile, missing part id")
			}
			log.Debug().Str("ID", argPartID).Str("Key", argProfileType).Msg("retrieving profile")
			profile, err := graphql.GetProfile(context.Background(), client, argPartID, argProfileType)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving profile")
			}
			if len(*profile) < 1 {
				fmt.Println("No documents found")
				os.Exit(0)
			}
			prettyJson, err := json.MarshalIndent(&profile, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("%s\n", string(prettyJson))
			os.Exit(0)
		}
		if argSHA256 != "" {
			log.Debug().Str("SHA256", argSHA256).Msg("retrieving part id by sha256")
			partID, err := graphql.GetPartIDBySha256(context.Background(), client, argSHA256)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argFVC != "" {
			log.Debug().Str("File Verification Code", argFVC).Msg("retrieving part id by file verification code")
			partID, err := graphql.GetPartIDByFVC(context.Background(), client, argFVC)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part id")
			}
			fmt.Printf("Part ID: %s \n", partID.String())
		}
		if argSearchQuery != "" {
			log.Debug().Str("Query", argSearchQuery).Msg("executing part search")
			response, err := graphql.Search(context.Background(), client, argSearchQuery)
			if err != nil {
				log.Fatal().Err(err).Msg("error searching for part")
			}

			prettyJson, err := json.MarshalIndent(response, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("Result: %s\n", string(prettyJson))
		}
		if argPartID != "" {
			log.Debug().Str("ID", argPartID).Msg("retrieving part by id")
			response, err := graphql.GetPartByID(context.Background(), client, argPartID)
			if err != nil {
				log.Fatal().Err(err).Msg("error getting part by id")
			}
			prettyJson, err := json.MarshalIndent(&response, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("%s\n", string(prettyJson))
		}
	case "query":
		if err := querySubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error parsing query subcommand, query subcommand usage: ccli query <GraphQL Query>")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		argQuery := querySubcommand.Arg(0)
		if argQuery == "" {
			log.Fatal().Msg("error executing user query, query subcommand usage: ccli query <GraphQL Query>")
		}
		if argQuery != "" {
			log.Debug().Msg("executing raw graphql query")
			response, err := graphql.Query(context.Background(), client, argQuery)
			if err != nil {
				log.Fatal().Err(err).Msg("error querying graphql")
			}

			// json result will be output in prettified format
			var data map[string]interface{}
			json.Unmarshal(response, &data)

			prettyJson, err := json.MarshalIndent(data, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Println(string(prettyJson))
		}
	case "upload":
		if err := uploadSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error parsing upload subcommand, upload subcommand usage: ccli upload <Path>")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		argPath := uploadSubcommand.Arg(0)
		if argPath == "" {
			log.Fatal().Msg("error executing upload, upload subcommand usage: ccli upload <Path>")
		}
		if argPath != "" {
			log.Debug().Msg("uploading file to server")
			response, err := graphql.UploadFile(http.DefaultClient, configData.ServerAddr, argPath, "")
			if err != nil {
				log.Fatal().Err(err).Msg("error uploading archive")
			}
			if response.StatusCode == 200 {
				fmt.Printf("Successfully uploaded package: %s\n", argPath)
			}
		}
	case "update":
		if err := updateSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error updating part, update subcommand usage: ./ccli update -part <Path>")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if argImportPath == "" {
			log.Fatal().Msg("error updating part, update subcommand usage: ./ccli update -part <Path>")
		}
		if argImportPath != "" {
			if argImportPath[len(argImportPath)-5:] != ".yaml" && argImportPath[len(argImportPath)-4:] != ".yml" {
				log.Fatal().Msg("error importing part, import path not a yaml file")
			}
			f, err := os.Open(argImportPath)
			if err != nil {
				log.Fatal().Err(err).Msg("error opening file")
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				log.Fatal().Err(err).Msg("error reading file")
			}
			var partData yaml.Part
			if err = yaml.Unmarshal(data, &partData); err != nil {
				log.Fatal().Err(err).Msg("error decoding file contents")
			}
			log.Debug().Msg("updating part")
			returnPart, err := graphql.UpdatePart(context.Background(), client, &partData)
			if err != nil {
				log.Fatal().Err(err).Msg("error updating part")
			}
			prettyJson, err := json.MarshalIndent(&returnPart, "", indent)
			if err != nil {
				log.Fatal().Err(err).Msg("error prettifying json")
			}
			fmt.Printf("Part successfully updated\n%s\n", string(prettyJson))
		}
	case "delete":
		if err := deleteSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error deleting part, delete subcommand usage: ./ccli delete -id <catalog_id>")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if argPartID == "" {
			log.Fatal().Msg("error deleting part, delete subcommand usage: ./ccli delete -id <catalog_id>")
		}
		if argPartID != "" {
			log.Debug().Str("ID", argPartID).Msg("deleting part")
			if err := graphql.DeletePart(context.Background(), client, argPartID); err != nil {
				log.Fatal().Err(err).Msg("error deleting part from catalog")
			}
			fmt.Printf("Successfully deleted id: %s from catalog\n", argPartID)
		}
	case "ping":
		if err := pingSubcommand.Parse(os.Args[2:]); err != nil {
			log.Fatal().Err(err).Msg("error pinging server, ping subcommand usage: ./ccli ping")
		}
		if argVerboseMode {
			zerolog.SetGlobalLevel(0)
		}
		if configData.ServerAddr == "" {
			log.Fatal().Msg("invalid configuration file, no server address located")
		}
		log.Debug().Str("Address", configData.ServerAddr).Msg("pinging server")
		resp, err := http.DefaultClient.Get(configData.ServerAddr)
		if err != nil {
			log.Fatal().Err(err).Msg("error contacting server")
		}
		resp.Body.Close()
		if resp.StatusCode != 200 && resp.StatusCode != 422 {
			log.Fatal().Msgf("error reaching server, status code: %d", resp.StatusCode)
		} else {
			fmt.Println("Ping Result: Success")
		}
	default:
		printHelp()
	}
}
