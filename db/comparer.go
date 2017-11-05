package db

import "bytes"

type Comparer interface {
	Compare(a, b []byte) int
	Name() string
}

var DefaultComparer Comparer = defaultCmp{}

type defaultCmp struct{}

func (defaultCmp) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func (defaultCmp) Name() string {
	return "BytesComparator"
}
