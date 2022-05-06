# simple-orderbook
A simple orderbook implement by golang

## Introduce
use red-black tree to implement the features:
* add order
* remove order

We should resort orders when a new order created or removed.
The red-black tree can keep O(nlogn) to insert and remove a node.

So we can keep two red-black trees to match buy orders and sell orders.

Finally, we can get the result like this:
```
-------
2088.02 3
2087.6 6
2086.55 4
-------
2086.55
-------
2086 3
2085.01 5
2082.34 1
2081.11 7
```
## use example
```golang
	orders := make([]*orderbook.Order, 0)
	for _, arg := range testArgs {
		orders = append(orders, generateOrders(arg.price, arg.amount, arg.dir))
	}

	engine := orderbook.NewMatchEngine()
	for _, order := range orders {
		engine.ProcessOrder(context.Background(), order)
	}
```
## bench test result
```shell
goos: linux
goarch: amd64
pkg: github.com/fghpdf/simple-orderbook/orderbook
BenchmarkMatchEngine-16           992701              9322 ns/op            1563 B/op         61 allocs/op
BenchmarkMatchEngine-16          1000000              9103 ns/op            1582 B/op         62 allocs/op
BenchmarkMatchEngine-16          1000000              9379 ns/op            1548 B/op         60 allocs/op
BenchmarkMatchEngine-16          1000000              9879 ns/op            1545 B/op         60 allocs/op
PASS
ok      github.com/fghpdf/simple-orderbook/orderbook    38.065s
```