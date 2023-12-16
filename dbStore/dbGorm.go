package dbstore

import (
	"context"
	"fmt"
	_ "myShopEcommerce/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"


	"gorm.io/driver/postgres"
)

type Sqllogger struct {
	logger.Interface
}

func (l Sqllogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n===================================================]\n", sql)
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

// func Db_Init_mgo_v3_myshpoing() *mongo.Database {

// 	db, errDB := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017")) // ใช้ได้

// 	if errDB != nil {
// 		Error_message := fmt.Sprintf("Server Error : %v", errDB)
// 		panic(Error_message)
// 	}

// 	db_Database := db.Database("myshop")
// 	fmt.Println("Start MongoDB V3 .....")

// 	return db_Database
// }

func Postgres_init() *gorm.DB {
	
	dsn := "host=localhost user=postgres password=P@ssw0rd dbname=myshoping port=5432 sslmode=disable TimeZone=Asia/Bangkok"

	dial := postgres.Open(dsn)

	db_Postgr, err := gorm.Open(dial, &gorm.Config{
		Logger: &Sqllogger{},
		// DryRun: true,
	})

	if err != nil {
		fmt.Println("Postgre Error: ", err)
		panic("Postgr Can not connect")
	}

	return db_Postgr
}
