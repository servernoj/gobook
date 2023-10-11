package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/servernoj/gobook/ch05/hit"
)

func run(args []string, out io.Writer, flagSet *flag.FlagSet) error {
	f := &flags{
		c:   runtime.NumCPU(),
		n:   100,
		rps: 0,
		m:   "GET",
		H:   []string{},
	}
	if err := f.parse(args, flagSet); err != nil {
		return err
	}
	fmt.Fprintf(out, "%+v\n", f)

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

	stat := client.Do()
	fmt.Println(stat)

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, flag.CommandLine); err != nil {
		os.Exit(1)
	}
}
