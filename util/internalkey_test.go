package util

import (
	"bytes"
	"testing"
)

func TestMakeKey(t *testing.T) {
	userKey := []byte{1, 2, 3, 4}
	k := makeInternalKey(userKey, KeyTypeSet, 12)

	if 0 != bytes.Compare(k.userKey(), userKey) {
		t.Errorf("userKey mismatch, expected %v actual %v", userKey, k.userKey())
	}
}

func TestCompareKey(t *testing.T) {
	key1 := []byte{2, 3, 4}
	key2 := []byte{1}

	k1 := makeInternalKey(key1, KeyTypeSet, 12)
	k2 := makeInternalKey(key2, KeyTypeSet, 1)

	cmp := internalKeyComparer{}
	comRes := cmp.Compare([]byte(k1), []byte(k2))
	if comRes <= 0 {
		t.Errorf("unexpected compare result - %v", comRes)
	}
}

func TestCompareType(t *testing.T) {
	key1 := []byte{2, 3, 4}

	k1 := makeInternalKey(key1, KeyTypeSet, 12)
	k2 := makeInternalKey(key1, KeyTypeDelete, 1)

	cmp := internalKeyComparer{}
	comRes := cmp.Compare([]byte(k1), []byte(k2))
	if comRes >= 0 {
		t.Errorf("unexpected compare result - %v", comRes)
	}
}

func TestCompareSeq(t *testing.T) {
	key1 := []byte{1}

	k1 := makeInternalKey(key1, KeyTypeSet, 12)
	k2 := makeInternalKey(key1, KeyTypeSet, 1)

	cmp := internalKeyComparer{}
	comRes := cmp.Compare([]byte(k1), []byte(k2))
	if comRes >= 0 {
		t.Errorf("unexpected compare result - %v", comRes)
	}
}
