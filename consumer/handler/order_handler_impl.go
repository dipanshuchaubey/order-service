package handler

import (
	"context"
	"fmt"

	"order-service/internal/biz/interfaces"
	"order-service/internal/constants"
	"order-service/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

type SyncOrderHandler struct {
	log  *log.Helper
	repo data.OrdersRepository
	osvc interfaces.OrdersHandlerInterface
}

func NewSyncOrderHandler(logger log.Logger, repo data.OrdersRepository, osvc interfaces.OrdersHandlerInterface) ISyncOrderHandler {
	return &SyncOrderHandler{log.NewHelper(logger), repo, osvc}
}

func (h *SyncOrderHandler) Handler(ctx context.Context, messageID string, message MessageData) error {
	h.log.WithContext(ctx).Info("SyncOrderHandler:: Received message: ", message)

	switch message.Event {
	case constants.OrderPlaced:
		err := h.handleOrderPlaced(ctx, message)
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func (h *SyncOrderHandler) handleOrderPlaced(ctx context.Context, message MessageData) error {
	h.log.WithContext(ctx).Info("SyncOrderHandler:: handleOrderPlaced :: ", message.OrderID)

	_, err := h.osvc.UpdateOrder(ctx, message.OrderID)
	if err != nil {
		errMsg := fmt.Sprintf("SyncOrderHandler:: handleOrderPlaced :: Error updating order: %s", err.Error())
		h.log.WithContext(ctx).Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	h.log.WithContext(ctx).Info("SyncOrderHandler:: handleOrderPlaced :: Order updated successfully for orderID: ", message.OrderID)
	return nil
}
