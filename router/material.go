package router

import (
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func MaterialRouter(r *gin.Engine, db *mongo.Database) {
	// 获取 material 表中数据
	r.GET("/materials", services.GetMaterial(db))
	// 修改 material 表中数据
	r.PUT("/materials", services.UpdateMaterial(db))
	// 新增 material 表中数据
	r.POST("/materials", services.AddMaterial(db))
	// 删除 material 表中数据
	r.DELETE("/materials", services.DeleteMaterial(db))
}
