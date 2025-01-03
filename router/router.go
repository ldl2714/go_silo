package router

import (
	"go_silo/modbus"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(client *modbus.ModbusClient, db *mongo.Database) *gin.Engine {
	r := gin.Default()
	// Belt-皮带
	BeltRouter(r, db, client)
	// Material-物料
	MaterialRouter(r, db)
	// Static-静态信息
	StaticRouter(r, db)
	// Event-事件
	EventRouter(r, db)
	//Level-料位
	LevelRouter(r, db)
	// Pid
	PidRouter(r, db, client)
	// disk
	DiskRouter(r, db)
	return r
}
