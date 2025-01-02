package services

import (
	"go_silo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLevel(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		filter := bson.D{}
		collection := db.Collection("level")

		// 添加排序选项，按时间字段降序排序
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "date", Value: -1}})

		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var levels []models.LevelModel
		for cursor.Next(ctx) {
			var level models.LevelModel
			if err := cursor.Decode(&level); err != nil {
				continue
			}
			levels = append(levels, level)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(levels) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No level models found"})
			return
		}

		// 获取 level 表中数据
		c.JSON(http.StatusOK, levels)
	}
}
