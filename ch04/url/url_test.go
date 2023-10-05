package url

import (
	"testing"
)

type Expectation func(url *URL, err error) string

func TestParse(t *testing.T) {

	tests := map[string]struct {
		input  string
		expect []Expectation
	}{
		"wrong url fails": {
			input: "wrong url",
			expect: []Expectation{
				func(url *URL, err error) string {
					if err == nil {
						return "err is not nil"
					}
					return ""
				},
			},
		},
		"scheme gets parsed": {
			input: "http://google.com",
			expect: []Expectation{
				func(url *URL, err error) string {
					if err != nil {
						return "err is nil"
					}
					return ""
				},
				func(url *URL, err error) string {
					if url != nil && url.Scheme != "http" {
						return "Scheme to be equal `http`"
					}
					return ""
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				violations := []string{}
				url, err := Parse(tt.input)
				for _, h := range tt.expect {
					if v := h(url, err); v != "" {
						violations = append(violations, v)
					}
				}
				if len(violations) > 0 {
					t.Logf("Parse(%q). Failed expections: %q\n", tt.input, violations)
					t.Fail()
				}
			},
		)
	}
}

func TestURLPort(t *testing.T) {
	tests := map[string]struct {
		input    *URL
		expected string
	}{
		"empty struct":               {&URL{}, ""},
		"missing port":               {&URL{Host: "google.com"}, ""},
		"numeric port over hostname": {&URL{Host: "google.com:80"}, "80"},
		"numeric port over ip4":      {&URL{Host: "1.2.3.4:80"}, "80"},
		"numeric port over ip6":      {&URL{Host: "2345:425:2CA1:0000:0000:567:5673:23b5:80"}, "80"},
		"service port over hostname": {&URL{Host: "google.com:http"}, "http"},
		"service port over ip6":      {&URL{Host: "2345:425:2CA1:0000:0000:567:5673:23b5:http"}, "http"},
	}
	for name, tt := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				if got, want := tt.input.Port(), tt.expected; got != want {
					t.Fatalf("%#v.Port(), got = %q, want = %q\n", *tt.input, got, want)
				}
			},
		)
	}
}

var tempGlobal string

func BenchmarkURLString(b *testing.B) {
	url := &URL{
		Scheme: "https",
		Host:   "google.com",
		Path:   "/?q=hello",
	}
	var tempLocal string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempLocal = url.String()
	}
	tempGlobal = tempLocal
}
