package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() *MySQL {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panicf("connect to mysql failed, err=%v", err)
	}
	return &MySQL{db: db}
}

type StockModel struct {
	ID        int64     `gorm:"column:id"`
	ProductID string    `gorm:"column:product_id"`
	Quantity  int32     `gorm:"column:quantity"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdateAt  time.Time `gorm:"column:updated_at"`
}

func (StockModel) TableName() string {
	return "o_stock"
}

func (d MySQL) BatchGetStockByID(ctx context.Context, productIDs []string) ([]StockModel, error) {
	var result []StockModel
	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}
