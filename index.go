package btree

type Index struct {
	Key     uint64
	Pointer int64
}

// NewIndex creates a new key/pointer index pair for the b-tree structure
func NewIndex(key uint64, pointer int64) *Index {
	i := Index{
		key,
		pointer,
	}
	return &i
}

func (i *Index) isEmptyOrDefault() bool {
	if i.Key == 0 && i.Pointer == 0 {
		return true
	}
	return false
}
