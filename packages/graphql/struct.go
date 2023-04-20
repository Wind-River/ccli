// This package implements graphql query and mutation data structures and handling for ccli utilizing hasura go-graphql-client library
package graphql

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Archive struct {
	Sha256     string    `graphql:"sha256"`
	Size       int64     `graphql:"Size"`
	PartID     uuid.UUID `graphql:"part_id"`
	Part       Part
	Md5        string `graphql:"md5"`
	Sha1       string `graphql:"sha1"`
	Name       string `graphql:"name"`
	InsertDate string `graphql:"insert_date"`
}

type Part struct {
	ID       uuid.UUID `graphql:"id"`
	PartType string    `graphql:"type"`
	Version  string    `graphql:"version"`
	Name     string    `graphql:"name"`
	//Label                string    `graphql:"label"`
	FamilyName           string `graphql:"family_name"`
	FileVerificationCode string `graphql:"file_verification_code"`
	Size                 int64  `graphql:"size"`
	License              string `graphql:"license"`
	LicenseRationale     string `graphql:"license_rationale"`
	//Description          string    `graphql:"description"`
	Comprised string    `graphql:"comprised"`
	Aliases   []string  `graphql:"aliases"`
	Profiles  []Profile `graphql:"profiles"`
}

type Profile struct {
	Key       string     `graphql:"key"`
	Documents []Document `graphql:"documents"`
}

type Document struct {
	Title    string          `graphql:"title"`
	Document json.RawMessage `graphql:"document"`
}
