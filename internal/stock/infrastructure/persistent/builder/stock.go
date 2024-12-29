package builder

import (
	"github.com/Nicknamezz00/gorder-v2/common/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stock struct {
	ID        []int64  `json:"ID,omitempty"`
	ProductID []string `json:"product_id,omitempty"`
	Quantity  []int32  `json:"quantity,omitempty"`
	Version   []int64  `json:"version,omitempty"`

	// extend fields
	OrderBy       string `json:"order_by,omitempty"`
	ForUpdateLock bool   `json:"for_update,omitempty"`
}

func NewStock() *Stock {
	return &Stock{}
}

func (s *Stock) FormatArg() (string, error) {
	return util.MarshalString(s)
}

func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
	db = s.fillWhere(db)
	if s.OrderBy != "" {
		db = db.Order(s.Order)
	}
	return db
}

func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
	if len(s.ID) > 0 {
		db = db.Where("ID in (?)", s.ID)
	}
	if len(s.ProductID) > 0 {
		db = db.Where("product_id in (?)", s.ProductID)
	}
	if len(s.Version) > 0 {
		db = db.Where("Version in (?)", s.Version)
	}
	if len(s.Quantity) > 0 {
		db = s.fillQuantityGT(db)
	}
	if s.ForUpdateLock {
		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
	return db
}

func (s *Stock) fillQuantityGT(db *gorm.DB) *gorm.DB {
	db = db.Where("Quantity >= ?", s.Quantity)
	return db
}

func (s *Stock) IDs(v ...int64) *Stock {
	s.ID = v
	return s
}

func (s *Stock) ProductIDs(v ...string) *Stock {
	s.ProductID = v
	return s
}

func (s *Stock) Order(v string) *Stock {
	s.OrderBy = v
	return s
}

func (s *Stock) Versions(v ...int64) *Stock {
	s.Version = v
	return s
}

func (s *Stock) QuantityGT(v ...int32) *Stock {
	s.Quantity = v
	return s
}

func (s *Stock) ForUpdate() *Stock {
	s.ForUpdateLock = true
	return s
}
