package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"

	zerrors "github.com/iobrother/zoo/core/errors"
)

var formDecoder *form.Decoder

func init() {
	formDecoder = form.NewDecoder()
	formDecoder.SetTagName("json")
}

type Context struct {
	*gin.Context
}

func (c *Context) SetError(err error) {
	c.Set("_zoo_error", err)
}

func (c *Context) GetError() error {
	v := c.Value("_zoo_error")
	if v == nil {
		return nil
	}
	err, ok := v.(error)
	if !ok {
		return nil
	}

	return err
}

func (c *Context) ShouldBind(v any) error {
	if err := c.Context.ShouldBind(v); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			return zerrors.NewWithStatusCode(400, 400, "Bad Request", err.Error())
		}
		return err
	}
	return nil
}

func (c *Context) ShouldBindUri(v any) error {
	m := make(map[string][]string)
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}

	if err := formDecoder.Decode(v, m); err != nil {
		return err
	}

	if err := binding.Validator.ValidateStruct(v); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			return zerrors.NewWithStatusCode(400, 400, "Bad Request", err.Error())
		}
		return err
	}

	return nil
}

func (c *Context) ShouldBindQuery(v any) error {
	values := c.Request.URL.Query()
	if err := formDecoder.Decode(v, values); err != nil {
		return err
	}

	if err := binding.Validator.ValidateStruct(v); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			return zerrors.NewWithStatusCode(400, 400, "Bad Request", err.Error())
		}
		return err
	}

	return nil
}

type HandlerFunc func(*Context)

func (s *Server) UseEx(middlewares ...HandlerFunc) {
	if len(middlewares) == 0 {
		return
	}
	handlers := make([]gin.HandlerFunc, 0, len(middlewares))
	for _, h := range middlewares {
		handler := func(c *gin.Context) {
			h(&Context{Context: c})
		}
		handlers = append(handlers, handler)
	}

	s.Use(handlers...)
}

func (s *Server) HandleEx(method, path string, handlers ...HandlerFunc) *Server {
	if len(handlers) == 0 {
		return s
	}
	gHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		handler := func(c *gin.Context) {
			h(&Context{Context: c})
		}
		gHandlers = append(gHandlers, handler)
	}

	s.Handle(method, path, gHandlers...)
	return s
}

func (s *Server) PostEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPost, path, handlers...)
}

func (s *Server) GetEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodGet, path, handlers...)
}

func (s *Server) DeleteEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodDelete, path, handlers...)
}

func (s *Server) PatchEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPatch, path, handlers...)
}

func (s *Server) PutEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPut, path, handlers...)
}

func (s *Server) OptionsEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodOptions, path, handlers...)
}

func (s *Server) HeadEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodHead, path, handlers...)
}

func (s *Server) AnyEx(path string, handlers ...HandlerFunc) *Server {
	anyMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
	for _, method := range anyMethods {
		s.HandleEx(method, path, handlers...)
	}

	return s
}

func (s *Server) MatchEx(methods []string, path string, handlers ...HandlerFunc) *Server {
	for _, method := range methods {
		s.HandleEx(method, path, handlers...)
	}

	return s
}
