package orderbook

import (
	"context"

	"github.com/shopspring/decimal"
)

func match(orders []*Order, direction DirectionEnum) []*Order {
	book := NewBook(direction)
	for _, order := range orders {
		book.Add(order)
	}

	var matched []*Order
	for book.getTop() != nil {
		order := book.getTop()
		matched = append(matched, order)
		book.Remove(order)
	}
	return matched
}

type MatchEngine struct {
	buyBook     *OrderBook
	sellBook    *OrderBook
	marketPrice decimal.Decimal
}

func NewMatchEngine() *MatchEngine {
	return &MatchEngine{
		buyBook:     NewBook(DirectionEnumBuy),
		sellBook:    NewBook(DirectionEnumSell),
		marketPrice: decimal.Zero,
	}
}

func (e *MatchEngine) ProcessOrder(ctx context.Context, order *Order) *MatchResponse {
	if order == nil {
		return nil
	}
	switch order.Direction {
	case DirectionEnumBuy:
		return e.processTaker(ctx, order, e.sellBook, e.buyBook)
	case DirectionEnumSell:
		return e.processTaker(ctx, order, e.buyBook, e.sellBook)
	}

	return nil
}

func (e *MatchEngine) processTaker(ctx context.Context, takerOrder *Order,
	makerBook *OrderBook, takerBook *OrderBook) *MatchResponse {
	res := &MatchResponse{
		TakerOrder:   takerOrder,
		MatchRecords: make([]*MatchedRecord, 0),
	}
	for true {
		makerOrder := makerBook.getTop()
		if makerOrder == nil {
			// Counterparty not exist
			break
		}

		if takerOrder.Direction == DirectionEnumBuy &&
			takerOrder.Price.LessThan(makerOrder.Price) {
			// Taker's buy price is less than first maker price
			// cannot fill
			break
		}

		if takerOrder.Direction == DirectionEnumSell &&
			takerOrder.Price.GreaterThan(makerOrder.Price) {
			// Taker's sell price is greater than first maker price
			// cannot fill
			break
		}

		// fill by maker price
		e.marketPrice = makerOrder.Price
		matchedAmount := decimal.Min(takerOrder.Amount, makerOrder.Amount)
		// fill log
		res.MatchRecords = append(res.MatchRecords, &MatchedRecord{
			FilledPrice:  makerOrder.Price,
			FilledAmount: matchedAmount,
			TakerOrder:   takerOrder,
			MakerOrder:   makerOrder,
		})
		// update amount
		takerOrder.Amount = takerOrder.Amount.Sub(matchedAmount)
		makerOrder.Amount = makerOrder.Amount.Sub(matchedAmount)
		// delete the order from makerOrder when filled completely
		if makerOrder.Amount.IsZero() {
			makerBook.Remove(makerOrder)
		}
		// break when taker order filled completely
		if takerOrder.Amount.IsZero() {
			break
		}
	}
	// put incomplete taker order into taker book
	if !takerOrder.Amount.IsZero() {
		takerBook.Add(takerOrder)
	}
	return res
}

func (e *MatchEngine) String() string {
	line := "\n-------\n"
	return line + e.sellBook.String() + line + e.marketPrice.String() + line + e.buyBook.String()
}
