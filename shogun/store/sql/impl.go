package sql

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/nwillems/shogun-graph/shogun/store"
	"github.com/nwillems/shogun-graph/shogun/store/sql/queries"
)

//go:embed migration/db_up.sql
var ddl string

type StorageImpl struct {
	db *sql.DB
	q  *queries.Queries
}

func NewStore(db *sql.DB) *StorageImpl {
	return &StorageImpl{
		db: db,
		q:  queries.New(db),
	}
}

func (s *StorageImpl) RunMigration(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	return nil
}

func (s *StorageImpl) Store(ctx context.Context, nodes []store.CreateNode, edges []store.CreateEdge) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	q := s.q.WithTx(tx)

	for _, node := range nodes {
		props, err := json.Marshal(node.Properties)
		if err != nil {
			return err
		}
		_, err = q.InsertNodes(ctx, queries.InsertNodesParams{
			NodeName:   node.Name,
			Properties: props,
		})
		if err != nil {
			return err
		}
	}

	for _, edge := range edges {
		props, err := json.Marshal(edge.Properties)
		if err != nil {
			return err
		}

		_, err = q.InsertEdges(ctx, queries.InsertEdgesParams{
			Source:     edge.Source,
			Target:     edge.Target,
			Type:       edge.Type,
			Properties: props,
		})
		if err != nil {
			return nil
		}
	}

	return tx.Commit()
}

func (s *StorageImpl) Query(ctx context.Context, q store.Query) (*store.Graph, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *StorageImpl) QueryNode(ctx context.Context, id int64) (*store.Graph, error) {
	return nil, fmt.Errorf("Not implemented")
}
