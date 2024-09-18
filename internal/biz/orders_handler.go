package biz

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "order-service/api/v1/order"
	"order-service/internal/biz/interfaces"
	"order-service/internal/constants"
	"order-service/internal/data"
	"order-service/internal/data/entity"
	"order-service/internal/publisher"
	"order-service/internal/redis"
	"order-service/internal/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type OrdersHandler struct {
	repo      data.OrdersRepository
	redis     redis.RedisHandlerInterface
	publisher publisher.PublisherInterface
	log       *log.Helper
}

func NewOrdersHandler(
	repo data.OrdersRepository,
	cache redis.RedisHandlerInterface,
	publisher publisher.PublisherInterface,
	logger log.Logger,
) interfaces.OrdersHandlerInterface {
	return &OrdersHandler{repo, cache, publisher, log.NewHelper(logger)}
}

func (h *OrdersHandler) GetOrdersForUser(ctx context.Context, userID string) ([]*v1.OrderData, error) {
	ctx, span := utils.Trace(ctx, "biz.GetOrdersForUser")
	defer span.End()

	h.log.WithContext(ctx).Infof("GetOrdersForUser:: Fetching orders for userID: %s", userID)

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
	ctx, span := utils.Trace(ctx, "biz.CreateOrder")
	defer span.End()

	h.log.WithContext(ctx).Infof("CreateOrder:: Create order request: %v", req)

	var order entity.OrdersEntity
	order.FromCreateOrderRequest(req)

	// TODO: Validate the request
	var orderData v1.OrderData
	order.ToProto(&orderData)

	// Cache the order
	valueBytes, marErr := json.Marshal(&order)
	if marErr != nil {
		errMsg := fmt.Sprintf(constants.ErrorMarshalling, marErr)
		h.log.WithContext(ctx).Error(errMsg)
		return nil, errors.New(500, constants.MarshalError, errMsg)
	}

	cacheErr := h.redis.Set(ctx, orderData.Id, valueBytes)
	if cacheErr != nil {
		h.log.Errorf("error caching order: %v", cacheErr)
	}

	// Publish OrderCreated event
	pubErr := h.publisher.PublishOrderEvents(ctx, constants.OrderPlaced, string(valueBytes), order.ID)
	if pubErr != nil {
		return nil, pubErr
	}

	h.log.WithContext(ctx).Infof("CreateOrder:: Order created successfully for orderID: %s", orderData.Id)
	return &v1.CreateOrderReply{Order: &orderData, Success: true}, nil
}

func (h *OrdersHandler) UpdateOrder(ctx context.Context, orderID string) (*v1.CreateOrderReply, error) {
	// Get order details from Cache
	cachedOrder, cacheErr := h.redis.Get(ctx, orderID)
	if cacheErr != nil {
		return nil, cacheErr
	}

	var order entity.OrdersEntity
	unMarErr := json.Unmarshal([]byte(cachedOrder), &order)
	if unMarErr != nil {
		h.log.Errorf("UpdateOrder:: %s", constants.ErrorUnmarshalling)
		return nil, unMarErr
	}

	order.Status = constants.OrderConfirmed

	// Insert order into DB
	h.log.WithContext(ctx).Infof("UpdateOrder:: Updating order for orderID: %s :: %v", orderID, order)
	orderDetails, dbErr := h.repo.CreateOrder(ctx, &order)

	if dbErr != nil {
		errMsg := fmt.Sprintf("UpdateOrder :: %s :: %v", constants.ErrorCreatingOrder, dbErr)
		h.log.Errorf(errMsg)
		return nil, errors.New(500, constants.MySQLError, errMsg)
	}

	// Update order in Cache
	var orderData v1.OrderData
	orderDetails.ToProto(&orderData)

	h.log.WithContext(ctx).Infof("UpdateOrder:: Order updated successfully for orderID: %s", orderID)
	return &v1.CreateOrderReply{Order: &orderData, Success: true}, nil
}
