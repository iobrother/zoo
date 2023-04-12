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

{{ range .Errors }}
func Is{{.CamelValue}}(err error) bool {
	e := errors.FromError(err)
	return e.Code == {{.Code}}
}

func Err{{.CamelValue}}(message ...string) *errors.Error {
	if len(message) > 0 {
	   return Err{{.CamelValue}}w(WithMessage(message[0]))
	}
    return Err{{.CamelValue}}w()
}

func Err{{.CamelValue}}f(format string, a ...any) *errors.Error {
	 return Err{{.CamelValue}}w(WithMessage(fmt.Sprintf(format, a...)))
}

func Err{{.CamelValue}}w(opt ...Option) *errors.Error {
    e := errors.NewWithStatusCode({{.StatusCode}}, {{.Code}}, "{{.Message}}", {{.Name}}_{{.Value}}.String())
	_apply(e, opt...)
	return e
}
{{- end }}