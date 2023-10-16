package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/servernoj/gobook/ch04/url"
	"github.com/servernoj/gobook/ch05/dynamicFlag"
)

const usageText = `
Usage:
  %s [options] URL
Options:
`

type flags struct {
	c       int
	n       int
	url     string
	m       string
	H       []string
	rps     int
	timeout time.Duration
}

func (f *flags) parse(args []string, flagSet *flag.FlagSet) error {

	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), usageText, os.Args[0])
		flagSet.PrintDefaults()
	}

	flagSet.Var((*dynamicFlag.Number)(&f.c), "c", "level of concurrency")
	flagSet.Var((*dynamicFlag.Number)(&f.rps), "r", "make up to requests per second")
	flagSet.Var((*dynamicFlag.Number)(&f.n), "n", "total number of requests")
	flagSet.Var((*dynamicFlag.Method)(&f.m), "m", "HTTP request method")
	flagSet.Var((*dynamicFlag.Header)(&f.H), "H", "request header")
	flagSet.DurationVar(&f.timeout, "t", f.timeout, "time allowance for all requests")

	if err := flagSet.Parse(args); err != nil {
		return err
	}
	f.url = flagSet.Arg(0)
	if err := f.validate(); err != nil {
		fmt.Fprintln(flagSet.Output(), err)
		flagSet.Usage()
		return err
	}
	return nil
}

func validateUrl(rawurl string) error {
	if strings.TrimSpace(rawurl) == "" {
		return errors.New("required")
	}
	url, err := url.Parse(rawurl)
	switch {
	case err != nil:
		{
			return errors.New("parse error")
		}
	case url.Scheme != "http" && url.Scheme != "https":
		{
			return errors.New("invalid scheme (must be 'http')")
		}
	case url.Host == "":
		{
			return errors.New("missing host")
		}
	}
	return nil
}

func (f *flags) validate() error {

	if urlErr := validateUrl(f.url); urlErr != nil {
		return fmt.Errorf("invalid value %q of URL: %w", f.url, urlErr)
	}
	if f.c > f.n {
		return errors.New("value of `c` must not be greater than value of `n`")
	}
	return nil
}
