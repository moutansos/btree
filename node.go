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

func (n *Node) query(key uint64) (index *Index, err error) {
	for i, d := range n.Data {
		if key < d.Key {
			nn, err := n.readLeftPtr(i)
			if err != nil {
				return nil, err
			}
			return nn.query(key)
		} else if (n.Data[i+1].isEmptyOrDefault() || len(n.Data) == i+1) && key > d.Key {
			nn, err := n.readRightPtr(i)
			if err != nil {
				return nil, err
			}
			return nn.query(key)
		} else if key == d.Key {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("The key was not found in the b-tree")
}

func (n *Node) insert(i *Index) (err error) {
	//TODO: Look into equations for managing b-tree height and inserting nodes into the
	//		tree in such a way that it doesn't turned into a list of linked arrays
	return nil
}

func (n *Node) nodeIsFull() bool {
	//Just check the last data point. If it is not zero then it is full
	if n.Data[len(n.Data)-1].Key != 0 {
		return true
	}
	return false
}

func (n *Node) readLeftPtr(index int) (newNode *Node, err error) {
	ptr := n.Pointers[index]
	if ptr != 0 {
		newNode, err = n.tree.ReadNode(ptr)
		if err != nil {
			return nil, err
		}
		return newNode, err
	}
	return nil, fmt.Errorf("the key was not found, the pointer was not referenced")
}

func (n *Node) readRightPtr(index int) (newNode *Node, err error) {
	return n.readLeftPtr(index + 1)
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
