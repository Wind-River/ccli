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

// This package implements graphql query and mutation data structures and handling for ccli utilizing hasura go-graphql-client library
package graphql

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

// Required to query graphql using custom scalar
type UUID string
type JSON string
type Upload os.File

type Archive struct {
	Sha256     string    `graphql:"sha256"`
	Size       int64     `graphql:"size"`
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
	ContentType          string    `graphql:"type" yaml:"content_type"`
	Version              string    `graphql:"version" yaml:"version"`
	Name                 string    `graphql:"name" yaml:"name"`
	Label                string    `graphql:"label" yaml:"label"`
	FamilyName           string    `graphql:"family_name" yaml:"family_name"`
	FileVerificationCode string    `graphql:"file_verification_code" yaml:"file_verification_code"`
	Size                 int64     `graphql:"size" yaml:"size"`
	License              string    `graphql:"license" yaml:"license"`
	LicenseRationale     string    `graphql:"license_rationale" yaml:"license_rationale"`
	Description          string    `graphql:"description" yaml:"description"`
	HomePage             string    `graphql:"home_page" yaml:"home_page"`
	Comprised            uuid.UUID `graphql:"comprised" yaml:"comprised"`
	Aliases              []string  `graphql:"aliases" yaml:"aliases"`
}

type PartInput struct {
	ID                   *UUID  `graphql:"id" json:"id"`
	Type                 string `graphql:"type" json:"type"`
	ContentType          string `graphql:"type" json:"content_type"`
	Name                 string `graphql:"name" json:"name"`
	Version              string `graphql:"version" json:"version"`
	Label                string `graphql:"label" json:"label"`
	FamilyName           string `graphql:"family_name" json:"family_name"`
	FileVerificationCode string `graphql:"file_verification_code" json:"file_verification_code"`
	License              string `graphql:"license" json:"license"`
	LicenseRationale     string `graphql:"license_rationale" json:"license_rationale"`
	Description          string `graphql:"description" json:"description"`
	HomePage             string `graphql:"home_page" json:"home_page"`
	Comprised            *UUID  `graphql:"comprised" json:"comprised"`
}

type NewPartInput struct {
	Type             string `graphql:"type" json:"type"`
	ContentType      string `graphql:"type" json:"content_type"`
	Name             string `graphql:"name" json:"name"`
	Version          string `graphql:"version" json:"version"`
	Label            string `graphql:"label" json:"label"`
	FamilyName       string `graphql:"family_name" json:"family_name"`
	License          string `graphql:"license" json:"license"`
	LicenseRationale string `graphql:"license_rationale" json:"license_rationale"`
	Description      string `graphql:"description" json:"description"`
	HomePage         string `graphql:"home_page" json:"home_page"`
	Comprised        *UUID  `graphql:"comprised" json:"comprised"`
}

type Profile []Document

type Document struct {
	Title    string          `graphql:"title" json:"title"`
	Document json.RawMessage `graphql:"document" json:"document"`
}
