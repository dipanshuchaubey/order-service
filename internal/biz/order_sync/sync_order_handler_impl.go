package order_sync

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

type SyncOrderHandler struct {
	log *log.Helper
}

func NewSyncOrderHandler(logger log.Logger) ISyncOrderHandler {
	return &SyncOrderHandler{log: log.NewHelper(logger)}
}

func (h *SyncOrderHandler) Handler(ctx context.Context, messageID string, message MessageData) error {
	fmt.Println("SyncOrderHandler:: Received message: ", message)
	return nil
}
