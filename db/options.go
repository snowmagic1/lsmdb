package db

type Compression uint8

const (
	CompressionNo     Compression = 0
	CompressionSnappy Compression = 1
)

type Options struct {
	BlockRestartInterval int
	BlockSize            int
	Comparer             Comparer
}

func (o *Options) GetBlockRestartInterval() int {
	if o == nil || o.BlockRestartInterval <= 0 {
		return 16
	}
	return o.BlockRestartInterval
}

func (o *Options) GetBlockSize() int {
	if o == nil || o.BlockSize <= 0 {
		return 4096
	}

	return o.BlockSize
}

func (o *Options) GetComparer() Comparer {
	if o == nil || o.Comparer == nil {
		return DefaultComparer
	}

	return o.Comparer
}
