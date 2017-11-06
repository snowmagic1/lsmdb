package table

import "encoding/binary"

type blockHandle struct {
	offset, length uint64
}

type block []byte

type indexEntry struct {
	bh     blockHandle
	keyLen int
}

const (
	blockTrailerLen = 5
	footerLen       = 48

	magic = "\x55\x93\x64\x9a\xfe\x32\x12\xac"
)

func encodeBlockHandle(dst []byte, b blockHandle) int {
	n := binary.PutUvarint(dst, b.offset)
	m := binary.PutUvarint(dst[n:], b.length)
	return n + m
}

func decodeBlockHandle(src []byte) (blockHandle, int) {
	offset, n := binary.Uvarint(src)
	length, m := binary.Uvarint(src[n:])
	if n == 0 || m == 0 {
		return blockHandle{}, 0
	}
	return blockHandle{offset, length}, n + m
}
