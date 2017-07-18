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
	tree, err := NewBTreeOnDisk(f, 100)
	if tree.File != f {
		t.Errorf("The file %v is invalid", tree.File)
	} else if tree.BlockSize != 100 {
		t.Errorf("The block size %v is invalid", tree.BlockSize)
	} else if err != nil {
		t.Error(err)
	}
}
