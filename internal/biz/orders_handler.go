package biz

import (
	"context"
	v1 "order-service/api/v1/order"
	"order-service/internal/biz/interfaces"
	"order-service/internal/data"
	"order-service/internal/data/entity"
	"order-service/internal/redis"

	"github.com/go-kratos/kratos/v2/log"
)

type OrdersHandler struct {
	repo  data.OrdersRepository
	redis redis.RedisHandlerInterface
	log   *log.Helper
}

func NewOrdersHandler(repo data.OrdersRepository, cache redis.RedisHandlerInterface, logger log.Logger) interfaces.OrdersHandlerInterface {
	return &OrdersHandler{repo, cache, log.NewHelper(logger)}
}

func (h *OrdersHandler) GetOrdersForUser(ctx context.Context, userID string) ([]*v1.OrderData, error) {
	res, err := h.repo.GetOrdersByCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var orders []*v1.OrderData
	for _, order := range res {
		info := v1.OrderData{}
		order.ToProto(&info)
		orders = append(orders, &info)
	}

	return orders, nil
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
	var order entity.OrdersEntity
	order.FromCreateOrderRequest(req)

	// Validate the request

	// Create the order
	createdOrder, err := h.repo.CreateOrder(ctx, &order)
	if err != nil {
		return nil, err
	}

	var orderData v1.OrderData
	createdOrder.ToProto(&orderData)

	// Cache the order
	cacheErr := h.redis.Set(ctx, orderData.Id, &orderData)
	if cacheErr != nil {
		h.log.Errorf("error caching order: %v", cacheErr)
	}

	return &v1.CreateOrderReply{Order: &orderData, Success: true}, nil
}
