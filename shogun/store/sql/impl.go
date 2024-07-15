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

func nodeToData(node queries.Node) (store.DataNode, error) {
	var props map[string]any
	err := json.Unmarshal(node.Properties, &props)

	return store.DataNode{
		ID:   node.ID,
		Name: node.NodeName,

		Properties: props,
	}, err
}

func edgeToData(edge queries.Edge) (store.DataEdge, error) {
	var props map[string]any
	err := json.Unmarshal(edge.Properties, &props)

	return store.DataEdge{
		Source: edge.Source,
		Target: edge.Target,
		Type:   edge.Type,

		Properties: props,
	}, err
}

func (s *StorageImpl) Query(ctx context.Context, q store.Query) (*store.Graph, error) {
	return nil, fmt.Errorf("Not implemented")
}

func Map[S, T any](ss []S, fs func(S) T) []T {
	result := make([]T, len(ss))
	for i, s := range ss {
		result[i] = fs(s)
	}

	return result
}

func FoldL[S, Acc any](xs []S, fs func(S, Acc) Acc, init Acc) Acc {
	result := init
	for _, x := range xs {
		result = fs(x, result)
	}

	return result
	// return FoldL(xs[1:], fs, fs(x[0])) // Missing bounds check
}

func MapErr[S, T any](ss []S, fse func(S) (T, error)) ([]T, error) {
	type ReturnVal struct {
		val T
		err error
	}
	type ReturnVal2 struct {
		val []T
		err error
	}

	ww := func(s S) func() ReturnVal { return func() ReturnVal { v, err := fse(s); return ReturnVal{v, err} } }
	t := Map(ss, ww)

	res := FoldL(t, func(f func() ReturnVal, acc ReturnVal2) ReturnVal2 {
		if acc.err != nil { // Short on error
			return acc
		}

		ret := f()
		return ReturnVal2{append(acc.val, ret.val), ret.err}
	}, ReturnVal2{make([]T, 0, len(t)), nil})

	return res.val, res.err
}

func (s *StorageImpl) QueryNode(ctx context.Context, id int64) (*store.Graph, error) {
	node, err := s.q.GetNode(ctx, id)
	if err != nil {
		return nil, err
	}

	dataNode, err := nodeToData(node)
	if err != nil {
		return nil, err
	}

	edges, err := s.q.GetNodeEdges(ctx, id)
	if err != nil {
		return nil, err
	}

	dataEdges, err := MapErr(edges, edgeToData)

	return &store.Graph{
		Nodes: []store.DataNode{dataNode},
		Edges: dataEdges,
	}, err
}
