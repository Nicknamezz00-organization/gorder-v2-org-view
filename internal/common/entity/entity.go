package entity

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
}

func (it Item) validate() error {
	//if err := util.AssertNotEmpty(it.ID, it.PriceID, it.Name); err != nil {
	//	return err
	//}
	var invalidFields []string
	if it.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	if it.Name == "" {
		invalidFields = append(invalidFields, "Name")
	}
	if it.PriceID == "" {
		invalidFields = append(invalidFields, "PriceID")
	}
	return fmt.Errorf("item=%v invalid, empty fields=[%s]", it, strings.Join(invalidFields, ","))
}

func NewItem(ID string, name string, quantity int32, priceID string) *Item {
	return &Item{ID: ID, Name: name, Quantity: quantity, PriceID: priceID}
}

func NewValidItem(ID string, name string, quantity int32, priceID string) (*Item, error) {
	item := NewItem(ID, name, quantity, priceID)
	if err := item.validate(); err != nil {
		return nil, err
	}
	return item, nil
}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}

func (iq ItemWithQuantity) validate() error {
	//if err := util.AssertNotEmpty(it.ID, it.PriceID, it.Name); err != nil {
	//	return err
	//}
	var invalidFields []string
	if iq.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	return errors.New(strings.Join(invalidFields, ","))
}

func NewItemWithQuantity(ID string, quantity int32) *ItemWithQuantity {
	return &ItemWithQuantity{ID: ID, Quantity: quantity}
}

func NewValidItemWithQuantity(ID string, quantity int32) (*ItemWithQuantity, error) {
	iq := NewItemWithQuantity(ID, quantity)
	if err := iq.validate(); err != nil {
		return nil, err
	}
	return iq, nil
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*Item
}

func NewValidOrder(ID string, customerID string, status string, paymentLink string, items []*Item) (*Order, error) {
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}
	return NewOrder(ID, customerID, status, paymentLink, items), nil
}
func NewOrder(ID string, customerID string, status string, paymentLink string, items []*Item) *Order {
	return &Order{ID: ID, CustomerID: customerID, Status: status, PaymentLink: paymentLink, Items: items}
}
