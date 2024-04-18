//go:generate wire ./...
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

func NewAppContext() (*AppContext, error) {
	panic(wire.Build(appContextProviderSet))
	return new(AppContext), nil
}
