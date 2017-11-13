package util

type Releaser interface {
	Release()
}

type BasicReleaser struct {
	releaser Releaser
	released bool
}

func (r *BasicReleaser) Released() bool {
	return r.released
}

func (r *BasicReleaser) Release() {
}

func (r *BasicReleaser) SetReleaser(releaser Releaser) {

}
