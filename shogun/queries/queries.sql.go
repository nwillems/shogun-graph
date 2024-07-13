// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package queries

import (
	"context"
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