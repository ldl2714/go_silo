package utils

import (
	"context"
	"fmt"
	"go_silo/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StartInsertingRealTime 每秒向 static 表中插入实时时间
func StartInsertingRealTime(db *mongo.Database) {
	go func() {
		collection := db.Collection("static")
		for {
			// 设置查询条件，查找 id 为 1 的文档
			filter := bson.D{{Key: "id", Value: 1}}
			// 设置当前时间
			update := bson.M{
				"$set": bson.M{
					"time": time.Now(),
				},
			}

			// 更新文档
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Printf("Error updating static: %v", err)
			} // else {
			// log.Printf("Updated static with current time")
			// }

			// 等待一秒钟
			time.Sleep(1 * time.Second)
		}
	}()
}

// GetShift 返回当前时间的班次（白班或夜班）
func GetShift() string {
	now := time.Now()
	hour := now.Hour()

	if hour >= 8 && hour < 17 {
		return "白班"
	}
	return "夜班"
}

// StartWebSocket 启动 websocket 任务
func StartWebSocket(db *mongo.Database) {
	// 设置 WebSocket 路由 	--- 实时趋势图
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(db, w, r)
	})
	// 设置 WebSocket 路由   --- 实时传递 belt 表中数据 给前端
	http.HandleFunc("/belts", func(w http.ResponseWriter, r *http.Request) {
		StartWebSocketBelt(db, w, r)
	})
	// 设置 WebSocket 路由   --- 实时传递 belt 表中数据 给前端
	http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		StartWebSocketAlert(db, w, r)
	})
	// 启动 Websocket 服务器，监听 8080 端口
	fmt.Println("WebSocket server is running on ws://localhost:8080/ws and ws://localhost:8080/belts and ws://localhost:8080/alerts")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe error:", err)
		}
	}()

}

// 升级器，用于将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

// handleWebSocket 每秒读取 belt 表中的数据并通过 WebSocket 发送给前端   --- 实时趋势图
func handleWebSocket(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var mu sync.Mutex
	done := make(chan struct{})
	var once sync.Once // 使用 sync.Once 确保只关闭一次

	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	go func() {
		defer closeDone() // 确保 goroutine 退出时关闭 done
		collectionBelt := db.Collection("belt")
		for {
			select {
			case <-done:
				return
			default:
				var belts []models.BeltModel_PV_SP
				cursor, err := collectionBelt.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{
					"id":      1,
					"specVol": 1,
					"vol":     1,
				}))
				if err != nil {
					log.Println("Error finding belt data:", err)
					return
				}
				if err = cursor.All(context.Background(), &belts); err != nil {
					log.Println("Error decoding belt data:", err)
					return
				}

				message := map[string]interface{}{
					"date":  time.Now().Format("2006-01-02 15:04:05"),
					"belts": belts,
				}

				mu.Lock()
				if err := conn.WriteJSON(message); err != nil {
					log.Println("Error writing JSON to WebSocket:", err)
					mu.Unlock()
					return
				}
				mu.Unlock()

				time.Sleep(1 * time.Second)
			}
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		closeDone() // WebSocket 关闭时调用
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			closeDone() // ReadMessage 出错时调用
			break
		}
	}
}

// 实时传递 belt 表中数据 给前端
func StartWebSocketBelt(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var mu sync.Mutex
	done := make(chan struct{})
	var once sync.Once // 使用 sync.Once 确保只关闭一次

	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	go func() {
		defer closeDone() // 确保 goroutine 退出时关闭 done
		collectionBelt := db.Collection("belt")
		for {
			select {
			case <-done:
				return
			default:
				var belts []models.BeltModel_Other
				cursor, err := collectionBelt.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{
					"id":           1,
					"parent":       1,
					"water":        1,
					"materialName": 1,
					"ratio":        1,
					"specVol":      1,
					"vol":          1,
					"diff":         1,
					"rate":         1,
					"wetAcc":       1,
					"dryAcc":       1,
				}))
				if err != nil {
					log.Println("Error finding belt data:", err)
					return
				}
				if err = cursor.All(context.Background(), &belts); err != nil {
					log.Println("Error decoding belt data:", err)
					return
				}

				message := map[string]interface{}{
					"belts": belts,
				}

				mu.Lock()
				if err := conn.WriteJSON(message); err != nil {
					log.Println("Error writing JSON to WebSocket:", err)
					mu.Unlock()
					return
				}
				mu.Unlock()

				time.Sleep(1 * time.Second)
			}
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		closeDone() // WebSocket 关闭时调用
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			closeDone() // ReadMessage 出错时调用
			break
		}
	}
}

// 报警信息
// StartAlarm 每秒向前端传递报警信息数据
func StartWebSocketAlert(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var mu sync.Mutex
	done := make(chan struct{})
	var once sync.Once // 使用 sync.Once 确保只关闭一次

	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	go func() {
		defer closeDone() // 确保 goroutine 退出时关闭 done
		collectionAlert := db.Collection("alert")
		for {
			select {
			case <-done:
				return
			default:
				var alerts []models.SendAlertModel
				cursor, err := collectionAlert.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{
					"id":        1,
					"diffAlert": 1,
				}))
				if err != nil {
					log.Println("Error finding alert data:", err)
					return
				}
				if err = cursor.All(context.Background(), &alerts); err != nil {
					log.Println("Error decoding alert data:", err)
					return
				}

				message := map[string]interface{}{
					"alerts": alerts,
				}

				mu.Lock()
				if err := conn.WriteJSON(message); err != nil {
					log.Println("Error writing JSON to WebSocket:", err)
					mu.Unlock()
					return
				}
				mu.Unlock()

				time.Sleep(1 * time.Second)
			}
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		closeDone() // WebSocket 关闭时调用
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			closeDone() // ReadMessage 出错时调用
			break
		}
	}
}
