package btree

import (
	"os"
)

type BTreeOnDisk struct {
	File string
}

func NewBTreeOnDisk(file string) (t *BTreeOnDisk, err error) {
	t = new(BTreeOnDisk)
	t.File = file
	return t, err
}

func (t *BTreeOnDisk) WriteNode(n *Node) error {
	data, err := n.ToBinary()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(t.File, os.O_RDWR, 0666) //TODO: Open for writing
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

func (t *BTreeOnDisk) NewNode() (n *Node, err error) {
	return NewNode(t)
}
