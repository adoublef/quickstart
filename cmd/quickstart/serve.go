package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"time"

	"go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/time/rate"
	"golang.org/x/sync/errgroup"
)

// Limits as recommended by Cloudflare
// source: https://developers.cloudflare.com/workers/platform/limits/
const (
	maxHeaderBytes = 32 * (1 << 10)
)

type serve struct {
	addr string
	rate rate.Rate
	// todo: tls
}

func (c *serve) parse(args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&c.addr, "addr", "0.0.0.0:0", "http listening port")
	// Cloudflare sets a 1000/min rate limit default
	fs.TextVar(&c.rate, "rate", rate.Rate{N: 1000, D: time.Minute}, "api rate limit")
	// throttle safe requests and limit non-safe requests
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 0 {
		// todo: print usage?
		return flag.ErrHelp
	}
	return nil
}

func (c *serve) run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	s := &http.Server{
		Addr:           c.addr,
		Handler:        http.Handler(c.rate.N, c.rate.D),
		BaseContext:    func(l net.Listener) context.Context { return ctx },
		MaxHeaderBytes: maxHeaderBytes,
		// todo: timeouts
	}
	s.RegisterOnShutdown(cancel)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		switch {
		case s.TLSConfig != nil:
			err = s.ListenAndServeTLS("", "")
		default:
			err = s.ListenAndServe()
		}
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			err = errors.Join(s.Close())
		}
		return err
	})

	return eg.Wait()
}
