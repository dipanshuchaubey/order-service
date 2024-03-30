package service

import (
	"context"

	v1 "order-service/api/v1/order"
	"order-service/internal/biz"
)

type OrderService struct {
	v1.UnimplementedOrderServer

	handler *biz.OrdersHandler
}

func NewOrderService(handler *biz.OrdersHandler) *OrderService {
	return &OrderService{handler: handler}
}

func (s *OrderService) GetAllOrders(ctx context.Context, req *v1.GetAllOrdersForUserRequest) (*v1.GetAllOrdersForUserReply, error) {
	res, err := s.handler.GetOrdersForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetAllOrdersForUserReply{Orders: res}, nil
}
