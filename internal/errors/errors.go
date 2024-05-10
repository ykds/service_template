package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

var errMap = map[int]Error{}

type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.code, e.message)
}

func NewError(code int, message string) Error {
	if _, ok := errMap[code]; ok {
		panic(fmt.Sprintf("错误码：%d 已定义", code))
	}
	e := Error{
		code:    code,
		message: message,
	}
	errMap[code] = e
	return e
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target error) bool {
	return errors.As(err, target)
}

func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func New(msg string) error {
	return errors.New(msg)
}
