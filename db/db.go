package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo = InitMongo()

func InitMongo() *mongo.Database {
	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接到 MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接 MongoDB 错误: %v", err)
	}
	fmt.Println("数据库连接成功")
	return client.Database("bx_silo")
}

// package db

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var Mongo = initMongo()

// func initMongo() *mongo.Database {

// 	// 设置超时
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

//		// 连接到 MongoDB
//		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
//		if err != nil {
//			log.Println("连接 MongoDB 错误", err)
//			return nil
//		}
//		fmt.Println("数据库连接成功")
//		return client.Database("bx_silo")
//	}
