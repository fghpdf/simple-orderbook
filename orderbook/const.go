package orderbook

type DirectionEnum int

const (
	DirectionEnumBuy DirectionEnum = iota + 1
	DirectionEnumSell
)

type SideEnum int

const (
	SideEnumLong SideEnum = iota + 1
	SideEnumShort
)

