package entity

const OrdersTableName = "orders"

func (*OrdersEntity) TableName() string {
	return OrdersTableName
}

type OrdersEntity struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CustomerID int64  `gorm:"column:customer_id"`
	CartID     int64  `gorm:"column:cart_id"`
	PaymentRef string `gorm:"column:payment_ref_no"`
	Status     string `gorm:"column:status"`
	ServerID   int64  `gorm:"column:server_id"`
	Base
}
