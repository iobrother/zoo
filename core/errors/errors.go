package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Error struct {
	Code     int32             `json:"code,omitempty"`
	Message  string            `json:"message,omitempty"`
	Detail   string            `json:"detail,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func New(code int32, message, detail string) *Error {
	e := &Error{
		Code:     code,
		Message:  message,
		Detail:   detail,
		Metadata: map[string]string{},
	}

	e.Metadata["_zoo_error_stack"] = stacktrace()
	return e
}

func Newf(code int32, message, format string, a ...any) *Error {
	e := &Error{
		Code:     code,
		Message:  message,
		Detail:   fmt.Sprintf(format, a...),
		Metadata: map[string]string{},
	}

	e.Metadata["_zoo_error_stack"] = stacktrace()
	return e
}

func Errorf(code int32, message, format string, a ...any) error {
	e := &Error{
		Code:     code,
		Message:  message,
		Detail:   fmt.Sprintf(format, a...),
		Metadata: map[string]string{},
	}

	e.Metadata["_zoo_error_stack"] = stacktrace()
	return e
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Code = 500
		e.Message = "Internal Server Error"
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

func (x *Error) clone() *Error {
	if x == nil {
		return nil
	}

	metadata := make(map[string]string, len(x.Metadata))
	for k, v := range x.Metadata {
		metadata[k] = v
	}
	return &Error{
		Code:     x.Code,
		Message:  x.Message,
		Detail:   x.Detail,
		Metadata: metadata,
	}
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

func (x *Error) Format(s fmt.State, verb rune) {
	copied := x.clone()
	delete(copied.Metadata, "_zoo_error_stack")
	msg := fmt.Sprintf("error: code = %d message = %s detail = %s metadata = %v", x.Code, x.Message, x.Detail, copied.Metadata)

	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, msg+"\n"+x.Stack())

		default:
			_, _ = io.WriteString(s, msg)
		}
	}
}

func WrapRpcError(err error) *Error {
	if err == nil {
		return nil
	}

	e := FromError(err)
	if e.Metadata["_zoo_error_stack"] != "" {
		e.Metadata["_zoo_error_stack"] = stacktrace() + "\n" + e.Metadata["_zoo_error_stack"]
	}

	return e
}
