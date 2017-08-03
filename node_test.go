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

	n.Data = [31]int64{23, 67, 78, 89}
	n.Pointers = [32]int64{1, 2, 3, 4, 5}

	_, err = n.ToBinary()
	if err != nil {
		t.Error(err)
	}
	//TODO: Check binary data from this method
}

func TestIsValidAddress(t *testing.T) {
	validAddrs := []int64{0, 504, 1008, 1512}
	invalidAddrs := []int64{-1, -4, 10, 2, 5032, 3432}

	for _, addr := range validAddrs {
		valid := IsValidAddress(addr)
		if !valid {
			t.Errorf("Valid node address of %v marked as invalid.")
			return
		}
	}

	for _, addr := range invalidAddrs {
		valid := IsValidAddress(addr)
		if valid {
			t.Errorf("Invalid node address of %v marked as valid.")
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

	n.Data[15] = 3253

	if n.IsEmpty() {
		t.Error("The node is supposed to have value and IsEmpty() returned that it does not!")
	}
}
