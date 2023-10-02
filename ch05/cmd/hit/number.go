package main

import (
	"errors"
	"fmt"
	"strconv"
)

type number int

func toNumber(value *int) *number {
	return (*number)(value)
}

func (n *number) Set(s string) error {
	value, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return errors.New("parse error")
	}
	if value <= 0 {
		return fmt.Errorf("must be positive")
	}
	*n = number(value)
	return nil
}

func (n *number) String() string {
	return strconv.Itoa(int(*n))
}
