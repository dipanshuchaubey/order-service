package order_sync

import "context"

type ISyncOrderHandler interface {
	Handler(ctx context.Context, messageID string, message MessageData) error
}

type MessageData struct {
	Event      string `json:"event"`
	OrderID    string `json:"order_id"`
	CartID     string `json:"cart_id"`
	PaymentRef string `json:"payment_ref"`
	OrderTime  int32  `json:"order_time"`
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
