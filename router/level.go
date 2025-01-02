package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func LevelRouter(r *gin.Engine, db *mongo.Database) {
	// 获取 level 表中数据
	r.GET("/levels", services.GetLevel(db))
}
