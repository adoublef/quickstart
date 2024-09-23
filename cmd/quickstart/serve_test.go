package main

import (
	"context"
	"testing"
	"time"

	"go.tmp/quickstart/internal/testing/is"
	"go.tmp/quickstart/internal/testing/wait"
	"golang.org/x/sync/errgroup"
)

func Test_serve_parse(t *testing.T) {}

func Test_serve_run(t *testing.T) {
	t.Run("OKGracefulShutdown", func(t *testing.T) {
		var s serve
		// random port?
		err := s.parse([]string{"-addr", ":8080"}, nil)
		is.OK(t, err)
		// cancellable context is needed
		ctx, cancel := context.WithCancel(context.Background())
		// errgroup makes handling errors in multiple go routines easier
		eg, ctx := errgroup.WithContext(ctx)
		// 1. start service
		eg.Go(func() error { return s.run(ctx) })
		// 1. wait for ready
		eg.Go(func() error {
			defer cancel() // could handle this elsewhere
			err := wait.ForHTTP(ctx, 30*time.Second, "http://localhost:8080/ready")
			if err != nil {
				return err
			}
			return nil
		})
		is.OK(t, eg.Wait()) // service is ready
	})
}
