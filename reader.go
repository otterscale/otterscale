package openhdc

import (
	"bytes"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

const defaultReadBatchSize = 1000

type Reader struct {
	opts readOptions
}

type readOptions struct {
	sourceName  string
	syncedAt    time.Time
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

func WithSourceName(n string) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.sourceName = n
	})
}

func WithSyncedAt(t time.Time) ReadOption {
	return newFuncReadOption(func(o *readOptions) {
		o.syncedAt = t
	})
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

func (r *Reader) BatchSize() int64 {
	return r.opts.batchSize
}

func (r *Reader) Keys() []string {
	return r.opts.keys
}

func (r *Reader) SkipKeys() []string {
	return r.opts.skipKeys
}

func (r *Reader) SyncMode(v string) property.SyncMode {
	if v == "" {
		return 0
	}
	for _, o := range r.opts.syncOptions {
		if o.GetKey() == v {
			return o.GetMode()
		}
	}
	return 0
}

func (r *Reader) SyncCursors(v string) []*workload.Sync_Option_Cursor {
	if v == "" {
		return nil
	}
	for _, o := range r.opts.syncOptions {
		if o.GetKey() == v {
			return o.GetCursors()
		}
	}
	return nil
}

func (r *Reader) Send(kind property.MessageKind, msgs chan<- *pb.Message, rec arrow.Record) error {
	msg, err := r.toMessage(kind, rec)
	if err != nil {
		return err
	}
	msgs <- msg
	return nil
}

func (r *Reader) toMessage(kind property.MessageKind, rec arrow.Record) (*pb.Message, error) {
	var buf bytes.Buffer
	w := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()), ipc.WithAllocator(memory.DefaultAllocator))
	defer w.Close()
	if err := w.Write(rec); err != nil {
		return nil, err
	}
	return pb.Message_builder{
		Kind:       &kind,
		Record:     buf.Bytes(),
		SourceName: &r.opts.sourceName,
		SyncedAt:   timestamppb.New(r.opts.syncedAt),
	}.Build(), nil
}
