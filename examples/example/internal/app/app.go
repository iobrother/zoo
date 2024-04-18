package app

import (
	"log"

	"github.com/google/wire"
	"github.com/iobrother/zoo/examples/example/internal/app/runtime"
)

var appContext *AppContext

var appContextProviderSet = wire.NewSet(
	runtime.RuntimeProviderSet,
	serviceProviderSet,
	wire.NewSet(wire.Struct(new(AppContext), "*")),
)

type AppContext struct {
	Runtime *runtime.Runtime
	Service *Service
}

func Init() error {
	var err error
	appContext, err = NewAppContext()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Context() *AppContext {
	return appContext
}
