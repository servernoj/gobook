package dynamicFlag

import (
	"errors"
	"regexp"
	"strings"
)

type Header []string

func (h *Header) Set(s string) error {
	re := regexp.MustCompile(`^\s*\S+\s*:\s*\S+\s*$`)
	if !re.MatchString(s) {
		return errors.New("invalid header pattern")
	}
	*h = append(*h, strings.TrimSpace(s))

	return nil
}

func (h *Header) String() string {
	return ""
}
