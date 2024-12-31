package convertor

import "sync"

var (
	orderConvertor *OrderConvertor
	orderOnce      sync.Once
)

var (
	itemConvertor *ItemConvertor
	itemOnce      sync.Once
)

var (
	itemWithQuantityConvertor *ItemWithQuantityConvertor
	itemWithQuantityOnce      sync.Once
)

func NewOrderConvertor() *OrderConvertor {
	orderOnce.Do(func() {
		orderConvertor = new(OrderConvertor)
	})
	return orderConvertor
}

func NewItemConvertor() *ItemConvertor {
	itemOnce.Do(func() {
		itemConvertor = new(ItemConvertor)
	})
	return itemConvertor
}

func NewItemWithQuantityConvertor() *ItemWithQuantityConvertor {
	itemWithQuantityOnce.Do(func() {
		itemWithQuantityConvertor = new(ItemWithQuantityConvertor)
	})
	return itemWithQuantityConvertor
}
