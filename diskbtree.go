package btree

import (
	"bytes"
	"encoding/binary"
	"os"
	"fmt"
)

type BTreeOnDisk struct {
	File      string
	BlockSize uint64
}

func NewBTreeOnDisk(file string, bsize uint64) (*BTreeOnDisk, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, bsize)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(buf.Bytes())

	t := new(BTreeOnDisk)
	t.File = file
	t.BlockSize = bsize
	return t, err
}

func (t *BTreeOnDisk) GetBlockSize() uint64 {
	return t.BlockSize
}

func (t *BTreeOnDisk) WriteNode(n *Node) error {
	return fmt.Errorf("Unimplemented")
}

func (t *BTreeOnDisk) NewNode() {
}
