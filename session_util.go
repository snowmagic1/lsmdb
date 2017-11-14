package lsmdb

import "sync/atomic"

func (s *session) allocFileNum() int64 {
	return atomic.AddInt64(&s.stNextFileNum, 1) // - 1
}

func (s *session) version() *version {
	s.vMu.Lock()
	defer s.vMu.Unlock()

	s.stVersion.incref()

	return s.stVersion
}
