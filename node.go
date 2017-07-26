package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const maxInt64 = 18446744073709551615

type Node struct {
	Pointers [32]uint64
	Data     [31]uint64
	
	Address  uint64
	tree     BTree
}

type binaryNode struct {
	Pointers [32]uint64
	Data     [31]uint64
}

func NewNode(t BTree) (*Node, error) {
	n := new(Node)
	n.tree = t
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

func (n *Node) Write() error {
	if n.tree != nil {
		return n.tree.WriteNode(n)
	} else {
		return fmt.Errorf("There was no tree attached to this node")
	}
}

func insertUint64at(ara []uint64, i int, val uint64) []uint64 {
	// https://github.com/golang/go/wiki/SliceTricks
	ara = append(ara, 0)
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}

func ToInt64Ara(i uint64) (ara []int64) {
	n := i / maxInt64
	r := i % maxInt64
	
	if r == 0 {
		ara = make([]int64, n)
	} else {
		ara = make([]int64, n + 1)
		ara[n] = r
	}

	for x := 0; x < n; x++ {
		ara[x] = maxInt64
	}

	return ara
}
