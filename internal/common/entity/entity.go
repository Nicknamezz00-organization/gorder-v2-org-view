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

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*Item
}
