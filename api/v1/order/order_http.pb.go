// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v4.25.3
// source: v1/order/order.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationOrderCreateOrder = "/order.v1.Order/CreateOrder"
const OperationOrderGetAllOrders = "/order.v1.Order/GetAllOrders"

type OrderHTTPServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderReply, error)
	GetAllOrders(context.Context, *GetAllOrdersForUserRequest) (*GetAllOrdersForUserReply, error)
}

func RegisterOrderHTTPServer(s *http.Server, srv OrderHTTPServer) {
	r := s.Route("/")
	r.GET("/order/{user_id}", _Order_GetAllOrders0_HTTP_Handler(srv))
	r.POST("/order", _Order_CreateOrder0_HTTP_Handler(srv))
}

func _Order_GetAllOrders0_HTTP_Handler(srv OrderHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetAllOrdersForUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOrderGetAllOrders)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAllOrders(ctx, req.(*GetAllOrdersForUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAllOrdersForUserReply)
		return ctx.Result(200, reply)
	}
}

func _Order_CreateOrder0_HTTP_Handler(srv OrderHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateOrderRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationOrderCreateOrder)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateOrder(ctx, req.(*CreateOrderRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateOrderReply)
		return ctx.Result(200, reply)
	}
}

type OrderHTTPClient interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest, opts ...http.CallOption) (rsp *CreateOrderReply, err error)
	GetAllOrders(ctx context.Context, req *GetAllOrdersForUserRequest, opts ...http.CallOption) (rsp *GetAllOrdersForUserReply, err error)
}

type OrderHTTPClientImpl struct {
	cc *http.Client
}

func NewOrderHTTPClient(client *http.Client) OrderHTTPClient {
	return &OrderHTTPClientImpl{client}
}

func (c *OrderHTTPClientImpl) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...http.CallOption) (*CreateOrderReply, error) {
	var out CreateOrderReply
	pattern := "/order"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationOrderCreateOrder))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *OrderHTTPClientImpl) GetAllOrders(ctx context.Context, in *GetAllOrdersForUserRequest, opts ...http.CallOption) (*GetAllOrdersForUserReply, error) {
	var out GetAllOrdersForUserReply
	pattern := "/order/{user_id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationOrderGetAllOrders))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
