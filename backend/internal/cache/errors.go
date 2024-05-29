package cache

import "errors"

var ErrNotFound = errors.New("not found")
var ErrExpired = errors.New("expired")
var ErrMaxItems = errors.New("maximum number of items reached")
var ErrExceededOrderSize = errors.New("exceeded order size")
var ErrExceededKeySize = errors.New("exceeded key size")
