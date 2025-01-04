package adapters

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/Nicknamezz00/gorder-v2/common/config"
	"github.com/Nicknamezz00/gorder-v2/stock/entity"
	"github.com/Nicknamezz00/gorder-v2/stock/infrastructure/persistent"
	"github.com/Nicknamezz00/gorder-v2/stock/infrastructure/persistent/builder"
	"github.com/spf13/viper"
	gormlogger "gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *persistent.MySQL {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		"",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	testDB := viper.GetString("mysql.dbname") + "_shadow"
	assert.NoError(t, db.Exec("DROP DATABASE IF EXISTS "+testDB).Error)
	assert.NoError(t, db.Exec("CREATE DATABASE IF NOT EXISTS "+testDB).Error)
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		testDB,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&persistent.StockModel{}))

	return persistent.NewMySQLWithDB(db)
}

func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := setupTestDB(t)

	// 准备初始数据
	var (
		testItem           = "item-1"
		initialStock int32 = 100
	)
	err := db.Create(ctx, nil, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  initialStock,
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	concurrentGoroutines := 10
	for i := 0; i < concurrentGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					// 模拟减少库存
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res cannot be empty")

	expectedStock := initialStock - int32(concurrentGoroutines)
	assert.Equal(t, expectedStock, res[0].Quantity)
}

func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := setupTestDB(t)

	// 准备初始数据
	var (
		testItem           = "item-1"
		initialStock int32 = 5
	)
	err := db.Create(ctx, nil, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  initialStock,
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	concurrentGoroutines := 100
	for i := 0; i < concurrentGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					// 模拟减少库存
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
		time.Sleep(20 * time.Millisecond)
	}

	wg.Wait()
	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res cannot be empty")

	assert.GreaterOrEqual(t, res[0].Quantity, int32(0))
}
