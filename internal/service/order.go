package service

import (
	"context"
	v1 "order-service/api/v1/order"
	"order-service/internal/biz/interfaces"
)

type OrderService struct {
	v1.UnimplementedOrderServer
	handler interfaces.OrdersHandlerInterface
}

func NewOrderService(handler interfaces.OrdersHandlerInterface) *OrderService {
	return &OrderService{handler: handler}
}

func (s *OrderService) GetAllOrders(ctx context.Context, req *v1.GetAllOrdersForUserRequest) (*v1.GetAllOrdersForUserReply, error) {
	res, err := s.handler.GetOrdersForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetAllOrdersForUserReply{Orders: res}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
	res, err := s.handler.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
