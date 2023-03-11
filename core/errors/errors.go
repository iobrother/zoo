package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

//go:generate protoc --go_out=:. --go_opt=paths=source_relative errors.proto

func New(code int32, message, detail string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func Newf(code int32, message, format string, a ...any) *Error {
	return New(code, message, fmt.Sprintf(format, a...))
}

func Errorf(code int32, message, format string, a ...any) error {
	return New(code, message, fmt.Sprintf(format, a...))
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Code = 500
		e.Detail = err
	}
	if e.Code == 0 {
		e.Code = 500
	}
	return e
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if e := new(Error); errors.As(err, &e) {
		return e
	}

	return Parse(err.Error())
}

func (x *Error) Error() string {
	b, _ := json.Marshal(x)
	return string(b)
}

func Code(err error) int {
	if err == nil {
		return 200
	}
	return int(FromError(err).Code)
}
