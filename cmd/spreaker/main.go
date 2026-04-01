package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/G10xy/spreaker-and-go/internal/cli"
)


var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cli.Execute(ctx, version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
