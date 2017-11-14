package lsmdb

type version struct {
	s   *session
	ref int
}

func newVersion(s *session) *version {
	return &version{s: s}
}

func (v *version) incref() {
	v.ref++
}

func (v *version) release() {
	v.ref--
}

func (v *version) pickMemdbLevel(umin, umax []byte, maxLevel int) (level int) {
	return 0
}
