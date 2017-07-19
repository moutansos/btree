package btree

type Node struct {
	pointers []uint64
	data	 []uint64
	tree	 *BTree
}

func NewNode(t *BTree) (*Node, error) {
	n := new(Node)
	n.tree = t
	return n, nil
}
