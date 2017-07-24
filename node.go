package btree

import (
	"bytes"
	"encoding/binary"
)

type Node struct {
	Pointers [32]uint64
	Data     [31]uint64
	tree     *BTree
}

type binaryNode struct {
	Pointers [32]uint64
	Data     [31]uint64
}

func NewNode(t BTree) (*Node, error) {
	n := new(Node)
	n.tree = &t
	return n, nil
}

func (n *Node) ToBinary() (result []byte, err error) {
	binNode := binaryNode{
		Pointers: n.Pointers,
		Data:     n.Data,
	}
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, binNode)
	if err != nil {
		return result, err
	}

	return buf.Bytes(), nil
}

func insertUint64at(ara []uint64, i int, val uint64) []uint64 {
	// https://github.com/golang/go/wiki/SliceTricks
	ara = append(ara, 0)
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}
