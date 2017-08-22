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
	n.Data[0] = Index{Key: 34, Pointer: 423}

	err = n.Write()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRemoveNode(t *testing.T) {
	f := path.Join(os.TempDir(), "test-remove-node.bin")
	//f := "test-remove-node.bin"

	//Create test data in the tree
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
	}

	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}
	n.Data[0] = Index{
		Key:     1,
		Pointer: 214,
	}

	err = n.Write()
	if err != nil {
		t.Error(err)
	}

	err = tree.RemoveNode(0)
	if err != nil {
		t.Error(err)
	}
}

func TestReadNode(t *testing.T) {
	f := path.Join(os.TempDir(), "test-read-node.bin")
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

	n.Data[0] = Index{Key: 2, Pointer: 345}
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
	} else if rn.Data[0].Key != 2 && rn.Data[0].Pointer != 345 {
		t.Errorf("Invalid data %v given by the read function at index 0. Expected Key: 2 and Pointer 345", rn.Data[0])
	} else if rn.Pointers[0] != 1 {
		t.Errorf("Invalid pointer %v given by the read function. Expected 1", rn.Pointers[0])
	} else if rn.Pointers[1] != 2 {
		t.Errorf("Invalid pointer %v given by the read function. Expected 2", rn.Pointers[1])
	}
}

func TestNextNodeAddress(t *testing.T) {
	f := path.Join(os.TempDir(), "test-next-node-address.bin")
	//f := "test-next-node-address.bin"

	//Create test data in the tree
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
	}
	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}

	err = n.Write()
	if err != nil {
		t.Error(err)
	}

	addr, err := tree.NextNodeAddress()
	if err != nil {
		t.Error(err)
	}

	if addr != 752 {
		t.Errorf("The address of %v is invalid. Expected 572", addr)
	}
}

func TestNextNodeAddressNewNodes(t *testing.T) {
	f := path.Join(os.TempDir(), "test-next-node-address-new-node.bin")
	//f := "test-next-node-address-new-node.bin"

	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
	}
	n1, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	} else if n1.Address != 0 {
		t.Errorf("Invalid address on first node. Expected 0 and got %v", n1.Address)
	}
	err = n1.Write()
	if err != nil {
		t.Error(err)
	}

	n2, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	} else if n2.Address != 752 {
		t.Errorf("Invalid address on first node. Expected 752 and got %v", n2.Address)
	}
	err = n1.Write()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateAvailableAddress(t *testing.T) {
	f := path.Join(os.TempDir(), "test-update-available-addresess.bin")
	//f := "test-update-available-addresess.bin"

	//Create test data in the tree
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
	}

	var nodes [3]*Node
	nodes[0], err = tree.NewNode()
	if err != nil {
		t.Error(err)
	}
	nodes[0].Data[0] = Index{
		Key:     23,
		Pointer: 564,
	}
	nodes[0].Pointers[0] = 234
	nodes[0].Pointers[0] = 345

	nodes[1], err = tree.NewNode()
	if err != nil {
		t.Error(err)
	}
	nodes[1].Data[0] = Index{
		Key:     67,
		Pointer: 563,
	}
	nodes[1].Pointers[0] = 23324
	nodes[1].Pointers[0] = 3543

	nodes[2], err = tree.NewNode()
	if err != nil {
		t.Error(err)
	}
	nodes[2].Data[0] = Index{
		Key:     23,
		Pointer: 564,
	}
	nodes[2].Pointers[0] = 234
	nodes[2].Pointers[0] = 345

	for _, n := range nodes {
		addr, err := tree.NextNodeAddress()
		if err != nil {
			t.Error(err)
		}

		n.Address = addr
		err = n.Write()
		if err != nil {
			t.Error(err)
		}
	}

	err = tree.RemoveNode(nodes[1].Address)
	if err != nil {
		t.Error(err)
	}

	//Finished setup - testing
	err = tree.UpdateAvailableAddresess()
	if err != nil {
		t.Error(err)
	}

	if tree.AvailableAddresses[0] != 752 {
		t.Error("the UpdateAvailableAddress function has not found the empty node.")
	}
}
