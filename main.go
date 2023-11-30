package main

import (
	_ "database/sql"
	"fmt"
	"strings"

	dbstore "myShopEcommerce/dbStore"
	"myShopEcommerce/handler"
	"myShopEcommerce/repository"
	ginpath "myShopEcommerce/router/gin_path"
	"myShopEcommerce/service"
	_ "myShopEcommerce/util"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	_ "github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")                               // อ้างการเข้าถึงข้อมูลด้วย . เช่น kafak.server
	viper.AutomaticEnv()                                   // เจอค่าใน config ก่อนจะเอาค่ามาใช้ ถ้าไม่เจอจะเอาใน yaml มาใช้
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // ใน shell ใช้ . ไม่ได้

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

//producer
func main() {
	
	ss := uuid.NewString()
	fmt.Println(ss)

	// ************************************** DB ****************************
	sqlClient := dbstore.DbGorm()

	// ********  Redis *******
	redisClient := initRedis()
	_ = redisClient

	// ***********  MongoDB ***********
	monGoClient := dbstore.Db_Init_mgo_v3_myshpoing()
	_ = monGoClient

	postgresClient := dbstore.Postgres_init()
	_ = postgresClient

	//******* Kafka Producer ********
	// producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	// if err != nil {
	// 	panic(err)
	// }

	// defer producer.Close()

	// *******  User ใช้ Kafka ********
	// ***************** User Handler *****************
	// kffka
	// handlerUser_Kafka := repository.NewUser_Repo_Kafka(producer)
	// handlerService := service.NewUser_Service_Kafka(handlerUser_Kafka)
	// NewUser_Handler_Kafka อยู่ใน handler ปกติ *******
	// handlerHandler_Userk := handler.NewUser_Handler_Kafka(handlerService)

	// **********************************   Handler  *********************************************
	// ***********  User Handler ****************
	handlerRepoUser := repository.NewUser_Repository(sqlClient)
	handlerService := service.NewUser_Service(handlerRepoUser)
	handlerHandler_User := handler.NewUser_Handler(handlerService)

	// *************  Store Handler  **************
	handlerRepo_Store := repository.NewStore_Repository(sqlClient)
	handlerService_Store := service.NewStore_Service(handlerRepo_Store)
	handlerHandler_Store := handler.NewStore_Handler(handlerService_Store)

	// *************** product ************
	// product Fuc New ต่างๆจะส่ง pointer เข้ามา มันต่างกันยังไง
	// REdis อยู่ layer นี้  - Getproduct กับ Search product
	handlerRepo_Product := repository.NewProduct_Repository(sqlClient)
	// handlerService_Product := service.NewProduct_Service(&handlerRepo_Product)
	// handlerHandler_Product := handler.NewProduct_Handler(&handlerService_Product)

	// *********** Redis **********
	// handlerRepo_Product := repository.NewProduct_Repository_redis(db)
	handlerService_Product := service.NewProduct_Service_Redis(&handlerRepo_Product, redisClient)
	handlerHandler_Product := handler.NewProduct_Handler_redis(&handlerService_Product)

	// ************* Cart Handler  ***********
	handlerRepo_Cart := repository.NewCart_Repo(sqlClient, monGoClient, postgresClient)
	handlerService_Cart := service.NewCart_Service(handlerRepo_Cart)
	handlerHandler_Cart := handler.NewCart_Handler(handlerService_Cart)

	//*********** Gin router *************
	ginpath.GinRouter_User(handlerHandler_User, handlerHandler_Store, handlerHandler_Cart, handlerHandler_Product)

}

func initRedis() *redis.Client {
	// addr := os.Getenv("REDIS_ADDR")
	addr := viper.GetString("redis.redis_add")
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})

}
