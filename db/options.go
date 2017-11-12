package db

type Compression uint8

const (
	DefaultCompression Compression = iota
	NoCompression
	SnappyCompression
	nCompression
)

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)

var (
	DefaultBlockCacher                   = LRUCacher
	DefaultBlockCacheCapacity            = 8 * MiB
	DefaultBlockRestartInterval          = 16
	DefaultBlockSize                     = 4 * KiB
	DefaultCompactionExpandLimitFactor   = 25
	DefaultCompactionGPOverlapsFactor    = 10
	DefaultCompactionL0Trigger           = 4
	DefaultCompactionSourceLimitFactor   = 1
	DefaultCompactionTableSize           = 2 * MiB
	DefaultCompactionTableSizeMultiplier = 1.0
	DefaultCompactionTotalSize           = 10 * MiB
	DefaultCompactionTotalSizeMultiplier = 10.0
	DefaultCompressionType               = SnappyCompression
	DefaultIteratorSamplingRate          = 1 * MiB
	DefaultOpenFilesCacher               = LRUCacher
	DefaultOpenFilesCacheCapacity        = 500
	DefaultWriteBuffer                   = 4 * MiB
	DefaultWriteL0PauseTrigger           = 12
	DefaultWriteL0SlowdownTrigger        = 8
)

var (
	// LRUCacher is the LRU-cache algorithm.
	LRUCacher int //= &CacherFunc{cache.NewLRU}

	// NoCacher is the value to disable caching algorithm.
	NoCacher int // = &CacherFunc{}
)

type Options struct {
	BlockRestartInterval int
	BlockSize            int
	Comparer             Comparer
	ReadOnly             bool
	NoSync               bool
	NoWriteMerge         bool
	WriteBuffer          int
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

func (o *Options) GetReadOnly() bool {
	if o == nil {
		return false
	}

	return o.ReadOnly
}

func (o *Options) GetNoWriteMerge() bool {
	if o == nil {
		return false
	}

	return o.NoWriteMerge
}

func (o *Options) GetNoSync() bool {
	if o == nil {
		return false
	}

	return o.NoSync
}

func (o *Options) GetWriteBuffer() int {
	if o == nil || o.WriteBuffer <= 0 {
		return DefaultWriteBuffer
	}

	return o.WriteBuffer
}

type WriteOptions struct {
	NoWriteMerge bool
	Sync         bool
}

func (wo *WriteOptions) GetNoWriteMerge() bool {
	if wo == nil {
		return false
	}

	return wo.NoWriteMerge
}

func (wo *WriteOptions) GetSync() bool {
	if wo == nil {
		return false
	}

	return wo.Sync
}
