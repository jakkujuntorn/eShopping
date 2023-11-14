package dbstore

import (
	"context"
	"fmt"
	_ "myShopEcommerce/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sqllogger struct {
	logger.Interface
}

func (l Sqllogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
sql,_ := fc()
fmt.Printf("%v\n===================================================]\n",sql)
}

func DbGorm() *gorm.DB {

	dsn := "root:P@ssw0rd@tcp(127.0.0.1:3306)/myshop?parseTime=true&loc=Local"
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: &Sqllogger{},
		// DryRun: true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		
		
	})
	if err != nil {
		panic("Db Can not connect")
	}

	// db.AutoMigrate(&models.UserRequest{})
	// db.AutoMigrate(models.Product{})
	// db.AutoMigrate(models.Picture_Product{})
	// db.AutoMigrate(models.Cart{})
	// db.AutoMigrate(models.CartItem{})

	return db
}
