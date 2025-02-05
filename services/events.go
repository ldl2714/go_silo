package services

import (
	"go_silo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 获取 event 表中数据
func GetEvent(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		filter := bson.D{}
		collection := db.Collection("event")

		// 添加排序选项，按时间字段降序排序
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "date", Value: -1}})

		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var events []models.EventModel
		for cursor.Next(ctx) {
			var event models.EventModel
			if err := cursor.Decode(&event); err != nil {
				continue
			}
			events = append(events, event)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(events) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No event models found"})
			return
		}

		// 获取 event 表中数据
		c.JSON(http.StatusOK, events)
	}

}
