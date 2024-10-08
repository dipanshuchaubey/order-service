package data

import (
	"context"
	"strconv"

	"order-service/internal/data/entity"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ordersRepository struct {
	db  *gorm.DB
	log *log.Helper
}

type OrdersRepository interface {
	GetAllOrders(ctx context.Context) ([]*entity.OrdersEntity, error)
	GetOrdersByCustomerID(ctx context.Context, customerID string) ([]*entity.OrdersEntity, error)
	GetOrderDetails(ctx context.Context, orderID string) (*entity.OrdersEntity, error)
	CreateOrder(context.Context, *entity.OrdersEntity) (*entity.OrdersEntity, error)
}

func NewOrdersRepository(data *Data, logger log.Logger) OrdersRepository {
	return &ordersRepository{data.db, log.NewHelper(logger)}
}

// GetAllOrders:: Get all orders
func (repo *ordersRepository) GetAllOrders(ctx context.Context) ([]*entity.OrdersEntity, error) {
	var orders []*entity.OrdersEntity
	if err := repo.db.Table(entity.OrdersTableName).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrdersByCustomerID:: Get all orders for a customer
func (repo *ordersRepository) GetOrdersByCustomerID(ctx context.Context, customerID string) ([]*entity.OrdersEntity, error) {
	customerIDInt64, convertErr := strconv.ParseInt(customerID, 10, 64)
	if convertErr != nil {
		return nil, convertErr
	}

	var orders []*entity.OrdersEntity
	if err := repo.db.Table(entity.OrdersTableName).Where("customer_id = ?", customerIDInt64).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderDetails:: Get order info by order ID
func (repo *ordersRepository) GetOrderDetails(ctx context.Context, orderID string) (*entity.OrdersEntity, error) {
	var order entity.OrdersEntity
	if err := repo.db.Table(entity.OrdersTableName).Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateOrder:: Place an order and set order status to 'PLACED'
func (repo *ordersRepository) CreateOrder(ctx context.Context, order *entity.OrdersEntity) (*entity.OrdersEntity, error) {
	var res entity.OrdersEntity
	if err := repo.db.Table(entity.OrdersTableName).Create(&order).Scan(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
