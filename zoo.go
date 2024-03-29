package zoo

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/iobrother/zoo/core/config"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/registry"
	"github.com/iobrother/zoo/core/registry/etcd"
	httpserver "github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/rpc/server"
	"github.com/iobrother/zoo/core/util/env"
)

type App struct {
	opts       Options
	zc         *zconfig
	rpcServer  *server.Server
	httpServer *httpserver.Server
}

func (a *App) GetHttpServer() *httpserver.Server {
	return a.httpServer
}

func (a *App) GetRpcServer() *server.Server {
	return a.rpcServer
}

type zconfig struct {
	App struct {
		Mode string
		Name string
	}
	Logger struct {
		Level      string `json:"level"`
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
		Compress   bool   `json:"compress"`
	}
	Http struct {
		Addr string
	}
	Rpc struct {
		Addr string
	}
	Tracer struct {
		Addr string
	}
	Registry struct {
		BasePath       string
		EtcdAddr       []string
		UpdateInterval int
	}
}

func New(opts ...Option) *App {
	options := newOptions(opts...)
	zc := &zconfig{}
	if err := config.Unmarshal(zc); err != nil {
		log.Fatal(err)
	}
	fmt.Println(zc)
	if zc.App.Name == "" {
		log.Fatal("config item app.name can't be empty")
	}

	env.Set(zc.App.Mode)

	level, err := zapcore.ParseLevel(zc.Logger.Level)
	if err != nil {
		level = log.InfoLevel
	}
	if env.IsDevelop() {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.NewTee([]io.Writer{os.Stderr, w}, level, log.WithCaller(true), log.Development())
		log.ResetDefault(l)
	} else {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.New(w, level, log.WithCaller(true))
		log.ResetDefault(l)
	}

	app := &App{
		opts: options,
		zc:   zc,
	}

	tracing := false
	if zc.Tracer.Addr != "" {
		setTracerProvider(zc.Tracer.Addr, zc.App.Name)
		tracing = true
	}

	if app.opts.InitRpcServer != nil {
		app.rpcServer = server.NewServer(
			server.Name(zc.App.Name),
			server.Addr(zc.Rpc.Addr),
			server.BasePath(zc.Registry.BasePath),
			server.UpdateInterval(zc.Registry.UpdateInterval),
			server.EtcdAddr(zc.Registry.EtcdAddr),
			server.Tracing(tracing),
			server.InitRpcServer(app.opts.InitRpcServer),
		)
	}
	mode := "debug"
	if env.IsProduct() || env.IsStaging() {
		mode = "release"
	}

	if app.opts.InitHttpServer != nil {
		var r registry.Registry
		var opts []etcd.Option
		if zc.Registry.BasePath != "" {
			opts = append(opts, etcd.BasePath(zc.Registry.BasePath))
		}
		if len(zc.Registry.EtcdAddr) > 0 {
			opts = append(opts, etcd.Addrs(zc.Registry.EtcdAddr...))
			r = etcd.NewRegistry(opts...)
		}

		app.httpServer = httpserver.NewServer(
			httpserver.Name(zc.App.Name),
			httpserver.Addr(zc.Http.Addr),
			httpserver.Mode(mode),
			httpserver.Tracing(tracing),
			httpserver.Registry(r),
			httpserver.InitHttpServer(app.opts.InitHttpServer),
		)
	}

	return app
}

func setTracerProvider(endpoint string, name string) *trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Fatal(err)
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func (a *App) Run() error {
	for _, f := range a.opts.BeforeStart {
		if err := f(); err != nil {
			return err
		}
	}

	if a.rpcServer != nil {
		if err := a.rpcServer.Start(); err != nil {
			return err
		}
	}

	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return err
		}
	}

	for _, f := range a.opts.AfterStart {
		if err := f(); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	log.Infof("received signal %s", <-ch)

	for _, f := range a.opts.BeforeStop {
		if err := f(); err != nil {
			return err
		}
	}

	if a.rpcServer != nil {
		_ = a.rpcServer.Stop()
	}

	if a.httpServer != nil {
		_ = a.httpServer.Stop()
	}

	for _, f := range a.opts.AfterStop {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
