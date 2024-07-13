package shogun

type (
	Node struct {
		Name string
		Properties JsonThing
	}
	Edge struct {
		Left NodeID
		Right NodeID

		Type string

		Properties JsonThing
	}
)

type Storage interface {
	Store(nodes []Node, edges []Edge) error
}
