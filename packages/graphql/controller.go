package graphqlController

import (
	"context"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

func GetPartID(ctx context.Context, client *graphql.Client, sha256 string) (*uuid.UUID, error) {
	var query struct {
		Archive struct {
			PartID uuid.UUID `graphql:"part_id"`
		} `graphql:"archive(sha256: $sha256)"`
	}

	variables := map[string]interface{}{
		"sha256": sha256,
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}
	return &query.Archive.PartID, nil
}

func Query(ctx context.Context, client *graphql.Client, query string) ([]byte, error) {
	response, err := client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}
