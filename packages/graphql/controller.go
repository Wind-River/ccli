package graphql

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
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
