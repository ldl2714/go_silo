package utils

import (
	"go_silo/modbus"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/robfig/cron/v3"
)

// StartSchedulers 启动所有定时任务
func StartSchedulers(client *modbus.ModbusClient, db *mongo.Database) {
	c := cron.New(cron.WithSeconds())

	// 每天的 8 点和 18 点调用 updateShiftService
	_, err := c.AddFunc("0 0 8,18 * * *", func() {
		err := GetShift()
		if err != "" {
			log.Printf("Error in updateShiftService: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling updateShiftService: %v", err)
	}

	// 每3秒调用一次的定时任务
	_, err = c.AddFunc("*/3 * * * * *", func() {
		// 使用持久的 ModbusClient 实例读取数据
		// modbus.ReadVol(client, db)
	})
	if err != nil {
		log.Fatalf("Error scheduling every second task: %v", err)
	}

	//每秒调用一次的定时任务
	_, err = c.AddFunc("* * * * * *", func() {
		// 使用持久的 ModbusClient 实例读取数据
		modbus.UpdateBeltTable(client, db)
	})
	if err != nil {
		log.Fatalf("Error scheduling every second task: %v", err)
	}
	c.Start()
}
