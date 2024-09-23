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
package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"wrs/catalog/ccli/packages/yaml"

	graphqlUpload "bitbucket.wrs.com/scm/weststar/graphql-upload-go.git"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
)

// Adds a profile document to a part and returns any errors that occur
func AddProfile(ctx context.Context, client *graphql.Client, id string, key string, document json.RawMessage) error {
	var mutation struct {
		AttachDocument bool `graphql:"attachDocument(id: $id, key: $key, document: $document)"`
	}

	variables := map[string]interface{}{
		"id":       UUID(id),
		"key":      key,
		"document": JSON(document),
	}
	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}
	return nil
}

// Delete part from catalog - TODO - this should be difficult to execute
func DeletePart(ctx context.Context, client *graphql.Client, id string, recursiveDelete bool, forceDelete bool) error {
	var mutation struct {
		DeletePart bool `graphql:"deletePart(part_id: $id, recursive:$recursiveDelete, force:$forceDelete)"`
	}

	variables := map[string]interface{}{
		"id":              UUID(id),
		"recursiveDelete": recursiveDelete,
		"forceDelete":     forceDelete,
	}

	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}
	return nil
}

// Retrieves a profile from the catalog
func GetProfile(ctx context.Context, client *graphql.Client, id string, key string) (*Profile, error) {
	var query struct {
		Profile `graphql:"profile(id:$id, key:$key)"`
	}

	variables := map[string]interface{}{
		"id":  UUID(id),
		"key": key,
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.Profile, nil
}

// Adds a logical part to the catalog using a yaml template format and returns the inserted part
func AddPart(ctx context.Context, client *graphql.Client, newPart yaml.Part) (*Part, error) {
	var newPartInput NewPartInput

	if err := YamlToNewPartInput(newPart, &newPartInput); err != nil {
		return nil, err
	}

	var mutation struct {
		Part `graphql:"createPart(partInput: $partInput)"`
	}

	variables := map[string]interface{}{
		"partInput": newPartInput,
	}

	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	//Alias insertion is handling with createAlias mutation
	if newPart.Aliases != nil && len(newPart.Aliases) != 0 {
		var aliasMutation struct {
			UUID `graphql:"createAlias(id: $id, alias: $alias)"`
		}

		for _, v := range newPart.Aliases {
			aliasVariables := map[string]interface{}{
				"id":    UUID(mutation.Part.ID.String()),
				"alias": v,
			}

			if err := client.Mutate(ctx, &aliasMutation, aliasVariables); err != nil {
				return nil, err
			}

		}
	}
	// Subparts are inserted utilizing partHasPart mutation
	if newPart.CompositeList != nil && len(newPart.CompositeList) != 0 {
		var compositeMutation struct {
			PartHasPart bool `graphql:"partHasPart(parent: $parent, child: $child, path: $path)"`
		}

		// Seen map prevents duplication of subpart paths in the catalog
		seen := make(map[string]bool)
		compositeList := []string{}

		for _, v := range newPart.CompositeList {
			if !seen[v] {
				compositeList = append(compositeList, v)
				seen[v] = true
			}
		}

		for _, v := range compositeList {
			compositeVariables := map[string]interface{}{
				"parent": UUID(mutation.ID.String()),
				"child":  UUID(v),
				"path":   v,
			}

			if err := client.Mutate(ctx, &compositeMutation, compositeVariables); err != nil {
				return nil, err
			}
		}
	}

	return &mutation.Part, nil
}

// retrieve part data from provided sha256 value
func GetPartIDBySha256(ctx context.Context, client *graphql.Client, sha256 string) (*uuid.UUID, error) {
	var query struct {
		Archive `graphql:"archive(sha256: $sha256)"`
	}

	variables := map[string]interface{}{
		"sha256": sha256,
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.PartID, nil
}

// retrieve part data from provided file verification code value
func GetPartIDByFVC(ctx context.Context, client *graphql.Client, fvc string) (*uuid.UUID, error) {
	var query struct {
		Part `graphql:"part(file_verification_code: $fvc)"`
	}

	variables := map[string]interface{}{
		"fvc": fvc,
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.Part.ID, nil
}

// Retrieves a part from the catalog using catalog id
func GetPartByID(ctx context.Context, client *graphql.Client, id string) (*Part, error) {

	var query struct {
		Part `graphql:"part(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": UUID(id),
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	return &query.Part, nil
}

// Retrieves a part from the catalog using file verification code
func GetPartByFVC(ctx context.Context, client *graphql.Client, fvc string) (*Part, error) {

	var query struct {
		Part `graphql:"part(file_verification_code: $fvc)"`
	}

	variables := map[string]interface{}{
		"fvc": fvc,
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	return &query.Part, nil
}

// Retrieves a part from the catalog using sha256
func GetPartBySHA256(ctx context.Context, client *graphql.Client, sha256 string) (*Part, error) {

	var query struct {
		Part `graphql:"part(sha256: $sha256)"`
	}

	variables := map[string]interface{}{
		"sha256": sha256,
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}
	return &query.Part, nil
}

// Retrieves a slice of parts from the catalog using find_archive query to search by name
func Search(ctx context.Context, client *graphql.Client, searchQuery string) (*[]Part, error) {

	var query struct {
		FindArchive []struct {
			Archive `graphql:"archive"`
		} `graphql:"find_archive(query: $searchQuery, method: $method)"`
	}

	variables := map[string]interface{}{
		"searchQuery": searchQuery,
		"method":      "fast",
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	var parts []Part
	for _, v := range query.FindArchive {
		parts = append(parts, v.Part)
	}
	return &parts, nil
}

// allow user defined queries to be executed by ccli
func Query(ctx context.Context, client *graphql.Client, query string) ([]byte, error) {
	response, err := client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// uploads an archive to the catalog using graphql-upload library
func UploadFile(httpClient *http.Client, uri string, path string, name string) (*graphqlUpload.Response, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	response, err := graphqlUpload.Upload(
		http.DefaultClient,
		uri,
		`
		mutation($file: Upload!){
		  uploadArchive(file: $file){
			  name
			  insert_date
			  sha256
			  sha1
			  part_id
		  }
		}
	  `,
		graphqlUpload.File{
			Name:     path,
			Variable: "file",
			Data:     f,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// updates a part record from the catalog using yaml template
func UpdatePart(ctx context.Context, client *graphql.Client, partData *yaml.Part) (*Part, error) {

	var partInput PartInput
	if partData.CatalogID == "" {
		if partData.FVC == "" && partData.Sha256 == "" {
			return nil, errors.New("error updating part, no part identifier provided")
		}
		if partData.FVC != "" {
			catalogID, err := GetPartIDByFVC(ctx, client, partData.FVC)
			if err != nil {
				return nil, err
			}
			if catalogID != nil {
				partInputID := UUID(catalogID.String())
				partInput.ID = &partInputID
			}
		} else if partData.Sha256 != "" {
			catalogID, err := GetPartIDBySha256(ctx, client, partData.Sha256)
			if err != nil {
				return nil, err
			}
			if catalogID != nil {
				partInputID := UUID(catalogID.String())
				partInput.ID = &partInputID
			}
		}
	}
	if partData.CatalogID != "" {
		partInputID := UUID(partData.CatalogID)
		partInput.ID = &partInputID
	}
	if partData.Name != "" {
		partInput.Name = partData.Name
	}
	if partData.Version != "" {
		partInput.Version = partData.Version
	}
	if partData.FamilyName != "" {
		partInput.FamilyName = partData.FamilyName
	}
	if partData.ContentType != "" {
		partInput.ContentType = partData.ContentType
	}
	if partData.Type != "" {
		partInput.Type = partData.Type
	}
	if partData.Label != "" {
		partInput.Label = partData.Label
	}
	if partData.License.LicenseExpression != "" {
		partInput.License = partData.License.LicenseExpression
	}
	if partData.License.AnalysisType != "" {
		partInput.LicenseRationale = partData.License.AnalysisType
	}
	if partData.Description != "" {
		partInput.Description = partData.Description
	}
	if partData.HomePage != "" {
		partInput.HomePage = partData.HomePage
	}
	if partData.ComprisedOf != "" {
		comprisedID := UUID(partData.ComprisedOf)
		partInput.Comprised = &comprisedID
	}

	var mutation struct {
		Part `graphql:"updatePart(partInput: $partInput)"`
	}

	variables := map[string]interface{}{
		"partInput": partInput,
	}

	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if partData.Aliases != nil {
		var aliasMutation struct {
			UUID `graphql:"createAlias(id: $id, alias: $alias)"`
		}

		for _, v := range partData.Aliases {
			aliasVariables := map[string]interface{}{
				"id":    *partInput.ID,
				"alias": v,
			}

			if err := client.Mutate(ctx, &aliasMutation, aliasVariables); err != nil {
				return nil, err
			}
		}
	}

	return &mutation.Part, nil
}

// Sets the fields of a part record including emptry values from the catalog using yaml template
func SetPart(ctx context.Context, client *graphql.Client, partData *yaml.Part) (*Part, error) {

	var partInput PartInput
	if partData.CatalogID == "" {
		if partData.FVC == "" && partData.Sha256 == "" {
			return nil, errors.New("error updating part, no part identifier provided")
		}
		if partData.FVC != "" {
			catalogID, err := GetPartIDByFVC(ctx, client, partData.FVC)
			if err != nil {
				return nil, err
			}
			if catalogID != nil {
				partInputID := UUID(catalogID.String())
				partInput.ID = &partInputID
			}
		} else if partData.Sha256 != "" {
			catalogID, err := GetPartIDBySha256(ctx, client, partData.Sha256)
			if err != nil {
				return nil, err
			}
			if catalogID != nil {
				partInputID := UUID(catalogID.String())
				partInput.ID = &partInputID
			}
		}
	}

	if partData.ComprisedOf != "" {
		comprisedID := UUID(partData.ComprisedOf)
		partInput.Comprised = &comprisedID
	}

	var mutation struct {
		Part `graphql:"setPart(partInput: $partInput)"`
	}

	variables := map[string]interface{}{
		"partInput": partInput,
	}

	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if partData.Aliases != nil {
		var aliasMutation struct {
			UUID `graphql:"createAlias(id: $id, alias: $alias)"`
		}

		for _, v := range partData.Aliases {
			aliasVariables := map[string]interface{}{
				"id":    *partInput.ID,
				"alias": v,
			}

			if err := client.Mutate(ctx, &aliasMutation, aliasVariables); err != nil {
				return nil, err
			}
		}
	}

	return &mutation.Part, nil
}

// Used to convert a part data structure into the structure expected by yaml i/o
func UnmarshalPart(part *Part, yamlPart *yaml.Part) error {
	yamlPart.Format = 1.0
	yamlPart.CatalogID = part.ID.String()
	yamlPart.FVC = part.FileVerificationCode
	yamlPart.Name = part.Name
	yamlPart.Version = part.Version
	yamlPart.FamilyName = part.FamilyName
	yamlPart.Type = part.PartType
	yamlPart.ContentType = part.ContentType
	yamlPart.Label = part.Label
	yamlPart.Description = part.Description
	yamlPart.HomePage = part.HomePage
	yamlPart.License.LicenseExpression = part.License
	yamlPart.License.AnalysisType = part.LicenseRationale
	if part.Size != 0 {
		yamlPart.Size = fmt.Sprint(part.Size)
	}
	yamlPart.Aliases = part.Aliases
	if part.Comprised != uuid.Nil {
		yamlPart.ComprisedOf = part.Comprised.String()
	}
	return nil
}

// Used to convert a part in yaml format into the format expected for new part mutation
func YamlToNewPartInput(yamlPart yaml.Part, newPartInput *NewPartInput) error {
	newPartInput.Type = yamlPart.Type
	newPartInput.ContentType = yamlPart.ContentType
	newPartInput.Name = yamlPart.Name
	newPartInput.Version = yamlPart.Version
	newPartInput.Label = yamlPart.Label
	newPartInput.FamilyName = yamlPart.FamilyName
	newPartInput.License = yamlPart.License.LicenseExpression
	newPartInput.LicenseRationale = yamlPart.License.AnalysisType
	newPartInput.Description = yamlPart.Description
	newPartInput.HomePage = yamlPart.HomePage
	if yamlPart.ComprisedOf != "" {
		comprisedUUID, err := uuid.Parse(yamlPart.ComprisedOf)
		if err != nil {
			return err
		}
		if comprisedUUID != uuid.Nil {
			graphqlUUID := UUID(comprisedUUID.String())
			newPartInput.Comprised = &graphqlUUID
		}
	}
	return nil
}
