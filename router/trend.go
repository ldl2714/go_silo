package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func TrendRouter(r *gin.Engine, db *mongo.Database, redisClient *redis.Client) {
	//获取 历史曲线
	r.GET("/trends", services.GetTrend(db, redisClient))
}
