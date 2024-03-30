package biz

import (
	"context"
	v1 "order-service/api/v1/order"
	"order-service/internal/data"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
)

type OrdersHandler struct {
	repo data.OrdersRepository
	log  *log.Helper
}

func NewOrdersHandler(repo data.OrdersRepository, logger log.Logger) *OrdersHandler {
	return &OrdersHandler{repo, log.NewHelper(logger)}
}

func (h *OrdersHandler) GetOrdersForUser(ctx context.Context, userID string) ([]*v1.OrderData, error) {
	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	res, err := h.repo.GetOrdersByCustomerID(ctx, userIDInt64)
	if err != nil {
		return nil, err
	}

	var orders []*v1.OrderData
	for _, order := range res {
		orders = append(orders, &v1.OrderData{
			Id:         order.ID,
			CartId:     order.CartID,
			CustomerId: order.CustomerID,
			Status:     order.Status,
			PaymentRef: order.PaymentRef,
			ServerId:   order.ServerID,
			CreatedAt:  order.CreatedAt.String(),
		})
	}

	return orders, nil
}
