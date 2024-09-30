package main

import (
	"context"
	"flag"
	"testing"
	"time"

	"go.tmp/quickstart/internal/testing/is"
	"go.tmp/quickstart/internal/testing/wait"
	"golang.org/x/sync/errgroup"
)

func Test_serve_parse(t *testing.T) {
	type testcase struct {
		in   []string
		want error
	}

	var tt = map[string]testcase{
		"OKRate": {
			in: []string{"-rate-limit", "1/10s"},
		},
		"ErrTooManyArgs": {
			in:   []string{"never"},
			want: flag.ErrHelp,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := (&serve{}).parse(tc.in, nil)
			is.NotOK(t, err, tc.want) // got;want
		})
	}
}

func Test_serve_run(t *testing.T) {
	t.Run("OKGracefulShutdown", func(t *testing.T) {
		var s serve
		// random port?
		err := s.parse([]string{"-address", ":8080"}, nil)
		is.OK(t, err)
		// cancellable context is needed
		ctx, cancel := context.WithCancel(context.Background())
		// errgroup makes handling errors in multiple go routines easier
		eg, ctx := errgroup.WithContext(ctx)
		// 1. start service
		eg.Go(func() error { return s.run(ctx) })
		ready := func() error {
			defer cancel()
			err := wait.ForHTTP(ctx, 30*time.Second, "http://localhost:8080/ready")
			if err != nil {
				return err
			}
			return nil
		}
		eg.Go(ready)
		is.OK(t, eg.Wait()) // service is ready
	})
}
