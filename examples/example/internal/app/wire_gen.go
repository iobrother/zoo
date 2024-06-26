// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/iobrother/zoo/core/config"
	"github.com/iobrother/zoo/examples/example/internal/app/runtime"
	"github.com/iobrother/zoo/examples/example/internal/service/example"
	"github.com/iobrother/zoo/examples/example/pkg/database/db"
	"github.com/iobrother/zoo/examples/example/pkg/database/redis"
)

// Injectors from wire.go:

func NewAppContext() (*AppContext, error) {
	configConfig := config.Default()
	runtimeConfig, err := runtime.NewConfig(configConfig)
	if err != nil {
		return nil, err
	}
	dbConfig := &runtimeConfig.Mysql
	gormDB, err := db.Open(dbConfig)
	if err != nil {
		return nil, err
	}
	redisConfig := &runtimeConfig.Redis
	client, err := redis.NewClient(redisConfig)
	if err != nil {
		return nil, err
	}
	runtimeRuntime := &runtime.Runtime{
		Config: runtimeConfig,
		DB:     gormDB,
		RC:     client,
	}
	exampleExample := example.New(runtimeRuntime)
	service := &Service{
		Example: exampleExample,
	}
	appAppContext := &AppContext{
		Runtime: runtimeRuntime,
		Service: service,
	}
	return appAppContext, nil
}
