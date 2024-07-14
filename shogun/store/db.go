package store

import "context"

type (
	NodeName       = string
	JsonProperties = map[string]any
)

type (
	CreateNode struct {
		Name       string
		Properties JsonProperties
	}
	CreateEdge struct {
		Source NodeName
		Target NodeName

		Type string

		Properties JsonProperties
	}
)

type (
	DataNodeRef = int64 // Could also be the name
	Query       struct {
		Nodes     []string
		EdgeTypes *([]string)
		Depth     *int
	}

	Graph struct {
		Nodes []DataNode
		Edges []DataEdge
	}
	DataEdge struct {
		Type       string
		Properties JsonProperties

		Source DataNodeRef
		Target DataNodeRef
	}
	DataNode struct {
		ID   int64
		Name string

		Properties JsonProperties
	}
)

type Storage interface {
	// In a transaction, creates the given nodes, and then creates the edges specifed
	Store(ctx context.Context, nodes []CreateNode, edges []CreateEdge) error

	// Queries the graph, using the given query
	Query(ctx context.Context, q Query) (*Graph, error)
	// Queries a single node, returns a graph with thaht single node and it's edges(ones with the node as source)
	QueryNode(ctx context.Context, id int64) (*Graph, error)
}
