package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

type testEnv struct {
	args     string
	out, err bytes.Buffer
}

func (te *testEnv) run() error {
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.SetOutput(&te.err)
	return run(strings.Fields(te.args), &te.out, flagSet)
}

func TestRunHappy(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args string
		out  string
	}{
		"url_only": {
			args: "http://google.com",
			out:  fmt.Sprintf("&{c:%d n:100 url:http://google.com}", runtime.NumCPU()),
		},
		"url+c+n": {
			args: "-c 5 -n 20 http://google.com",
			out:  "&{c:5 n:20 url:http://google.com}",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			te := testEnv{
				args: tt.args,
			}
			if err := te.run(); err != nil {
				t.Fatal(err)
			}
			if got, want := te.out.String(), tt.out; !strings.Contains(got, want) {
				t.Fatalf("Want %q to be part of %q\n", want, got)
			}
		})
	}
}

func TestRunSad(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args string
		err  string
	}{
		"no_url": {
			args: "",
			err:  `invalid value "" of URL: required`,
		},
		"too_much_concurrency": {
			args: fmt.Sprintf("-c %d http://google.com", 101),
			err:  "value of `c` must not be greater than value of `n`",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			te := testEnv{
				args: tt.args,
			}
			err := te.run()
			if err == nil {
				t.Fatal("error not reported")
			}
			if got, want := err.Error(), tt.err; !strings.Contains(got, want) {
				t.Fatalf("Want %q to be part of %q\n", want, got)
			}
		})
	}
}
