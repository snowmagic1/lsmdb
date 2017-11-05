package util

import (
	"bytes"
)

type internalKey []byte
type internalKeyType uint8

const (
	KeyTypeDelete internalKeyType = 0
	KeyTypeSet    internalKeyType = 1
	KeyTypeMax    internalKeyType = 1

	// type | seqNum
	KeyHeaderLength int = 1 + 8
)

func makeInternalKey(userKey []byte, userKeyType internalKeyType, seqNum uint64) internalKey {
	iKey := make(internalKey, len(userKey)+KeyHeaderLength)
	i := copy(iKey, userKey)
	iKey[i+0] = uint8(userKeyType)
	iKey[i+1] = uint8(seqNum)
	iKey[i+2] = uint8(seqNum >> 1 * 8)
	iKey[i+3] = uint8(seqNum >> 2 * 8)
	iKey[i+4] = uint8(seqNum >> 3 * 8)
	iKey[i+5] = uint8(seqNum >> 4 * 8)
	iKey[i+6] = uint8(seqNum >> 5 * 8)
	iKey[i+7] = uint8(seqNum >> 6 * 8)

	return iKey
}

func (k internalKey) valid() bool {
	i := len(k) - KeyHeaderLength
	return i >= 0 && internalKeyType(k[i]) <= KeyTypeMax
}

func (k internalKey) userKey() []byte {
	return k[:len(k)-KeyHeaderLength]
}

func (k internalKey) keyType() internalKeyType {
	i := len(k) - KeyHeaderLength
	return internalKeyType(k[i])
}

func (k internalKey) seqNum() uint64 {
	i := len(k) - 7
	seq := uint64(k[i+0])
	seq |= uint64(k[i+1]) << 1 * 8
	seq |= uint64(k[i+2]) << 2 * 8
	seq |= uint64(k[i+3]) << 3 * 8
	seq |= uint64(k[i+4]) << 4 * 8
	seq |= uint64(k[i+5]) << 5 * 8
	seq |= uint64(k[i+6]) << 6 * 8

	return seq
}

type internalKeyComparer struct{}

func (c internalKeyComparer) Compare(a, b []byte) int {
	ak, bk := internalKey(a), internalKey(b)

	if cmp := bytes.Compare(ak.userKey(), bk.userKey()); cmp != 0 {
		return cmp
	}

	if as, bs := ak.seqNum(), bk.seqNum(); as < bs {
		return 1
	} else if as > bs {
		return -1
	}

	if at, bt := ak.keyType(), bk.keyType(); at < bt {
		return 1
	} else if at < bt {
		return -1
	}

	return 0
}
