package orderbook

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_match(t *testing.T) {
	orders := []*Order{
		{
			Nonce:      "3400-NEW",
			Price:      decimal.NewFromInt(3400),
			Amount: decimal.NewFromInt(1),
			CreateAt:   time.Unix(1650612313, 0),
		},
		{
			Nonce:      "3300-NEW",
			Price:      decimal.NewFromInt(3300),
			Amount: decimal.NewFromInt(2),
			CreateAt:   time.Unix(1650612313, 0),
		},
		{
			Nonce:      "3300-OLD",
			Price:      decimal.NewFromInt(3300),
			Amount: decimal.NewFromInt(2),
			CreateAt:   time.Unix(1650612300, 0),
		},
		{
			Nonce:      "3200-OLD",
			Price:      decimal.NewFromInt(3200),
			Amount: decimal.NewFromInt(2),
			CreateAt:   time.Unix(1650612300, 0),
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
