package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"flag"
	"net/http"
	"wrs/catalog/cli/graphqlController"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed data/busybox-1.35.0.json
var testData []byte

// var argLogLevel int
var argEndpoint string
var argLogLevel int
var sha256 string

func init() {
	flag.StringVar(&argEndpoint, "endpoint", "https://aws-ip-services-tk-dev.wrs.com/api/graphql", "GraphQL Endpoint")
	flag.IntVar(&argLogLevel, "log", int(zerolog.InfoLevel), "Log Level Minimum")
	flag.StringVar(&sha256, "part", "", "Graphql Part Query")
}

func main() {
	flag.Parse()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	client := graphql.NewClient(argEndpoint, httpClient)

	if sha256 != "" {
		partID, err := graphqlController.GetPartID(context.Background(), client, sha256)
		if err != nil {
			log.Fatal().Err(err).Msg("error retrieving part id")
		}
	}

}
