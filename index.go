package btree

type Index struct {
	Key     uint64
	Pointer int64
}

func (i *Index) isEmptyOrDefault() bool {
	if i.Key == 0 && i.Pointer == 0 {
		return true
	}
	return false
}
