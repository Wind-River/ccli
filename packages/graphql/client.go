// Provides access to hasura go graphql client and implements graphql query and mutation functionality.
package graphql

import (
	"github.com/hasura/go-graphql-client"
)

// port hasura client generation into ccli graphql package
func GetNewClient(url string, httpClient graphql.Doer) *graphql.Client {
	return graphql.NewClient(url, httpClient)
}
