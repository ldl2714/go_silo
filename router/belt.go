package router

import (
	"go_silo/modbus"
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func BeltRouter(r *gin.Engine, db *mongo.Database, client *modbus.ModbusClient) {
	//获取皮带信息
	r.GET("/belts", services.GetBelt(db))
	//修改皮带配比
	r.PUT("/belts/ratio", services.UpdateBeltRatio(db, client))
	//修改皮带物料
	r.PUT("/belts/materialId", services.UpdateBeltMaterialId(db))
}
