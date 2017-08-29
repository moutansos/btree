package btree

import (
	"math/rand"
	"testing"
)

func TestNewNode(t *testing.T) {
	tree := new(BTreeOnDisk)
	_, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertInt64at(t *testing.T) {
	ara := []int64{23, 45, 56, 78, 9}
	ara = insertInt64at(ara, 1, 67)
	if ara[1] != 67 {
		t.Error("Invalid value at the insertion point")
	} else if ara[0] != 23 {
		t.Error("Invalid value at position 0")
	} else if ara[2] != 45 {
		t.Error("Invalid value at position 2")
	}
}

func TestToBinary(t *testing.T) {
	tree := new(BTreeOnDisk)

	n, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}

	n.Data = [31]Index{
		Index{Key: 2, Pointer: 23},
		Index{Key: 3, Pointer: 67},
		Index{Key: 4, Pointer: 78},
		Index{Key: 6, Pointer: 89},
	}
	n.Pointers = [32]int64{1, 2, 3, 4, 5}

	_, err = n.ToBinary()
	if err != nil {
		t.Error(err)
	}
	//TODO: Check binary data from this method
}

func TestIsValidAddress(t *testing.T) {
	validAddrs := []int64{0, 752, 1504, 2256}
	invalidAddrs := []int64{-1, -4, 10, 2, 5032, 3432}

	for _, addr := range validAddrs {
		valid := IsValidAddress(addr)
		if !valid {
			t.Errorf("Valid node address of %v marked as invalid.", addr)
			return
		}
	}

	for _, addr := range invalidAddrs {
		valid := IsValidAddress(addr)
		if valid {
			t.Errorf("Invalid node address of %v marked as valid.", addr)
			return
		}
	}
}

func TestIsEmptyNode(t *testing.T) {
	tree := new(BTreeOnDisk)
	n, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}

	if !n.IsEmpty() {
		t.Error("The node is supposed to be empty and is not!")
	}

	n.Pointers[0] = 32

	if n.IsEmpty() {
		t.Error("The node is supposed to have value and IsEmpty() returned that it it does not!")
	} else {
		//Undo for data portion of the test
		n.Pointers[0] = 0
	}

	n.Data[15] = Index{Key: 2, Pointer: 3253}

	if n.IsEmpty() {
		t.Error("The node is supposed to have value and IsEmpty() returned that it does not!")
	}
}

func TestNodeIsFull(t *testing.T) {
	tree := new(BTreeOnDisk)
	n, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}

	if n.nodeIsFull() {
		t.Error("node is empty but reports as full")
	}

	//Seed data
	for i := 0; i < len(n.Data); i++ {
		n.Data[i].Key = rand.Uint64()
		for n.Data[i].Key == 0 {
			n.Data[i].Key = rand.Uint64()
		}
	}

	if !n.nodeIsFull() {
		t.Error("node is full but reports as not full")
	}
}

func TestNodeSize(t *testing.T) {
	tree := new(BTreeOnDisk)
	n, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}

	if n.size() != 0 {
		t.Errorf("node size should be zero but is actually %v", n.size())
	}

	n.Data[0] = Index{Key: 324, Pointer: 2}
	n.Data[1] = Index{Key: 325, Pointer: 2}
	n.Data[2] = Index{Key: 327, Pointer: 2}
	n.Data[3] = Index{Key: 343, Pointer: 2}

	if n.size() != 4 {
		t.Errorf("node size should be 4 but is actually %v", n.size())
	}
}

func TestFindMedianDataPoints(t *testing.T) {
	tree := new(BTreeOnDisk)
	n, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}

	n.Data[0] = Index{Key: 324, Pointer: 2}
	n.Data[1] = Index{Key: 325, Pointer: 2}
	n.Data[2] = Index{Key: 327, Pointer: 2}
	n.Data[3] = Index{Key: 343, Pointer: 2}

	left, right, err := n.findMedianDataPoints()
	if err != nil {
		t.Error(err)
	} else if left != 1 {
		t.Errorf("was expecting 1 for the left data point, got %v", left)
	} else if right != 2 {
		t.Errorf("was expecting 2 for the right data point, got %v", right)
	}

	//Check on an odd number of elements

	n.Data[4] = Index{Key: 432, Pointer: 2}
	n.Data[5] = Index{Key: 463, Pointer: 2}
	n.Data[6] = Index{Key: 784, Pointer: 2}
	left, right, err = n.findMedianDataPoints()
	if err != nil {
		t.Error(err)
	} else if left != 2 {
		t.Errorf("was expecting 2 for the left data point, got %v", left)
	} else if right != 3 {
		t.Errorf("was expecting 3 for the right data point, got %v", right)
	}
}
