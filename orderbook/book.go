package orderbook

import (
	"fmt"
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

func (ob *OrderBook) String() string {
	if ob.book.Size() == 0 {
		return "(empty)"
	}
	var orders []string
	for _, value := range ob.book.Values() {
		order := value.(*Order)
		orders = append(orders, fmt.Sprintf("%s %s",
			order.Price.String(), order.Amount.String()))
	}
	if ob.direction == DirectionEnumSell {
		for i, j := 0, len(orders)-1; i < j; i, j = i+1, j-1 {
			orders[i], orders[j] = orders[j], orders[i]
		}
	}
	return strings.Join(orders, "\n")
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
