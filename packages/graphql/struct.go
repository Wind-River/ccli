// This package implements graphql query and mutation data structures and handling for ccli utilizing hasura go-graphql-client library
package graphql

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Required to query graphql using custom scalar
type UUID string
type JSON string

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
	ID                   uuid.UUID `graphql:"id" yaml:"id"`
	PartType             string    `graphql:"type" yaml:"type"`
	Version              string    `graphql:"version" yaml:"version"`
	Name                 string    `graphql:"name" yaml:"name"`
	Label                string    `graphql:"label" yaml:"label"`
	FamilyName           string    `graphql:"family_name" yaml:"family_name"`
	FileVerificationCode string    `graphql:"file_verification_code" yaml:"file_verification_code"`
	Size                 int64     `graphql:"size" yaml:"size"`
	License              string    `graphql:"license" yaml:"license"`
	LicenseRationale     string    `graphql:"license_rationale" yaml:"license_rationale"`
	Description          string    `graphql:"description" yaml:"description"`
	Comprised            uuid.UUID `graphql:"comprised" yaml:"comprised"`
	Aliases              []string  `graphql:"aliases" yaml:"aliases"`
}

type Profile struct {
	Key       string     `graphql:"key"`
	Documents []Document `graphql:"documents"`
}

type Document struct {
	Title    string          `graphql:"title"`
	Document json.RawMessage `graphql:"document"`
}
