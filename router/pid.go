package router

import (
	"go_silo/modbus"
	"go_silo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PidRouter(r *gin.Engine, db *mongo.Database, client *modbus.ModbusClient) {
	// 获取 pid 表中数据
	r.GET("/pids", services.GetPid(db))
	// 修改 pid 表中数据
	r.PUT("/pids", services.PutPid(client))
}
