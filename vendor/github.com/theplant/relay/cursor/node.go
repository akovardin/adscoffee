package cursor

type Node[T any] interface {
	RelayNode() T
}

type SelfNode[T any] struct {
	Node T
}

func (n *SelfNode[T]) RelayNode() T {
	return n.Node
}

func (n *SelfNode[T]) MarshalJSON() ([]byte, error) {
	return JSONMarshal(n.Node)
}

type NodeWrapper[T any] struct {
	Object any
	Unwrap func() T
}

func (n *NodeWrapper[T]) RelayNode() T {
	return n.Unwrap()
}

func (n *NodeWrapper[T]) MarshalJSON() ([]byte, error) {
	return JSONMarshal(n.Object)
}
