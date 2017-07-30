package btree

type BTree interface {
	NewNode() (n *Node, err error)
	WriteNode(n *Node) error
	ReadNode(address int64) (n *Node, err error)
}
