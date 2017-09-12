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
	File               string
	AvailableAddresses []int64
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
	n.tree = t
	return n, nil
}

// RemoveNode removes a node from a b-tree structure by writing all of
// the bytes in the section to zero. Then it adds the address of the
// removed node to the cache of available addresses.
func (t *BTreeOnDisk) RemoveNode(addr int64) (err error) {
	if !IsValidAddress(addr) {
		return fmt.Errorf("the provided address of %v is invalid", addr)
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

	blankNode, err := t.NewNode()
	if err != nil {
		return err
	}
	blankNode.Address = addr
	t.AvailableAddresses = append(t.AvailableAddresses, addr)
	return blankNode.Write()
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

// AddressIsAvailable checks the input address in the node and returns true
// if the node at the specific address is empty and therefore available.
func (t *BTreeOnDisk) AddressIsAvailable(addr int64) (available bool, err error) {
	for _, e := range t.AvailableAddresses {
		if e == addr {
			return true, nil
		}
	}
	return false, nil
}

func (t *BTreeOnDisk) NextNodeAddress() (int64, error) {
	if len(t.AvailableAddresses) > 0 {
		val := t.AvailableAddresses[0]
		t.AvailableAddresses = t.AvailableAddresses[1:]
		return val, nil
	}

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

func (t *BTreeOnDisk) UpdateAvailableAddresess() (err error) {
	stat, err := os.Stat(t.File)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	size := stat.Size()

	var i int64
	for i = 0; i < size; i = i + 752 { //Iterate through every node
		n, err := t.ReadNode(i)
		if err != nil {
			return err
		}

		isAvailable, err := t.AddressIsAvailable(i)
		if err != nil {
			return fmt.Errorf("unable to check for available address")
		}

		if n.IsEmpty() && !isAvailable {
			t.AvailableAddresses = append(t.AvailableAddresses, i)
		}
	}
	return nil
}

func (t *BTreeOnDisk) QueryIndex(key uint64) (index *Index, err error) {
	n, err := t.ReadNode(0)
	if !n.IsEmpty() {
		return n.query(key)
	}
	return nil, fmt.Errorf("the b-tree is empty")
}
