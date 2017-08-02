package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const maxInt64 = 18446744073709551615

type Node struct {
	Pointers [32]int64
	Data     [31]int64

	Address int64
	tree    BTree
}

type binaryNode struct { //504 bytes
	Pointers [32]int64
	Data     [31]int64
	//TODO: Add error detection and recovery
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
	}
	return fmt.Errorf("There was no tree attached to this node")
}

func IsValidAddress(addr int64) bool {
	if addr >= 0 && addr%504 == 0 {
		return true
	}
	return false
}

func insertUint64at(ara []uint64, i int, val uint64) []uint64 {
	// https://github.com/golang/go/wiki/SliceTricks
	ara = append(ara, 0)
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}
