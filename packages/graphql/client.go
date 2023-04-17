package graphql

import (
	"github.com/hasura/go-graphql-client"
)

// port hasura client generation into ccli graphql package
func GetNewClient(url string, httpClient graphql.Doer) *graphql.Client {
	return graphql.NewClient(url, httpClient)
}
