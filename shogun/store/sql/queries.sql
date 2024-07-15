

-- name: GetNode :one
SELECT n.id, n.node_name, n.properties
FROM nodes n
WHERE n.id = sqlc.arg('node_id');

-- name: GetNodeEdges :many
SELECT e.*
FROM edges e
WHERE e.source = sqlc.arg('node_id') OR e.target = sqlc.arg('node_id');

-- name: InsertNodes :one
INSERT INTO nodes(node_name, properties )
VALUES (sqlc.arg('node_name'), sqlc.arg('properties'))
RETURNING *;

-- name: InsertEdges :one
INSERT INTO edges(source, target, type, properties)
SELECT source_node.id, target_node.id, sqlc.arg('type'), sqlc.arg('properties')
FROM (
  SELECT 'A' jj, id
  FROM nodes sn
  WHERE sn.node_name = sqlc.arg('source')
  LIMIT 1
) AS source_node
LEFT JOIN (
  SELECT 'A' jj, id
  FROM nodes tn
  WHERE tn.node_name = sqlc.arg('target')
  LIMIT 1
) AS target_node ON source_node.jj = target_node.jj
RETURNING *;

