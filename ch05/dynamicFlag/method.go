package dynamicFlag

import (
	"errors"
	"slices"
	"strings"
)

type Method string

func (m *Method) Set(s string) error {
	allowed := []string{"GET", "POST", "PUT"}
	s = strings.ToUpper(strings.TrimSpace(s))
	if slices.Index(allowed, s) == -1 {
		return errors.New("unsupported method")
	}
	*m = Method(s)
	return nil
}

func (m *Method) String() string {
	return string(*m)
}
