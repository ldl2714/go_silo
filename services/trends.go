package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetTrend 获取历史曲线
func GetTrend(db *mongo.Database, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		dateDay := c.Query("dateDay")

		cacheVolKey := dateDay + ":" + id + ":vol"
		cacheSpecVolKey := dateDay + ":" + id + ":specVol"

		cachedVolValue, err := redisClient.Get(c.Request.Context(), cacheVolKey).Result()
		if err == nil && cachedVolValue != "" {
			cacheSpecVolValue, _ := redisClient.Get(c.Request.Context(), cacheSpecVolKey).Result()
			c.JSON(200, gin.H{
				"source": "cache",
				"data": []gin.H{
					{
						"date":    dateDay,
						"id":      id,
						"vol":     cachedVolValue,
						"specVol": cacheSpecVolValue,
					},
				},
			})
			return
		}

		collection := db.Collection("trend")
		filter := bson.M{"date": dateDay, "id": id}
		cursor, err := collection.Find(c, filter)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to query database"})
			return
		}
		defer cursor.Close(c)

		var results []bson.M
		if err = cursor.All(c, &results); err != nil {
			c.JSON(500, gin.H{"error": "failed to parse database results"})
			return
		}

		c.JSON(200, gin.H{
			"source": "database",
			"data":   results,
		})
	}
}
