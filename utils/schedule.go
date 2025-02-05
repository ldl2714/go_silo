package utils

// import (
// 	"go_silo/modbus"
// 	"go_silo/tasks"
// 	"log"

// 	"go.mongodb.org/mongo-driver/mongo"

// 	"github.com/robfig/cron/v3"
// )

// StartSchedulers 启动所有定时任务
import (
	"go_silo/modbus"
	"go_silo/tasks"
	"log"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/robfig/cron/v3"
)

// StartSchedulers 启动所有定时任务
func StartSchedulers(client *modbus.ModbusClient, db *mongo.Database, redisClient *redis.Client) {
	c := cron.New(cron.WithSeconds())

	// 每天的 8 点和 18 点调用 updateShiftService
	_, err := c.AddFunc("0 59 7,17 * * *", func() {
		// 调用 GetShift
		err := GetShift()
		if err != "" {
			log.Printf("Error in GetShift: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling tasks: %v", err)
	}

	// 每天的 8 点和 18 点调用 ProductionTable
	_, err = c.AddFunc("0 0 8,18 * * *", func() {
		// _, err = c.AddFunc("0 9 11 * * *", func() {
		// 调用 ProductionTable
		tasks.ProductionTable(db)
	})

	if err != nil {
		log.Fatalf("Error scheduling tasks: %v", err)
	}

	// 每3秒调用一次的定时任务
	_, err = c.AddFunc("*/3 * * * * *", func() {
		// 使用持久的 ModbusClient 实例读取数据
		modbus.ReadVol(client, db)
		modbus.MaterialLevel(client, db)
		modbus.ReadPid(client, db)
	})
	if err != nil {
		log.Fatalf("Error scheduling every second task: %v", err)
	}

	//每秒调用一次的定时任务
	_, err = c.AddFunc("* * * * * *", func() {
		// 使用持久的 ModbusClient 实例读取数据
		tasks.UpdateBeltVol(db)
		tasks.ProductionAryRunTime(db)
		tasks.PostNewDateByTrend(db, redisClient)
		tasks.AlertTask(db)
	})

	if err != nil {
		log.Fatalf("Error scheduling every second task: %v", err)
	}
	// 这个表达式代表每天的 00:02 执行任务
	_, err = c.AddFunc("0 2 0 * * *", func() {
		tasks.SaveHistoryVol(db, redisClient)
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}
	c.Start()
}
