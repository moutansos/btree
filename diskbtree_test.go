package btree

import (
	"os"
	"path"
	"testing"
)

func TestNewBTreeOnDisk(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "btree.bin")
	//f := "btree.bin"
	tree, err := NewBTreeOnDisk(f)
	if tree.File != f {
		t.Errorf("The file %v is invalid", tree.File)
	} else if err != nil {
		t.Error(err)
	}
}

func TestWriteNode(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "btree-write.bin")
	//f = "btree-write.bin"
	dtree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := dtree.NewNode()
	if err != nil {
		t.Error(err)
		return
	}
	n.Pointers[0] = 1
	n.Pointers[1] = 2
	n.Data[0] = 423

	err = n.Write()
	if err != nil {
		t.Error(err)
		return
	}
}
