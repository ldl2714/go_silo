package utils

import (
	"go_silo/modbus"
	"log"

	"github.com/robfig/cron/v3"
)

// 定义一个全局的 ModbusClient 实例
var globalClientPLC *modbus.ModbusClient

// StartSchedulers 启动所有定时任务
func StartSchedulers() {
	c := cron.New(cron.WithSeconds())

	// 初始化 realPLC 并连接到 PLC
	globalClientPLC = modbus.Modbus()
	if globalClientPLC == nil {
		log.Fatalf("Failed to create Modbus client")
	}

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

	// 每秒调用一次的定时任务
	_, err = c.AddFunc("*/1 * * * * *", func() {
		// 使用持久的 globalClientPLC 实例读取数据
		modbus.ReadRegisters(globalClientPLC)
	})
	if err != nil {
		log.Fatalf("Error scheduling every second task: %v", err)
	}

	c.Start()
}
