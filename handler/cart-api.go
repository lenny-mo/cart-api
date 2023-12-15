package handler

import (
	"context"
	"fmt"

	"github.com/lenny-mo/cart-api/proto/cartapi"

	"github.com/lenny-mo/cart/proto/cart"
)

// CartApi 实现下面的接口
// Server API for CartApi service∂
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
	cartRes, err := c.CarService.FindAll(ctx, &cart.FindAllCartRequest{
		UserId: req.Userid,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 复制全部item
	for _, v := range cartRes.CartItems {
		res.CartItems = append(res.CartItems, &cartapi.CartItem{
			Skuid:    v.Skuid,
			Quantity: v.Quantity,
			Time:     v.Time,
			Status:   cartapi.CartStatus(v.Status),
		})
	}
	return nil
}

func (c *CartAPI) Add(ctx context.Context, req *cartapi.AddCartRequest, res *cartapi.AddCartResponse) error {
	cartRes, err := c.CarService.Add(ctx, &cart.AddCartRequest{
		UserId: req.UserId,
		Item: &cart.CartItem{
			Skuid:    req.Item.Skuid,
			Quantity: req.Item.Quantity,
			Time:     req.Item.Time,
			Status:   cart.CartStatus(req.Item.Status),
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	res.Msg = cartRes.Msg
	return nil
}

func (c *CartAPI) Update(ctx context.Context, req *cartapi.UpdateRequest, res *cartapi.UpdateResponse) error {
	cartRes, err := c.CarService.Update(ctx, &cart.UpdateRequest{
		UserId: req.UserId,
		Item: &cart.CartItem{
			Skuid:    req.Item.Skuid,
			Quantity: req.Item.Quantity,
			Time:     req.Item.Time,
			Status:   cart.CartStatus(req.Item.Status),
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	res.Msg = cartRes.Code
	return nil
}

func (c *CartAPI) Delete(ctx context.Context, req *cartapi.DeleteRequest, res *cartapi.DeleteResponse) error {
	cartRes, err := c.CarService.Delete(ctx, &cart.DeleteRequest{
		Userid: req.Userid,
		Skuid:  req.Skuid,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	res.Msg = cartRes.Msg
	return nil
}
