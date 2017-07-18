package btree

import "testing"

func TestNewBTreeInMem(t *testing.T) {
	tree, err := NewBTreeInMem(100)
	if err != nil {
		t.Error(err)
	} else if tree.data[0] != byte(100) {
		t.Errorf("Wrong first byte: %v", tree.data[0])
	}
}
