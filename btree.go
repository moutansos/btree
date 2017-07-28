package btree

type BTree interface {
	NewNode() (n *Node, err error)
	WriteNode(n *Node) error
}
