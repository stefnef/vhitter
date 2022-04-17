package errorDto

import "errors"

type ErrorConnection struct {
	Msg   string
	Cause error
}

func (e *ErrorConnection) Error() string {
	return e.Msg
}

func (m *ErrorConnection) Is(target error) bool { return target.Error() == m.Msg }

var (
	ErrorNotFound     = errors.New("not found")
	ErrorUnauthorized = errors.New("unauthorized")
)
