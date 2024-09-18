package biz

import (
	"context"
	"fmt"
	v1 "order-service/api/v1/order"
	"order-service/internal/biz/interfaces"
	"order-service/internal/data/entity"
	"order-service/mocks"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupOrdersHandler(repo *mocks.OrdersRepository, redis *mocks.RedisHandlerInterface, pub *mocks.PublisherInterface) (context.Context, interfaces.OrdersHandlerInterface) {
	ctx := context.Background()

	return ctx, NewOrdersHandler(repo, redis, pub, log.DefaultLogger)
}

func TestOrderHandler_GetOrdersForUser(t *testing.T) {
	repo := mocks.NewOrdersRepository(t)
	redis := mocks.NewRedisHandlerInterface(t)
	pub := mocks.NewPublisherInterface(t)

	ctx, handler := setupOrdersHandler(repo, redis, pub)

	ordersResponse := []*entity.OrdersEntity{
		{
			ID:         "order-1",
			CustomerID: 10,
			CartID:     23,
			Status:     "SUCCESS",
			PaymentRef: "payment-1",
		},
		{
			ID:         "order-2",
			CustomerID: 10,
			CartID:     43,
			Status:     "FAILED",
			PaymentRef: "payment-2",
		},
	}

	type args struct {
		OrderID string
	}

	tests := []struct {
		title       string
		args        args
		want        []*v1.OrderData
		wantErr     bool
		errResponse error
		mock        func(repo *mocks.OrdersRepository)
	}{
		{
			title: "Success - GetOrdersForUser",
			args: args{
				OrderID: "order-1",
			},
			want: []*v1.OrderData{
				{
					Id:         "order-1",
					CustomerId: 10,
					CartId:     23,
					Status:     "SUCCESS",
					PaymentRef: "payment-1",
					CreatedAt:  "0001-01-01 00:00:00 +0000 UTC",
				},
				{
					Id:         "order-2",
					CustomerId: 10,
					CartId:     43,
					Status:     "FAILED",
					PaymentRef: "payment-2",
					CreatedAt:  "0001-01-01 00:00:00 +0000 UTC",
				},
			},
			wantErr:     false,
			errResponse: nil,
			mock: func(repo *mocks.OrdersRepository) {
				repo.On("GetOrdersByCustomerID", mock.Anything, mock.Anything).Return(ordersResponse, nil).Once()
			},
		},
		{
			title: "Failure - Error in GetOrdersByCustomerID",
			args: args{
				OrderID: "order-1",
			},
			want:        nil,
			wantErr:     true,
			errResponse: fmt.Errorf("something went wrong in GetOrdersByCustomerID"),
			mock: func(repo *mocks.OrdersRepository) {
				repo.On("GetOrdersByCustomerID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("something went wrong in GetOrdersByCustomerID")).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock(repo)

			response, err := handler.GetOrdersForUser(ctx, test.args.OrderID)

			if err == nil {
				assert.False(t, test.wantErr)
				assert.Equal(t, response, test.want)
			}

			if err != nil {
				if !test.wantErr {
					assert.Fail(t, fmt.Sprintf("Expected no errors in GetOrdersForUser, received: %v", err.Error()))
				}

				assert.Equal(t, err.Error(), test.errResponse.Error())
			}
		})
	}
}
