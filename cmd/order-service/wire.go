//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"order-service/internal/biz"
	"order-service/internal/conf"
	"order-service/internal/data"
	"order-service/internal/publisher"
	"order-service/internal/redis"
	"order-service/internal/server"
	"order-service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Publisher, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, redis.ProviderSet, publisher.ProviderSet, newApp))
}
