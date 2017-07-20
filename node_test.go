package btree

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	tree := new(BTree)
	_, err := NewNode(tree)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertUint64at(t *testing.T) {
	ara := []uint64{23, 45, 56, 78, 9}
	ara = insertUint64at(ara, 1, 67)
	if ara[1] != 67 {
		t.Error("Invalid value at the insertion point")
	} else if ara[0] != 23 {
		t.Error("Invalid value at position 0")
	} else if ara[2] != 45 {
		t.Error("Invalid value at position 2")
	}
}
