package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/servernoj/gobook/url"
)

const usageText = `
Usage:
  %s [options] URL
Options:%s`

type flags struct {
	c   int
	n   int
	url string
}

func (f *flags) parse() error {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, os.Args[0], "\n")
		flag.PrintDefaults()
	}
	flag.Var(toNumber(&f.c), "c", "level of concurrency")
	flag.Var(toNumber(&f.n), "n", "total number of requests")
	flag.Parse()
	f.url = flag.Arg(0)

	if err := f.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
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
	case url.Scheme != "http":
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
