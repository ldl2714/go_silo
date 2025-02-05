package services

import (
	"fmt"
	"go_silo/modbus"
	"go_silo/models"
	"go_silo/precision_conversion"
	"go_silo/profiles"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPid 获取 pid 表中数据
func GetPid(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 设置一个空的查询条件 (相当于 find({}) )
		filter := bson.D{}
		// 获取集合
		collection := db.Collection("pid")

		// 查询数据
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println("Error finding pids:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx) //确保查询结束时关闭 cursor

		// 存储结果
		var pids []models.PidModel
		for cursor.Next(ctx) {
			// 解码单个 pid 到结构体
			var pid models.PidModel
			if err := cursor.Decode(&pid); err != nil {
				log.Println("Error decoding pid:", err)
				continue // 如果解码失败，跳过当前记录
			}
			pid.PID_SP = math.Round(pid.PID_SP*100) / 100
			if err := cursor.Decode(&pid); err != nil {
				log.Println("Error decoding pid:", err)
				continue // 如果解码失败，跳过当前记录
			}
			pids = append(pids, pid)
		}
		// 如果 cursor 在遍历时发生错误，打印错误并返回
		if err := cursor.Err(); err != nil {
			log.Println("Cursor error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 如果没有找到结果，返回提示信息
		if len(pids) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No static models found"})
			return
		}
		// 返回查询结果，200 状态码
		c.JSON(http.StatusOK, pids)
	}
}

type DiskPidPlcData struct {
	ID    string      `json:"id"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// PutPid 修改 pid 表中数据
func PutPid(client *modbus.ModbusClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data DiskPidPlcData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		// 查找配置文件
		var conf *profiles.ParameterProfile
		for _, profile := range profiles.ParameterProfiles {
			if profile.ID == data.ID {
				conf = &profile
				break
			}
		}
		if conf == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		field := strings.ToUpper(data.Field)
		address, err := getFieldAddress(conf, field)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pidFormatType := profiles.PidFormatTypeMap[field]

		switch pidFormatType {
		case profiles.Boolean:
			var value bool
			switch v := data.Value.(type) {
			case bool:
				value = v
			case float64:
				value = v != 0
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value type for boolean"})
				return
			}
			err = client.WriteCoil(uint16(address), value)
		case profiles.Number:
			value, ok := data.Value.(float64)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value type for number"})
				return
			}
			err = client.WriteRegisters(uint16(address), precision_conversion.Transform32FloatTo16BitSmall(float32(value)))
		case profiles.Time:
			value, ok := data.Value.(float64)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value type for time"})
				return
			}
			err = client.WriteRegisters(uint16(address), precision_conversion.TransToTimeWrite(uint32(value)))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown PID format type"})
			return
		}

		if err != nil {
			log.Printf("Failed to write to PLC: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to PLC"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

// getFieldAddress 获取配置文件中的字段地址
func getFieldAddress(conf *profiles.ParameterProfile, field string) (int, error) {
	r := reflect.ValueOf(conf)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return 0, fmt.Errorf("invalid field")
	}
	return int(f.Int()), nil
}
