package graphqlController

import (
	"context"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

func getPartID(ctx context.Context, client graphql.Client, sha256 string) (*uuid.UUID, error) {
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
