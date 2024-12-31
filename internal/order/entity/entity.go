package entity

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}
