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
}

// NewNode creates a new node using the specified b-tree structure
func NewNode(t BTree) (*Node, error) {
	n := new(Node)
	n.tree = t
	return n, nil
}

// ToBinary changes this node from a in memory native structure into
// an array of binary bytes to be written to a file or stored in a
// block of memory.
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

// IsEmpty returns true if the node is empty. For speed it only
// checks if the first data entry in the node is zero. If by some
// error any of the other parts of the data array end up with
// values and the first remains zero the node will be effectively
// empty.
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

func (n *Node) remove(key uint64) (err error) {
	for i, d := range n.Data {
		if key < d.Key {
			nn, err := n.readLeftPtr(i)
			if err != nil {
				return err
			}
			return nn.remove(key)
		} else if (n.Data[i+1].isEmptyOrDefault() || len(n.Data) == i+1) && key > d.Key {
			nn, err := n.readRightPtr(i)
			if err != nil {
				return err
			}
			return nn.remove(key)
		} else if key == d.Key {
			//TODO: Remove here
		}
	}
	return fmt.Errorf("The key was not found in the b-tree")
}

func (n *Node) insert(i *Index) (err error) {
	//TODO: Increase insert performance
	if n.nodeIsFull() {
		next, err := n.splitIntoTwoSubnodes()
		if err != nil {
			return err
		}
		return next.insert(i)
	} else if n.IsEmpty() {
		n.Data[0] = *i
		return n.Write()
	}

	for x, d := range n.Data {
		if i.Key < d.Key {
			//Insert or recurse if the data goes to the left of the current value
			if n.Pointers[x] == 0 { //Insert into this node
				n.insertThisNodeLeft(i, x)
				return n.Write()
			}
			//Recurse
			nn, err := n.readLeftPtr(x)
			if err != nil {
				return err
			}
			return nn.insert(i)
		} else if (n.Data[x+1].isEmptyOrDefault() || len(n.Data) == x+1) && i.Key > d.Key {
			//If the next value is empty or the end of the block, and the new value
			// is greater than the current key, then insert or recurse
			if n.Pointers[x+1] == 0 { //Insert into this node
				n.insertThisNodeRight(i, x)
				return n.Write()
			}
			//Recurse
			nn, err := n.readRightPtr(x)
			if err != nil {
				return err
			}
			return nn.insert(i)
		} else if i.Key == d.Key {
			return fmt.Errorf("the key of %v was already in the b-tree", i.Key)
		}
	}

	return fmt.Errorf("there was an internal logic error inserting into this node, key of inserting index: %v, first data key in node: %v", i.Key, n.Data[0].Key)
}

func (n *Node) insertThisNodeLeft(i *Index, o int) {
	n.Data = insertIndexAt(n.Data, o, *i)
	n.Pointers = insertInt64at(n.Pointers, o, 0)
}

func (n *Node) insertThisNodeRight(i *Index, o int) {
	n.Data = insertIndexAt(n.Data, o+1, *i)
	n.Pointers = insertInt64at(n.Pointers, o+1, 0)
}

// Only run on nodes that are full
func (n *Node) splitIntoTwoSubnodes() (new *Node, err error) {
	median, err := n.findMedianDataPoint()
	if err != nil {
		return nil, err
	}

	//Copy the values to the left node
	leftNode, err := n.tree.NewNode()
	if err != nil {
		return nil, err
	}

	for i, e := range n.Data[:median] {
		leftNode.Data[i] = e
		leftNode.Pointers[i] = n.Pointers[i]
	}
	leftNode.Pointers[median] = n.Pointers[median]
	err = leftNode.Write()
	if err != nil {
		return nil, err
	}

	//Copy the values to the right node
	rightNode, err := n.tree.NewNode()
	if err != nil {
		return nil, err
	}

	rightNode.Pointers[0] = n.Pointers[median+1]
	for i, e := range n.Data[median+1 : n.size()] {
		rightNode.Data[i] = e
		rightNode.Pointers[i+1] = n.Pointers[median+2+i]
	}
	err = rightNode.Write()
	if err != nil {
		return nil, err
	}

	//Clear the old node and assign the new subnodes
	var medianVal = n.Data[median]
	n.clear()
	n.Data[0] = medianVal

	n.Pointers[0] = leftNode.Address
	n.Pointers[1] = rightNode.Address

	err = n.Write()
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) findMedianDataPoint() (medianIndex int, err error) {
	size := n.size()
	if size < 3 {
		return -1, fmt.Errorf("the b-tree only had %v elements, no median data point exists", size)
	}

	medianIndex = size / 2

	return medianIndex, nil
}

func (n *Node) size() int {
	for i, el := range n.Data {
		if el.Key == 0 {
			return i
		}
	}

	return len(n.Data)
}

func (n *Node) nodeIsFull() bool {
	//Just check the last data point. If it is not zero then it is full
	if n.Data[len(n.Data)-1].Key != 0 {
		return true
	}
	return false
}

func (n *Node) clear() {
	for i := 0; i < len(n.Data); i++ {
		n.Data[i].Key = 0
		n.Data[i].Pointer = 0
	}

	for i := 0; i < len(n.Pointers); i++ {
		n.Pointers[i] = 0
	}
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

// IsValidAddress indicates if the given value is a valid node address
// nodes are ussualy 752 bytes in lenght and therefore addresses occur
// every 752 bytes.
func IsValidAddress(addr int64) bool {
	if addr >= 0 && addr%752 == 0 {
		return true
	}
	return false
}

func insertInt64at(ara [32]int64, i int, val int64) [32]int64 {
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}

func insertIndexAt(ara [31]Index, i int, val Index) [31]Index {
	copy(ara[i+1:], ara[i:])
	ara[i] = val
	return ara
}

func removeInt64at(ara [32]int64, i int) [32]int64 {
	s := append(ara[:i], ara[i+1:]...)
	newAra := new([32]int64)
	for o, e := range s {
		newAra[o] = e
	}
	return *newAra
}

func removeIndexAt(ara [31]Index, i int) [31]Index {
	s := append(ara[:i], ara[i+1:]...)
	newAra := new([31]Index)
	for o, e := range s {
		newAra[o] = e
	}
	return *newAra
}
