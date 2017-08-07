package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const maxInt64 = 18446744073709551615

// Node is a structure that represents a node when in memory ouside the tree.
// It is used for creating and editing nodes and is then written from there.
type Node struct {
	Pointers [32]int64
	Data     [31]Index

	Address int64
	tree    BTree
}

type binaryNode struct { //752 bytes
	Pointers [32]int64
	Data     [31]Index
	//TODO: Change data type to index type that contains the index and the pointer to the data
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

func (n *Node) IsEmpty() bool {
	for _, p := range n.Pointers {
		if p != 0 {
			return false
		}
	}

	for _, d := range n.Data {
		if !d.isEmptyOrDefault() {
			return false
		}
	}

	return true
}

func IsValidAddress(addr int64) bool {
	if addr >= 0 && addr%752 == 0 {
		return true
	}
	return false
}

func insertInt64at(ara []int64, i int, val int64) []int64 {
	// https://github.com/golang/go/wiki/SliceTricks
	ara = append(ara, 0)
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}
