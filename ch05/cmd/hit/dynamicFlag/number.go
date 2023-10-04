package dynamicFlag

import (
	"errors"
	"fmt"
	"strconv"
)

type Number int

func (n *Number) Set(s string) error {
	value, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return errors.New("parse error")
	}
	if value <= 0 {
		return fmt.Errorf("must be positive")
	}
	*n = Number(value)
	return nil
}

func (n *Number) String() string {
	return strconv.Itoa(int(*n))
}
