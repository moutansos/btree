package btree

type BTree interface {
	NewNode()
	GetBlockSize() uint64
}
