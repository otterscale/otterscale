package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *App {
	o := options{
		ctx:      context.Background(),
		sigs:     []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL},
		timeout:  10 * time.Second,
		logLevel: slog.LevelDebug,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	// stderr: log
	// stdout: fmt.Print*
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: o.logLevel})
	slog.SetDefault(slog.New(handler))
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

func (a *App) Run() error {
	nctx, stop := signal.NotifyContext(a.ctx, a.opts.sigs...)
	defer stop()

	eg, ctx := errgroup.WithContext(nctx)
	for _, server := range a.opts.servers {
		// https://go.dev/blog/loopvar-preview
		eg.Go(func() error {
			<-ctx.Done()
			sctx, cancel := context.WithTimeout(a.opts.ctx, a.opts.timeout)
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
