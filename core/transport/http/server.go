package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http/middleware/logging"
	"github.com/iobrother/zoo/core/transport/http/middleware/tracing"
)

type Server struct {
	opts Options
	*gin.Engine
	server *http.Server
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	srv := &Server{
		opts: options,
	}

	gin.SetMode(srv.opts.Mode)
	r := gin.New()
	srv.Engine = r
	srv.server = &http.Server{Handler: srv.Engine}
	return srv
}

func (s *Server) Init(opts ...Option) {
	for _, opt := range opts {
		opt(&s.opts)
	}
}

func (s *Server) Start() error {
	if s.opts.Tracing {
		s.Use(tracing.Trace(s.opts.Name))
	}
	s.Use(logging.Log())

	if s.opts.InitHttpServer != nil {
		if err := s.opts.InitHttpServer(s); err != nil {
			return err
		}
	}

	l, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	a := l.Addr().String()
	log.Infof("Server [GIN] listening on %s", a)
	go func() {
		if err := s.server.Serve(l); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err.Error())
			}
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
