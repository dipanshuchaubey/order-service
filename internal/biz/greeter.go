package biz

import (
	"context"

	v1 "order-service/api/v1/helloworld"
	"order-service/internal/data"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

type GreeterUsecase struct {
	repo data.OrdersRepository
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo data.OrdersRepository, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo, log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return &Greeter{Hello: g.Hello}, nil
}

func (uc *GreeterUsecase) SayHello(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Info("SayHello")
	return &Greeter{Hello: g.Hello}, nil
}
