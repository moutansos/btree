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

