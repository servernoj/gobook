package url_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/servernoj/gobook/ch04/url"
)

func TestURLHostname(t *testing.T) {
	tests := map[string]struct {
		input    *url.URL
		expected string
	}{
		"nil struct":   {nil, ""},
		"empty struct": {&url.URL{}, ""},
	}
	for name, tt := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				if got, want := tt.input.Port(), tt.expected; got != want {
					t.Fatalf("%#v.Hostname(), got = %q, want = %q\n", *tt.input, got, want)
				}
			},
		)
	}
}

func ExampleURL() {
	u, err := url.Parse("http://google.com:80/?q=hello")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	fmt.Println(u)
	// output:
	// https://google.com:80/?q=hello
}

func ExampleURL_fields() {
	u, _ := url.Parse("http://google.com:80/?q=hello")
	fmt.Println(u.Port())
	fmt.Println(u.Hostname())
	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	fmt.Println(u.Path)
	// output:
	// 80
	// google.com
	// http
	// google.com:80
	// /?q=hello
}
