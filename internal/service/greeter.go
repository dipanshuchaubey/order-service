package service

import (
	"context"

	v1 "order-service/api/helloworld/v1"
	"order-service/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *GreeterService) HealthCheck(ctx context.Context, in *v1.HealthCheckRequest) (*v1.HealthCheckReply, error) {
	return &v1.HealthCheckReply{Status: "Service is up and running!!"}, nil
}
