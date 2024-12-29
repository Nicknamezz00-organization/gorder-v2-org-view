package builder

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stock struct {
	id        []int64
	productID []string
	quantity  []int32
	version   []int64

	// extend fields
	order     string
	forUpdate bool
}

func NewStock() *Stock {
	return &Stock{}
}

func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
	db = s.fillWhere(db)
	if s.order != "" {
		db = db.Order(s.order)
	}
	return db
}

func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
	if len(s.id) > 0 {
		db = db.Where("id in (?)", s.id)
	}
	if len(s.productID) > 0 {
		db = db.Where("product_id in (?)", s.productID)
	}
	if len(s.version) > 0 {
		db = db.Where("version in (?)", s.version)
	}
	if len(s.quantity) > 0 {
		db = s.fillQuantityGT(db)
	}

	if s.forUpdate {
		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
	return db
}

func (s *Stock) fillQuantityGT(db *gorm.DB) *gorm.DB {
	db = db.Where("quantity >= ?", s.quantity)
	return db
}

func (s *Stock) IDs(v ...int64) *Stock {
	s.id = v
	return s
}

func (s *Stock) ProductIDs(v ...string) *Stock {
	s.productID = v
	return s
}

func (s *Stock) Order(v string) *Stock {
	s.order = v
	return s
}

func (s *Stock) Versions(v ...int64) *Stock {
	s.version = v
	return s
}

func (s *Stock) QuantityGT(v ...int32) *Stock {
	s.quantity = v
	return s
}

func (s *Stock) ForUpdate() *Stock {
	s.forUpdate = true
	return s
}
