package btree

// BTree is an interface into a b-tree collection.
// There are two types of b-trees available in this library. The BTreeOnDisk and
// the BTreeInMemory. These are accessible by this interface.
type BTree interface {
	NewNode() (n *Node, err error)
	WriteNode(n *Node) error
	ReadNode(address int64) (n *Node, err error)
}