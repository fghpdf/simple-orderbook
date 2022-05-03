package orderbook

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	Nonce     string // order id
	Trader    string
	Direction DirectionEnum
	Side      SideEnum
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Slippage  uint
	CreateAt  time.Time
	UpdateAt  time.Time
	ExpireAt  time.Time
}

type OrderKey struct {
	Nonce    string // order id
	Price    decimal.Decimal
	CreateAt time.Time
}

type MatchResponse struct {
	TakerOrder   *Order
	MatchRecords []*MatchedRecord
}

type MatchedRecord struct {
	FilledPrice  decimal.Decimal
	FilledAmount decimal.Decimal
	TakerOrder   *Order
	MakerOrder   *Order
}
