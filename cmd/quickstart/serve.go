package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/time/rate"
	"golang.org/x/sync/errgroup"
)

var cmdServe = &serve{}

type serve struct {
	addr string
	rate rate.Rate
	// todo: tls
}

func (c *serve) parse(args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&c.addr, "address", "0.0.0.0:0", "http listening port")
	// Cloudflare sets a 1000/min rate limit default
	fs.TextVar(&c.rate, "rate-limit", rate.Rate{N: 1000, D: time.Minute}, "api rate limit")
	// throttle safe requests and limit non-safe requests
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), `
The serve command initialises and runs a HTTP server.

Usage:
	%s serve [arguments]

Arguments:
`[1:], os.Args[0])
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	return nil
}

func (c *serve) run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// this can be handle by a flag?
	shutdown, err := setupOTel(ctx)
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:           c.addr,
		Handler:        http.Handler(c.rate.N, c.rate.D),
		BaseContext:    func(l net.Listener) context.Context { return ctx },
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		// todo: ReadHeaderTimeout uses ReadTimeout if not set
		ReadTimeout:  http.DefaultReadTimeout,
		WriteTimeout: http.DefaultWriteTimeout,
		IdleTimeout:  http.DefaultIdleTimeout,
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

		return shutdown(ctx)
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
