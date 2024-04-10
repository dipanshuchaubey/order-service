package handler

import (
	"context"
	"fmt"

	"order-service/internal/constants"
	"order-service/internal/data"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type SyncOrderHandler struct {
	log  *log.Helper
	repo data.OrdersRepository
}

func NewSyncOrderHandler(logger log.Logger, repo data.OrdersRepository) ISyncOrderHandler {
	return &SyncOrderHandler{log: log.NewHelper(logger), repo: repo}
}

func (h *SyncOrderHandler) Handler(ctx context.Context, messageID string, message MessageData) error {
	fmt.Println("SyncOrderHandler:: Received message: ", message)
	orders, err := h.repo.GetAllOrders(ctx)
	if err != nil {
		errMsg := fmt.Sprintf(constants.ErrorFetchingOrders, err)
		h.log.WithContext(ctx).Errorf(errMsg)
		return errors.New(500, constants.MySQLFetchError, errMsg)
	}

	fmt.Println("SyncOrderHandler:: Orders fetched successfully: ", orders)
	return nil
}
