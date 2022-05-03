package orderbook

import (
	"strings"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type OrderBook struct {
	direction DirectionEnum
	book      *rbt.Tree
}

func NewBook(direction DirectionEnum) *OrderBook {
	compare := BuyComparator
	if direction == DirectionEnumSell {
		compare = SellComparator
	}
	return &OrderBook{
		direction: direction,
		book:      rbt.NewWith(compare),
	}
}

func (ob *OrderBook) Add(order *Order) {
	ob.book.Put(&OrderKey{
		Nonce:    order.Nonce,
		Price:    order.Price,
		CreateAt: order.CreateAt,
	}, order)
}

func (ob *OrderBook) Remove(order *Order) {
	ob.book.Remove(&OrderKey{
		Nonce:    order.Nonce,
		Price:    order.Price,
		CreateAt: order.CreateAt,
	})
}

func (ob *OrderBook) getTop() *Order {
	if ob.book.Size() == 0 {
		return nil
	}
	return ob.book.Left().Value.(*Order)
}

func SellComparator(a, b interface{}) int {
	aAsserted := a.(*OrderKey)
	bAsserted := b.(*OrderKey)
	if strings.EqualFold(aAsserted.Nonce, bAsserted.Nonce) {
		return 0
	}
	// low price to high price
	cmp := aAsserted.Price.Cmp(bAsserted.Price)
	// older order to newer order
	if cmp == 0 {
		if aAsserted.CreateAt.Before(bAsserted.CreateAt) {
			return -1
		}
		return 1
	}

	return cmp
}

func BuyComparator(a, b interface{}) int {
	aAsserted := a.(*OrderKey)
	bAsserted := b.(*OrderKey)
	if strings.EqualFold(aAsserted.Nonce, bAsserted.Nonce) {
		return 0
	}
	// high price to low price
	cmp := bAsserted.Price.Cmp(aAsserted.Price)
	// older order to newer order
	if cmp == 0 {
		if aAsserted.CreateAt.Before(bAsserted.CreateAt) {
			return -1
		}
		return 1
	}

	return cmp
}
