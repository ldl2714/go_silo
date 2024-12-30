package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func EventRouter(r *gin.Engine, db *mongo.Database) {
	// 获取 event 表中数据
	r.GET("/events", services.GetEvent(db))
}
