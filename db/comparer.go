package db

import "bytes"

type Comparer interface {
	Compare(a, b []byte) int
	Name() string
	Separator(dst, a, b []byte) []byte
}

var DefaultComparer Comparer = defaultCmp{}

type defaultCmp struct{}

func (defaultCmp) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func (defaultCmp) Name() string {
	return "BytesComparator"
}

func (defaultCmp) Separator(dst, a, b []byte) []byte {
	i, n := SharedPrefixLen(a, b), len(dst)
	dst = append(dst, a...)
	if len(b) > 0 {
		if i == len(a) {
			return dst
		}
		if i == len(b) {
			panic("a < b is a precondition, but b is a prefix of a")
		}
		if a[i] == 0xff || a[i]+1 >= b[i] {
			return dst
		}
	}
	i += n
	for ; i < len(dst); i++ {
		if dst[i] != 0xff {
			dst[i]++
			return dst[:i+1]
		}
	}
	return dst
}

func SharedPrefixLen(a, b []byte) int {
	i, n := 0, len(a)
	if n > len(b) {
		n = len(b)
	}

	for i < n && a[i] == b[i] {
		i++
	}

	return i
}
