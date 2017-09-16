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

func TestInsertIndexAt(t *testing.T) {
	ara := [31]Index{
		Index{Key: 32, Pointer: 43},
		Index{Key: 53, Pointer: 423},
		Index{Key: 79, Pointer: 324},
		Index{Key: 83, Pointer: 432},
		Index{Key: 93, Pointer: 493},
	}
	ara = insertIndexAt(ara, 2, Index{Key: 5, Pointer: 32})
	if ara[1].Key != 53 {
		t.Error("invalid value right before insertion point")
	} else if ara[2].Key != 5 {
		t.Error("invalid value at the insertion point")
	} else if ara[3].Key != 79 {
		t.Error("invalid value right after the insertion point")
	}
}

func TestNodeSize(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "test-node-size.bin")
	//f = "test-node-size.bin"
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
		return
	}

	s := n.size()
	if s != 0 {
		t.Errorf("node size was supposed to be zero but was actually %v", s)
	}

	testKeys := []uint64{2, 4, 5, 8, 10, 67, 89}
	for _, key := range testKeys {
		err = n.insert(NewIndex(key, 1))
		if err != nil {
			t.Error(err)
			return
		}
	}

	s = n.size()
	if s != 7 {
		t.Errorf("node size was supposed to be seven but was actually %v", s)
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

func TestInsertIndex(t *testing.T) {
	dir := os.TempDir()
	f := path.Join(dir, "test-node-insert-index.bin")
	//f = "test-node-insert-index.bin"
	tree, err := NewBTreeOnDisk(f)
	if err != nil {
		t.Error(err)
		return
	}

	n, err := tree.NewNode()
	if err != nil {
		t.Error(err)
	}

	i1 := Index{Key: 30, Pointer: 78}
	err = n.insert(&i1)
	if err != nil {
		t.Error(err)
	} else if n.Data[0].Key != i1.Key && n.Data[0].Pointer != i1.Pointer {
		t.Errorf("invalid insert of the first index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i1.Key, i1.Pointer, n.Data[0].Key, n.Data[0].Pointer)
	}

	i2 := Index{Key: 45, Pointer: 89}
	err = n.insert(&i2)
	if err != nil {
		t.Error(err)
	} else if n.Data[0].Key != i1.Key && n.Data[0].Pointer != i1.Pointer {
		t.Errorf("invalid insert of the first index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i1.Key, i1.Pointer, n.Data[0].Key, n.Data[0].Pointer)
	} else if n.Data[1].Key != i2.Key && n.Data[1].Pointer != i2.Pointer {
		t.Errorf("invalid insert of the second index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i2.Key, i2.Pointer, n.Data[1].Key, n.Data[1].Pointer)
	}

	i3 := Index{Key: 5, Pointer: 67}
	err = n.insert(&i3)
	if err != nil {
		t.Error(err)
	} else if n.Data[0].Key != i3.Key && n.Data[0].Pointer != i3.Pointer {
		t.Errorf("invalid insert of the first index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i3.Key, i3.Pointer, n.Data[0].Key, n.Data[0].Pointer)
	} else if n.Data[1].Key != i1.Key && n.Data[1].Pointer != i1.Pointer {
		t.Errorf("invalid insert of the second index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i1.Key, i1.Pointer, n.Data[1].Key, n.Data[1].Pointer)
	} else if n.Data[2].Key != i2.Key && n.Data[2].Pointer != i2.Pointer {
		t.Errorf("invalid insert of the third index. Expected key of %v and pointer of %v, got key of %v and pointer of %v", i2.Key, i2.Pointer, n.Data[2].Key, n.Data[2].Pointer)
	}
}
