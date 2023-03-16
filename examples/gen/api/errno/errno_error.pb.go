// Code generated by protoc-gen-zoo-errno. DO NOT EDIT.
// versions:
// - protoc-gen-zoo-errno v0.1.0
// - protoc                  (unknown)
// source: errno/errno.proto

package errno

import (
	fmt "fmt"
	errors "github.com/iobrother/zoo/core/errors"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = fmt.Errorf
var _ = errors.New

type Option interface {
	apply(*errors.Error)
}

type optFunc func(e *errors.Error)

func (o optFunc) apply(e *errors.Error) { o(e) }

func WithMessage(s string) Option {
	return optFunc(func(e *errors.Error) {
		if s != "" {
			e.Message = s
		}
	})
}

func WithDetail(s string) Option {
	return optFunc(func(e *errors.Error) {
		if s != "" {
			e.Detail = s
		}
	})
}

func WithMetadata(k string, v string) Option {
	return optFunc(func(e *errors.Error) {
		if k != "" && v != "" {
			e.Metadata[k] = v
		}
	})
}

func _apply(e *errors.Error, opts ...Option) {
	for _, opt := range opts {
		opt.apply(e)
	}
}

func IsInternalServerError(err error) bool {
	e := errors.FromError(err)
	return e.Code == 500
}

func ErrInternalServerError(detail string) *errors.Error {
	return errors.New(500, "服务器错误", detail)
}

func ErrInternalServerErrorf(format string, a ...any) *errors.Error {
	return errors.New(500, "服务器错误", fmt.Sprintf(format, a...))
}

func ErrInternalServerErrorw(opt ...Option) *errors.Error {
	e := errors.New(500, "服务器错误", ErrorReason_INTERNAL_SERVER_ERROR.String())
	_apply(e, opt...)
	return e
}
func IsDbError(err error) bool {
	e := errors.FromError(err)
	return e.Code == 100101
}

func ErrDbError(message ...string) *errors.Error {
	if len(message) > 0 {
		return ErrDbErrorw(WithMessage(message[0]))
	}
	return ErrDbErrorw()
}

func ErrDbErrorf(format string, a ...any) *errors.Error {
	return ErrDbErrorw(WithMessage(fmt.Sprintf(format, a...)))
}

func ErrDbErrorw(opt ...Option) *errors.Error {
	e := errors.New(100101, "数据库错误", ErrorReason_DB_ERROR.String())
	_apply(e, opt...)
	return e
}
func IsOrderNotExist(err error) bool {
	e := errors.FromError(err)
	return e.Code == 100201
}

func ErrOrderNotExist(message ...string) *errors.Error {
	if len(message) > 0 {
		return ErrOrderNotExistw(WithMessage(message[0]))
	}
	return ErrOrderNotExistw()
}

func ErrOrderNotExistf(format string, a ...any) *errors.Error {
	return ErrOrderNotExistw(WithMessage(fmt.Sprintf(format, a...)))
}

func ErrOrderNotExistw(opt ...Option) *errors.Error {
	e := errors.New(100201, "订单不存在", ErrorReason_ORDER_NOT_EXIST.String())
	_apply(e, opt...)
	return e
}
