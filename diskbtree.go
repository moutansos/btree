package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// BTreeOnDisk is a structure that references a b-tree structure that
// resides on disk instead of in memory.
type BTreeOnDisk struct {
	File string
}

// NewBTreeOnDisk creates a new b-tree that resides on disk. The
// structure uses internal pointers to bytes in the file. It uses
// these to work like memory pointers.
func NewBTreeOnDisk(file string) (t *BTreeOnDisk, err error) {
	t = new(BTreeOnDisk)
	t.File = file

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return t, nil
	}

	err = os.Remove(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// WriteNode writes the specified node to disk. It takes a single
// parameter node. It uses the address inside the n *Node parameter
// and confirms that it is a valid pointer.
func (t *BTreeOnDisk) WriteNode(n *Node) error {
	if !IsValidAddress(n.Address) {
		return fmt.Errorf("Invalid address. Cannot write node at %v", n.Address)
	}

	data, err := n.ToBinary()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(t.File, os.O_RDWR, 0666)
	if os.IsNotExist(err) {
		f, err = os.Create(t.File)
	}
	defer f.Close()

	if err != nil {
		return err
	}

	_, err = f.WriteAt(data, n.Address)
	return err
}

// ReadNode reads the node from disk. The parameter takes a positive
// integer and returns an error if the address is invalid. The function
// returns two parameters n *Node which is the node and err of type
// error.
func (t *BTreeOnDisk) ReadNode(address int64) (n *Node, err error) {
	if !IsValidAddress(address) {
		return nil, fmt.Errorf("Invalid address. Cannot read node at %v", address)
	}

	f, err := os.Open(t.File)
	if os.IsNotExist(err) {
		return nil, err
	}
	defer f.Close()

	if err != nil {
		return nil, err
	}

	data := make([]byte, 752)

	//TODO: Validate address
	_, err = f.Seek(address, 0)
	if err != nil {
		return nil, err
	}

	_, err = f.Read(data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	bn := new(binaryNode)

	err = binary.Read(buf, binary.LittleEndian, bn)
	if err != nil {
		return nil, err
	}

	n = new(Node)
	n.Pointers = bn.Pointers
	n.Data = bn.Data
	n.Address = address
	return n, nil
}

func (t *BTreeOnDisk) RemoveNode(addr int64) (err error) {
	if !IsValidAddress(addr) {
		return fmt.Errorf("the provided address of %v is invalid")
	}

	f, err := os.Open(t.File)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	treeSize := info.Size()
	if addr > treeSize {
		return fmt.Errorf("The provided address is larger than the tree")
	}

	return nil
}

// NewNode calls the standalone NewNode function and gives it the
// calling binary tree.
func (t *BTreeOnDisk) NewNode() (n *Node, err error) {
	addr, err := t.NextNodeAddress()
	if err != nil {
		return nil, err
	}
	n, err = NewNode(t)
	n.Address = addr
	return n, err
}

func (t *BTreeOnDisk) NextNodeAddress() (int64, error) {
	stat, err := os.Stat(t.File)
	if os.IsNotExist(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}

	addr := stat.Size()
	if IsValidAddress(addr) {
		return addr, nil
	}
	return -1, fmt.Errorf("the address %v was invalid and indicates a corrupt b-tree structure", addr)
}
