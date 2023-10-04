package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

func run(args []string, out io.Writer, flagSet *flag.FlagSet) error {
	f := &flags{
		c: runtime.NumCPU(),
		n: 100,
		t: 5 * time.Second,
		m: "GET",
		H: []string{},
	}
	if err := f.parse(args, flagSet); err != nil {
		return err
	}
	fmt.Fprintf(out, "%+v\n", f)
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, flag.CommandLine); err != nil {
		os.Exit(1)
	}
}
