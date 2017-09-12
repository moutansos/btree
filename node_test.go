package btree

import (
	"math/rand"
	"os"
	"path"
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
	ara := [32]int64{23, 45, 56, 78, 9}
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

	median, err := n.findMedianDataPoint()
	if err != nil {
		t.Error(err)
	} else if median != 2 {
		t.Errorf("was expecting 2 for the median data point, got %v", median)
	}

	//Check on an odd number of elements

	n.Data[4] = Index{Key: 432, Pointer: 2}
	n.Data[5] = Index{Key: 463, Pointer: 2}
	n.Data[6] = Index{Key: 784, Pointer: 2}
	median, err = n.findMedianDataPoint()
	if err != nil {
		t.Error(err)
	} else if median != 3 {
		t.Errorf("was expecting 3 for the median point, got %v", median)
	}
}

func TestSplitNode(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "test-node-split.bin")
	//f = "test-node-split.bin"
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}

	n.Pointers[0] = 345
	n.Data[0] = Index{Key: 324, Pointer: 2}
	n.Pointers[1] = 7438
	n.Data[1] = Index{Key: 325, Pointer: 3}
	n.Pointers[2] = 3243
	n.Data[2] = Index{Key: 327, Pointer: 4}
	n.Pointers[3] = 4737
	n.Data[3] = Index{Key: 343, Pointer: 5}
	n.Pointers[4] = 435
	n.Data[4] = Index{Key: 352, Pointer: 6}
	n.Pointers[5] = 3490

	err = n.Write()
	if err != nil {
		t.Error(err)
	}

	n, err = n.splitIntoTwoSubnodes()
	if err != nil {
		t.Error(err)
	}

	if n.Data[0].Key != 327 {
		t.Errorf("The computed top key is not right. Expected 327, got %v", n.Data[0].Key)
	}

	leftNode, err := n.readLeftPtr(0)
	if err != nil {
		t.Error(err)
	} else if leftNode.Data[0].Key != 324 ||
		leftNode.Pointers[0] != 345 {
		t.Errorf("the left key has invalid data at index 0, expected data key to be 324 and left pointer to be 345, was actually %v and %v", leftNode.Data[0].Key, leftNode.Pointers[0])
	} else if leftNode.Data[1].Key != 325 ||
		leftNode.Pointers[1] != 7438 {
		t.Errorf("the left key has invalid data at index 1, expected data key to be 325 and left pointer to be 7438, was actually %v and %v", leftNode.Data[1].Key, leftNode.Pointers[1])
	} else if leftNode.Pointers[2] != 3243 {
		t.Errorf("the left key has an invalid right pointer at index 1, expected 3243, got %v", leftNode.Pointers[2])
	}

	rightNode, err := n.readRightPtr(0)
	if err != nil {
		t.Error(err)
	} else if rightNode.Data[0].Key != 343 ||
		rightNode.Pointers[0] != 4737 {
		t.Errorf("the left key has invalid data at index 0, expected data key to be 343 and left pointer to be 4737, was actually %v and %v", rightNode.Data[0].Key, rightNode.Pointers[0])
	} else if rightNode.Data[1].Key != 352 ||
		rightNode.Pointers[1] != 435 {
		t.Errorf("the left key has invalid data at index 1, expected data key to be 352 and left pointer to be 435, was actually %v and %v", rightNode.Data[1].Key, rightNode.Pointers[1])
	} else if rightNode.Pointers[2] != 3490 {
		t.Errorf("the left key has an invalid right pointer at index 1, expected 3490, got %v", rightNode.Pointers[2])
	}
}

func TestQuery(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "test-node-query.bin")
	//f = "test-node-query.bin"
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}

	n.Pointers[0] = 752
	n.Data[0] = Index{Key: 23, Pointer: 98}
	n.Pointers[1] = 32423

	err = n.Write()
	if err != nil {
		t.Error(err)
	}

	n2, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}
	n2.Data[0] = Index{Key: 10, Pointer: 78}
	n2.Data[1] = Index{Key: 12, Pointer: 93}

	err = n2.Write()
	if err != nil {
		t.Error(err)
	}

	i, err := n.query(12)
	if err != nil {
		t.Error(err)
	} else if i.Pointer != 93 {
		t.Errorf("the returned index should have a pointer of 93 but had %v", i.Pointer)
	}
}
