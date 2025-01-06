package openhdc

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/openhdc/openhdc/api/property/v1"
)

const (
	defaultStopTimeout = time.Second * 2
	defaultLogLevel    = slog.LevelDebug
)

type App struct {
	opts appOptions

	ctx    context.Context
	cancel context.CancelFunc
}

type appOptions struct {
	id          string
	ctx         context.Context
	sigs        []os.Signal
	stopTimeout time.Duration

	kind    property.WorkloadKind
	name    string
	version string

	servers []*Server

	logLevel slog.Leveler
}

var defaultAppOptions = appOptions{
	id:          uuid.NewString(),
	ctx:         context.Background(),
	sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL},
	stopTimeout: defaultStopTimeout,
	logLevel:    defaultLogLevel,
}

type AppOption interface {
	apply(*appOptions)
}

type funcAppOption struct {
	f func(*appOptions)
}

var _ AppOption = (*funcAppOption)(nil)

func (fro *funcAppOption) apply(ro *appOptions) {
	fro.f(ro)
}

func newFuncAppOption(f func(*appOptions)) *funcAppOption {
	return &funcAppOption{
		f: f,
	}
}

func WithKind(k property.WorkloadKind) AppOption {
	return newFuncAppOption(func(o *appOptions) {
		o.kind = k
	})
}

func WithName(n string) AppOption {
	return newFuncAppOption(func(o *appOptions) {
		o.name = n
	})
}

func WithVersion(v string) AppOption {
	return newFuncAppOption(func(o *appOptions) {
		o.version = v
	})
}

func WithServers(srvs ...*Server) AppOption {
	return newFuncAppOption(func(o *appOptions) {
		o.servers = srvs
	})
}

func WithLogLevel(l slog.Leveler) AppOption {
	return newFuncAppOption(func(o *appOptions) {
		o.logLevel = l
	})
}

func New(opt ...AppOption) *App {
	opts := defaultAppOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	ctx, cancel := context.WithCancel(opts.ctx)
	a := &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   opts,
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: opts.logLevel})
	slog.SetDefault(slog.New(handler))
	return a
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(a.ctx, a.opts.sigs...)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	for _, server := range a.opts.servers {
		eg.Go(func() error {
			<-ctx.Done()
			sctx, cancel := context.WithTimeout(a.opts.ctx, a.opts.stopTimeout)
			defer cancel()
			return server.Stop(sctx)
		})
		eg.Go(func() error {
			return server.Start(a.opts.ctx)
		})
	}
	return eg.Wait()
}

func (a *App) Stop() (err error) {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
