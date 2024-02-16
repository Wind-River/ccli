// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.

// Provides access to hasura go graphql client and implements graphql query and mutation functionality.
package graphql

import (
	"github.com/hasura/go-graphql-client"
)

// port hasura client generation into ccli graphql package
func GetNewClient(url string, httpClient graphql.Doer) *graphql.Client {
	return graphql.NewClient(url, httpClient)
}
