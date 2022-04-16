package errorDto

import "errors"

var (
	ErrorNotFound     = errors.New("not found")
	ErrorUnauthorized = errors.New("unauthorized")
)
