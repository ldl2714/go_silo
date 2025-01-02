package main

import (
	"go_silo/db"
	"go_silo/modbus"
	"go_silo/router"
	"go_silo/utils"
	"log"
)

func main() {
	// 初始化modbus客户端
	modbusClient := modbus.Modbus()

	// 初始化数据库连接
	database := db.Mongo

	//每秒向Static表中插入一次实时时间
	utils.StartInsertingRealTime(database)

	// 启动所有定时任务
	utils.StartSchedulers(modbusClient, database)

	// 设置路由并传递数据库实例
	r := router.SetupRouter(modbusClient, database)
	// pkgB.CallA()

	// 启动服务器
	if err := r.Run(":2714"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
