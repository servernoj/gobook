package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/servernoj/gobook/ch05/hit"
)

func run(args []string, out io.Writer, flagSet *flag.FlagSet) error {
	f := &flags{
		c:       runtime.NumCPU(),
		n:       100,
		rps:     0,
		m:       "GET",
		H:       []string{},
		timeout: 10 * time.Minute,
	}
	if err := f.parse(args, flagSet); err != nil {
		return err
	}

	requestTemplate, err := http.NewRequest(f.m, f.url, http.NoBody)
	if err != nil {
		return err
	}
	client := hit.Client{
		RequestTemplate:  requestTemplate,
		NumberOfRequests: f.n,
		Concurrency:      f.c,
		RPS:              f.rps,
	}
	ctx, cancel := context.WithTimeout(
		context.Background(),
		f.timeout,
	)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()

	stat := client.Do(ctx)
	fmt.Println(stat)

	contextError := ctx.Err()
	if errors.Is(contextError, context.DeadlineExceeded) {
		return fmt.Errorf("timed out in %s", f.timeout)
	}

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, flag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
