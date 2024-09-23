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
	"golang.org/x/sync/errgroup"
)

type serve struct {
	addr string
	// tls
}

func (s *serve) parse(args []string, _ func(string) string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&s.addr, "addr", "0.0.0.0:0", "http listening port")
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	return nil
}

func (c *serve) run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	s := &http.Server{
		Addr:        c.addr,
		Handler:     http.Handler(),
		BaseContext: func(l net.Listener) context.Context { return ctx },
		// todo: timeouts
		// todo: maxHeaderBytes
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
