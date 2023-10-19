package shortener

import (
	"errors"
)

var ErrTest = errors.New("test error")
var ErrNotFound = errors.New("key not found")
var ErrAlreadyExists = errors.New("key already exists")
var ErrInvalidMethod = errors.New("invalid method")
var ErrBadRequest = errors.New("bad request")
