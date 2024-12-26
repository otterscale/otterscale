package openhdc

import (
	"github.com/openhdc/openhdc/api/workload/v1"
)

const defaultReadBatchSize = 1000

type Reader struct {
	opts readOptions
}

type readOptions struct {
	batchSize   int64
	keys        []string
	skipKeys    []string
	syncOptions []*workload.Sync_Option
}

var defaultReadOptions = readOptions{
	batchSize: defaultReadBatchSize,
}

type ReadOption interface {
	apply(*readOptions)
}

type funcReadOption struct {
	f func(*readOptions)
}

var _ ReadOption = (*funcReadOption)(nil)

func (fro *funcReadOption) apply(ro *readOptions) {
	fro.f(ro)
}

func newFuncReadOption(f func(*readOptions)) *funcReadOption {
	return &funcReadOption{
		f: f,
	}
}

func WithBatchSize(s int64) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.batchSize = s
	})
}

func WithKeys(ks ...string) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.keys = ks
	})
}

func WithSkipKeys(sks ...string) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.skipKeys = sks
	})
}

func WithOptions(opts ...*workload.Sync_Option) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.syncOptions = opts
	})
}

func NewReader(opt ...ReadOption) *Reader {
	opts := defaultReadOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	r := &Reader{
		opts: opts,
	}
	return r
}
