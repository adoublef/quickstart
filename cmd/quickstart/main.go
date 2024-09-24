package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
)

type cmd interface {
	parse(args []string, getenv func(string) string) error
	run(ctx context.Context) error
}

var cmds = map[string]cmd{
	"serve": cmdServe,
}

func main() {
	err := run(context.Background(), os.Args[1:], os.Getenv)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(2)
	} else if err != nil {
		// give a better error prefix here
		fmt.Fprintf(os.Stderr, "ERRO: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, getenv func(string) string) error {
	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}
	for name, c := range cmds {
		// todo: handle help
		// todo: handle version
		if cmd == name {
			err := c.parse(args, getenv)
			if err != nil {
				return err
			}
			return c.run(ctx)
		}
	}
	// wrap/handle properly
	return flag.ErrHelp
}
