package handler

import (
	"cart-api/proto/cartapi"
	"context"

	"github.com/lenny-mo/cart/proto/cart"
)

// CartApi 实现下面的接口
// Server API for CartApi service
//
//	type CartApiHandler interface {
//		FindAll(context.Context, *FindAllRequest, *FindAllResponse) error
//		Add(context.Context, *AddCartRequest, *AddCartResponse) error
//		Update(context.Context, *UpdateRequest, *UpdateResponse) error
//		Delete(context.Context, *DeleteRequest, *DeleteResponse) error
//	}
type CartAPI struct {
	CarService cart.CartService
}

// FindAll 获取用户的购物车列表
func (c *CartAPI) FindAll(ctx context.Context, req *cartapi.FindAllRequest, res *cartapi.FindAllResponse) error {

	return nil
}

func (c *CartAPI) Add(context.Context, *cartapi.AddCartRequest, *cartapi.AddCartResponse) error {
	return nil
}

func (c *CartAPI) Update(ctx context.Context, req *cartapi.UpdateRequest, res *cartapi.UpdateResponse) error {
	return nil
}

func (c *CartAPI) Delete(ctx context.Context, req *cartapi.DeleteRequest, res *cartapi.DeleteResponse) error {
	return nil
}
