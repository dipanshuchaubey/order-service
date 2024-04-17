package entity

import (
	v1 "order-service/api/v1/order"

	"github.com/oklog/ulid/v2"
)

const OrdersTableName = "orders"

func (*OrdersEntity) TableName() string {
	return OrdersTableName
}

type OrdersEntity struct {
	ID         string `gorm:"column:id;primaryKey"`
	CustomerID int64  `gorm:"column:customer_id"`
	CartID     int64  `gorm:"column:cart_id"`
	PaymentRef string `gorm:"column:payment_ref_no"`
	Status     string `gorm:"column:status"`
	ServerID   int64  `gorm:"column:server_id"`
	Base
}

func (entity *OrdersEntity) ToProto(data *v1.OrderData) {
	data.Id = entity.ID
	data.CartId = entity.CartID
	data.CustomerId = entity.CustomerID
	data.Status = entity.Status
	data.PaymentRef = entity.PaymentRef
	data.ServerId = entity.ServerID
	data.CreatedAt = entity.CreatedAt.String()
}

func (entity *OrdersEntity) FromCreateOrderRequest(data *v1.CreateOrderRequest) {
	entity.ID = ulid.Make().String()
	entity.CustomerID = data.CustomerId
	entity.PaymentRef = data.PaymentRef
	entity.CartID = data.CartId
	entity.Status = "CREATED"
	entity.ServerID = 1
}
