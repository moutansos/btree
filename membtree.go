package btree

import (
	"bytes"
	"encoding/binary"
)

// BTreeInMemory is a structure that references a b-tree structure that
// resides in memory instead of on disk.
type BTreeInMemory struct {
	data []byte
}

func NewBTreeInMem(size uint64) (*BTreeInMemory, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, size)
	if err != nil {
		return nil, err
	}

	tree := new(BTreeInMemory)
	tree.data = appendRangeBytes(tree.data, buf.Bytes())
	return tree, nil
}

func appendRangeBytes(d []byte, n []byte) []byte {
	for _, b := range n {
		d = append(d, b)
	}
	return d
}
