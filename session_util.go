package lsmdb

import "sync/atomic"

func (s *session) allocFileNum() int64 {
	return atomic.AddInt64(&s.stNextFileNum, 1) // - 1
}
