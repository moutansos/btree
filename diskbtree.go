package btree

import (
	"bytes"
	"encoding/binary"
	"os"
)

type BTreeOnDisk struct {
	File string
}

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

func (t *BTreeOnDisk) WriteNode(n *Node) error {
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

func (t *BTreeOnDisk) ReadNode(address int64) (n *Node, err error) {
	f, err := os.Open(t.File)
	if os.IsNotExist(err) {
		return nil, err
	}
	defer f.Close()

	if err != nil {
		return nil, err
	}

	data := make([]byte, 504)

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

func (t *BTreeOnDisk) NewNode() (n *Node, err error) {
	return NewNode(t)
}
