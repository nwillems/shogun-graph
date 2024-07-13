
-- name: GetNode :one
SELECT n.id, n.node_name, n.properties
FROM nodes n
WHERE n.id = sqlc.arg('node_id')
