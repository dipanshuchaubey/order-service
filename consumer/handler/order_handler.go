package handler

import (
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewSyncOrderHandler)

type ISyncOrderHandler interface {
	Handler(ctx context.Context, messageID string, message MessageData) error
}

type MessageData struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type MessageMeta struct {
	OrderID    string `json:"ID"`
	CartID     int64  `json:"CartID"`
	PaymentRef string `json:"PaymentRef"`
}

type MessageBody struct {
	Type           string `json:"Type"`
	MessageId      string `json:"MessageId"`
	SequenceNumber string `json:"SequenceNumber"`
	TopicArn       string `json:"TopicArn"`
	Subject        string `json:"Subject"`
	Message        string `json:"Message"`
	Timestamp      string `json:"Timestamp"`
}
