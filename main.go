package main

import (
	"context"
	"log"
	"time"

	"github.com/fghpdf/simple-orderbook/orderbook"
	"github.com/shopspring/decimal"
)

const (
	sell = orderbook.DirectionEnumSell
	buy  = orderbook.DirectionEnumBuy
)

// 方向 价格 数量
type testArg struct {
	dir    orderbook.DirectionEnum
	price  float64
	amount int64
}

func main() {
	testArgs := []*testArg{
		{buy, 2082.34, 1},
		{sell, 2087.6, 2},
		{buy, 2087.8, 1},
		{sell, 2087.6, 2},
		{buy, 2087.8, 1},
		{buy, 2085.01, 5},
		{sell, 2088.02, 3},
		{sell, 2087.60, 6},
		{buy, 2081.11, 7},
		{buy, 2086.0, 3},
		{buy, 2088.33, 1},
		{sell, 2086.54, 2},
		{sell, 2086.55, 5},
		{buy, 2086.55, 3},
	}

	orders := make([]*orderbook.Order, 0)
	for _, arg := range testArgs {
		orders = append(orders, generateOrders(arg.price, arg.amount, arg.dir))
	}

	engine := orderbook.NewMatchEngine()
	for _, order := range orders {
		res := engine.ProcessOrder(context.Background(), order)
		log.Println(res.TakerOrder.Direction, res.TakerOrder.Price, res.TakerOrder.Amount)
	}
}

func generateOrders(price float64, amount int64, dir orderbook.DirectionEnum) *orderbook.Order {
	return &orderbook.Order{
		Price:     decimal.NewFromFloat(price),
		Amount:    decimal.NewFromInt(amount),
		Direction: dir,
		CreateAt:  time.Now(),
	}
}
