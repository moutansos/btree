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

func insertUint64at(ara []uint64, i int, val uint64) []uint64 {
	// https://github.com/golang/go/wiki/SliceTricks
	ara = append(ara, 0)
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}
