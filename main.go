package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	graphqlController "wrs/catalog/ccli/packages/graphql"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed data/busybox-1.35.0.json
//var testData []byte

// var argLogLevel int
var argEndpoint string
var argLogLevel int
var sha256 string
var query string

func init() {
	flag.StringVar(&argEndpoint, "endpoint", "https://aws-ip-services-tk-dev.wrs.com/api/graphql", "GraphQL Endpoint")
	flag.IntVar(&argLogLevel, "log", int(zerolog.InfoLevel), "Log Level Minimum")
	flag.StringVar(&sha256, "sha256", "", "SHA256 Part Query")
	flag.StringVar(&query, "query", "", "Graphql Query")
}

func main() {
	flag.Parse()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	client := graphql.NewClient(argEndpoint, httpClient)

	command := os.Args[1]
	switch command {
	case "export":
		fmt.Println("Exporting")
	case "add":
		fmt.Println("Adding")
	default:
		if sha256 != "" {
			partID, err := graphqlController.GetPartID(context.Background(), client, sha256)
			if err != nil {
				log.Fatal().Err(err).Msg("error retrieving part id")
			}
			// log.Info().Str("Part ID", partID.String()).Send()
			fmt.Printf("Part ID: %s \n", partID.String())
		}

		if query != "" {
			response, err := graphqlController.Query(context.Background(), client, query)
			if err != nil {
				log.Fatal().Err(err).Msg("error querying graphql")
			}
			fmt.Printf("Query Response: %s\n", response)
		}
	}
}
