package handler

import (
	"context"
	"encoding/json"
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
	payload, payloadErr := h.extractDataFromMessage(message)
	if payloadErr != nil {
		errMsg := fmt.Sprintf("SyncOrderHandler:: handleOrderPlaced :: Error extracting data from message: %s", payloadErr.Error())
		h.log.WithContext(ctx).Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	h.log.WithContext(ctx).Info("SyncOrderHandler:: handleOrderPlaced :: ", payload.OrderID)

	_, err := h.osvc.UpdateOrder(ctx, payload.OrderID)
	if err != nil {
		errMsg := fmt.Sprintf("SyncOrderHandler:: handleOrderPlaced :: Error updating order: %s", err.Error())
		h.log.WithContext(ctx).Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}

func (h *SyncOrderHandler) extractDataFromMessage(message MessageData) (MessageMeta, error) {
	var data MessageMeta

	err := json.Unmarshal([]byte(message.Data), &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
