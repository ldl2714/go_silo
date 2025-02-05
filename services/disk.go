package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetDisk 获取 disk 表中数据
func GetDisk(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 disk 表中数据
		var results []struct {
			ID    int32 `json:"id" bson:"id"`
			Disk1 int32 `json:"disk1" bson:"disk1"`
			Disk2 int32 `json:"disk2" bson:"disk2"`
			Disk3 int32 `json:"disk3" bson:"disk3"`
			Disk4 int32 `json:"disk4" bson:"disk4"`
		}

		cursor, err := db.Collection("level").Find(c, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c)

		if err = cursor.All(c, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}

}

// PutDisk 更新 disk 表中数据
func PutDisk(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input map[string]interface{}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, ok := input["id"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		filter := bson.M{"id": int32(id)}
		update := bson.M{"$set": input}

		result, err := db.Collection("level").UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
	}
}
