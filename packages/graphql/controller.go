package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"wrs/catalog/ccli/packages/yaml"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	graphqlUpload "gitlab.devstar.cloud/WestStar/libraries/go/graphql-upload.git/code"
)

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

func UploadFile(httpClient *http.Client, uri string, path string, name string) (*http.Response, error) {
	return graphqlUpload.UploadFile(httpClient, uri, path, name)
}

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

func UnmarshalPart(part *Part, yamlPart *yaml.Part) error {
	yamlPart.Format = 1.0
	yamlPart.CatalogID = part.ID.String()
	yamlPart.FVC = part.FileVerificationCode
	yamlPart.Name = part.Name
	yamlPart.Version = part.Version
	yamlPart.Label = part.Label
	yamlPart.Description = part.Description
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
