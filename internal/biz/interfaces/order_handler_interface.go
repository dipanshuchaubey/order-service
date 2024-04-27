package interfaces

import (
	"context"
	v1 "order-service/api/v1/order"
)

type OrdersHandlerInterface interface {
	GetOrdersForUser(ctx context.Context, userID string) ([]*v1.OrderData, error)
	CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderReply, error)
	UpdateOrder(ctx context.Context, orderID string) (*v1.CreateOrderReply, error)
}
