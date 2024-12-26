package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func StaticRouter(r *gin.Engine, db *mongo.Database) {
	// 获取 static 表中数据
	r.GET("/statics", services.GetStatic(db))
	// 修改 static 表中数据
	r.PUT("/statics", services.UpdateStatic(db))
}
