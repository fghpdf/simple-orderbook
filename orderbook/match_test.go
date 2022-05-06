package orderbook

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	orders := []*Order{
		{
			Nonce:    "3400-NEW",
			Price:    decimal.NewFromInt(3400),
			Amount:   decimal.NewFromInt(1),
			CreateAt: time.Unix(1650612313, 0),
		},
		{
			Nonce:    "3300-NEW",
			Price:    decimal.NewFromInt(3300),
			Amount:   decimal.NewFromInt(2),
			CreateAt: time.Unix(1650612313, 0),
		},
		{
			Nonce:    "3300-OLD",
			Price:    decimal.NewFromInt(3300),
			Amount:   decimal.NewFromInt(2),
			CreateAt: time.Unix(1650612300, 0),
		},
		{
			Nonce:    "3200-OLD",
			Price:    decimal.NewFromInt(3200),
			Amount:   decimal.NewFromInt(2),
			CreateAt: time.Unix(1650612300, 0),
		},
	}

	type args struct {
		orders    []*Order
		direction DirectionEnum
	}
	tests := []struct {
		name string
		args args
		want []*Order
	}{
		{
			"Buy 4 orders",
			args{
				orders,
				DirectionEnumBuy,
			},
			[]*Order{
				orders[0],
				orders[2],
				orders[1],
				orders[3],
			},
		},
		{
			"Sell 4 orders",
			args{
				orders,
				DirectionEnumSell,
			},
			[]*Order{
				orders[3],
				orders[2],
				orders[1],
				orders[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := match(tt.args.orders, tt.args.direction)
			assert.Equal(t, tt.want, got)
		})
	}
}

const (
	sell = DirectionEnumSell
	buy  = DirectionEnumBuy
)

func TestMatchEngine(t *testing.T) {
	type argus struct {
		dir    DirectionEnum
		price  float64
		amount int64
	}
	testArgs := []*argus{
		{buy, 2082.34, 1},
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
	orders := make([]*Order, 0)
	for _, arg := range testArgs {
		orders = append(orders, generateTestOrders(arg.price, arg.amount, arg.dir))
	}
	engine := NewMatchEngine()
	for _, order := range orders {
		engine.ProcessOrder(context.Background(), order)
	}

	assert.Equal(t, "2086.55", engine.GetPrice().String())
	assert.Equal(t, "2086", engine.GetBuyBook().getTop().Price.String())
	assert.Equal(t, "3", engine.GetBuyBook().getTop().Amount.String())
	assert.Equal(t, "2086.55", engine.GetSellBook().getTop().Price.String())
	assert.Equal(t, "4", engine.GetSellBook().getTop().Amount.String())
}

func BenchmarkMatchEngine(b *testing.B) {
	b.ResetTimer()
	engine := NewMatchEngine()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		engine.ProcessOrder(ctx, generateTestOrders(rand.Float64()*10+2800, rand.Int63()*100,
			DirectionEnum(rand.Intn(1)+1)))
	}
}

func generateTestOrders(price float64, amount int64, dir DirectionEnum) *Order {
	return &Order{
		Nonce:     uuid.New().String(),
		Price:     decimal.NewFromFloat(price),
		Amount:    decimal.NewFromInt(amount),
		Direction: dir,
		CreateAt:  time.Now(),
	}
}
