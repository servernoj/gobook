package hit

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Result struct {
	Duration time.Duration
	Bytes    int
	Err      error
}

type Stat struct {
	Count       int
	Slowest     time.Duration
	Fastest     time.Duration
	Duration    time.Duration
	Errors      int
	RPS         float64
	SuccessRate float64
	Bytes       int
}

func (s *Stat) Process(r *Result) {
	s.Count++
	if s.Fastest == 0 || (r.Duration > 0 && r.Duration < s.Fastest) {
		s.Fastest = r.Duration
	}
	if r.Duration > s.Slowest {
		s.Slowest = r.Duration
	}
	if r.Err != nil {
		s.Errors++
	}
	s.Bytes += r.Bytes
}

func (s *Stat) PostProcess(d time.Duration) {
	s.RPS = float64(s.Count) / d.Seconds()
	s.SuccessRate = float64(s.Count-s.Errors) / float64(s.Count)
	s.Duration = d.Round(time.Millisecond)
}

func (s *Stat) Fprint(out io.Writer) {
	h := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	h("%-30s%d\n", "Requests sent:", s.Count)
	h("%-30s%.1f\n", "RPS:", s.RPS)
	h("%-30s%s\n", "Duration:", s.Duration)
	h("%-30s%d\n", "Bytes received:", s.Bytes)
	h("%-30s%.0f%%\n", "Success rate:", s.SuccessRate*100)
	if s.Count > 1 {
		h("%-30s%d\n", "Slowest request, ms:", s.Slowest.Milliseconds())
		h("%-30s%d\n", "Fastest request, ms:", s.Fastest.Milliseconds())
	}
}

func (s *Stat) String() string {
	var sb strings.Builder
	s.Fprint(&sb)
	return sb.String()
}
