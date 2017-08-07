package btree

import (
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
