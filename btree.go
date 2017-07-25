package btree

type BTree interface {
	NewNode()
	GetBlockSize() uint64
	WriteNode(n *Node) error
}
