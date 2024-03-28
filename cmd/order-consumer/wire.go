//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"order-service/internal/biz/order_sync"
	"order-service/internal/conf"
	"order-service/internal/consumer"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Data, *conf.Consumer, *conf.Consumer_Queue, log.Logger) (*consumer.OrderConsumer, func(), error) {
	panic(wire.Build(order_sync.ProviderSet, consumer.ProviderSet))
}
