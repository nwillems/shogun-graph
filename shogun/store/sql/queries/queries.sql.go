// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package queries

import (
	"context"
	"encoding/json"
)

const getNode = `-- name: GetNode :one
SELECT n.id, n.node_name, n.properties
FROM nodes n
WHERE n.id = ?1
`

func (q *Queries) GetNode(ctx context.Context, nodeID int64) (Node, error) {
	row := q.db.QueryRowContext(ctx, getNode, nodeID)
	var i Node
	err := row.Scan(&i.ID, &i.NodeName, &i.Properties)
	return i, err
}

const insertEdges = `-- name: InsertEdges :one
INSERT INTO edges(source, target, type, properties)
SELECT source_node.id, target_node.id, ?1, ?2
FROM (
  SELECT 'A' jj, id
  FROM nodes sn
  WHERE sn.node_name = ?3
  LIMIT 1
) AS source_node
LEFT JOIN (
  SELECT 'A' jj, id
  FROM nodes tn
  WHERE tn.node_name = ?4
  LIMIT 1
) AS target_node ON source_node.jj = target_node.jj
RETURNING id, source, target, type, properties
`

type InsertEdgesParams struct {
	Type       string
	Properties json.RawMessage
	Source     string
	Target     string
}

func (q *Queries) InsertEdges(ctx context.Context, arg InsertEdgesParams) (Edge, error) {
	row := q.db.QueryRowContext(ctx, insertEdges,
		arg.Type,
		arg.Properties,
		arg.Source,
		arg.Target,
	)
	var i Edge
	err := row.Scan(
		&i.ID,
		&i.Source,
		&i.Target,
		&i.Type,
		&i.Properties,
	)
	return i, err
}

const insertNodes = `-- name: InsertNodes :one
INSERT INTO nodes(node_name, properties )
VALUES (?1, ?2)
RETURNING id, node_name, properties
`

type InsertNodesParams struct {
	NodeName   string
	Properties json.RawMessage
}

func (q *Queries) InsertNodes(ctx context.Context, arg InsertNodesParams) (Node, error) {
	row := q.db.QueryRowContext(ctx, insertNodes, arg.NodeName, arg.Properties)
	var i Node
	err := row.Scan(&i.ID, &i.NodeName, &i.Properties)
	return i, err
}