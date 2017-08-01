package btree

import (
	"os"
	"path"
	"testing"
)

func TestNewBTreeOnDisk(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "test-new-btree.bin")
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
	f := path.Join(dir, "test-btree-write.bin")
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

func TestReadNode(t *testing.T) {
	f := path.Join(os.TempDir(), "test-read-node.go")
	//f := "test-read-node.bin"

	//Create test data in the tree
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
	}
	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}

	n.Data[0] = 332
	n.Pointers[0] = 1
	n.Pointers[1] = 2

	n.Address = 0
	err = n.Write()
	if err != nil {
		t.Error(err)
	}

	//Read and check input
	rn, err := tree.ReadNode(0)
	if err != nil {
		t.Error(err)
	} else if rn.Address != 0 {
		t.Errorf("Invalid address %v given by the read function. Expected 0", rn.Address)
	} else if rn.Data[0] != 332 {
		t.Errorf("Invalid data %v given by the read function at index 0. Expected 332", rn.Data[0])
	} else if rn.Pointers[0] != 1 {
		t.Errorf("Invalid pointer %v given by the read function. Expected 1", rn.Pointers[0])
	} else if rn.Pointers[1] != 2 {
		t.Errorf("Invalid pointer %v given by the read function. Expected 2", rn.Pointers[1])
	}
}
