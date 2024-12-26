package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StartInsertingRealTime 每秒向 static 表中插入实时时间
func StartInsertingRealTime(db *mongo.Database) {
	go func() {
		collection := db.Collection("static")
		for {
			// 设置查询条件，查找 id 为 1 的文档
			filter := bson.D{{Key: "id", Value: 1}}
			// 设置当前时间
			update := bson.M{
				"$set": bson.M{
					"time": time.Now(),
				},
			}

			// 更新文档
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Printf("Error updating static: %v", err)
			} // else {
			// log.Printf("Updated static with current time")
			// }

			// 等待一秒钟
			time.Sleep(1 * time.Second)
		}
	}()
}
