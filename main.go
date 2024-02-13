package main

import (
	_ "embed"
	"fmt"
	"os"

	"wrs/catalog/ccli/packages/cmd"
	"wrs/catalog/ccli/packages/config"
	"wrs/catalog/ccli/packages/graphql"
	"wrs/catalog/ccli/packages/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var configFile config.ConfigData
var indent string

// initialize configuration file and flag values
func init() {
	/*configFile, err := os.Open("ccli_config.yml")
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
	*/
	viper.SetConfigFile("ccli_config.yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("User configuration file not found. Please create ccli_config.yml and copy the contents of ccli_config.DEFAULT.yml.")
			os.Exit(1)
		} else {
			os.Exit(1)
		}
	}
	if err := viper.Unmarshal(&configFile); err != nil {
		fmt.Println("Could not unmarshal config file parameters")
		os.Exit(1)
	}
	indentString := ""
	for i := 0; i < int(configFile.JsonIndent); i++ {
		indentString += " "
	}
	indent = indentString
}

func main() {
	// set global log level to value found in configuration file
	zerolog.SetGlobalLevel(zerolog.Level(configFile.LogLevel))

	if configFile.LogFile == "" || configFile.LogFile[len(configFile.LogFile)-4:] != ".txt" {
		fmt.Println("*** ERROR - Error reading config file, log file must be a .txt file")
		os.Exit(1)
	}
	// open log file and set logging output
	logFile, err := os.Create(configFile.LogFile)
	if err != nil {
		fmt.Printf("*** ERROR - Error opening log file: %s\n", err.Error())
	}
	defer logFile.Close()
	multiLogger := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	log.Logger = log.Output(multiLogger)
	log.Debug().Msgf("log file opened at: %s", configFile.LogFile)

	if configFile.ServerAddr == "" {
		log.Fatal().Msg("invalid configuration file, no server address located")
	}
	resp, err := http.DefaultClient.Get(configFile.ServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("error contacting server")
	}
	resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 422 {
		log.Fatal().Msgf("server connection error, check config file and network configuration: Status Code (%d)\n", resp.StatusCode)
	}
	log.Debug().Msgf("successfully connected to server")

	client := graphql.GetNewClient(configFile.ServerAddr, http.DefaultClient)

	rootCmd := cmd.Example(&configFile)
	rootCmd.AddCommand(cmd.Ping(&configFile))
	rootCmd.AddCommand(cmd.Upload(&configFile))
	rootCmd.AddCommand(cmd.Update(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Query(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Find(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Export(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Add(&configFile, client, indent))
	rootCmd.AddCommand(cmd.Delete(&configFile, client, indent))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("*** ERROR:", err)
		os.Exit(1)
	}
}
