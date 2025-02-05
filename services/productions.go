package services

import (
	"go_silo/models"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetProduction 获取 production 表中数据
func GetProduction(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Query("date")
		shift := c.Query("shift")
		team := c.Query("team")
		log.Println(date, shift, team)
		filter := bson.M{
			"date":  date,
			"shift": shift,
			"team":  team,
		}

		// 查询数据
		cursor, err := db.Collection("production").Find(c, filter)
		if err != nil {
			log.Println("Error finding production:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c) // 确保查询结束时关闭 cursor

		// 存储结果
		var productions []models.ProductionModel
		for cursor.Next(c) {
			// 解码单个 production 到结构体
			var production models.ProductionModel
			if err := cursor.Decode(&production); err != nil {
				log.Println("Error decoding production:", err)
				continue // 如果解码失败，跳过当前记录
			}

			// 保留两位小数
			production.Water = math.Round(production.Water*100) / 100
			production.SpecVol = math.Round(production.SpecVol*100) / 100
			production.Vol = math.Round(production.Vol*100) / 100
			production.Diff = math.Round(production.Diff*100) / 100
			production.Rate = math.Round(production.Rate*100) / 100
			production.WetAcc = math.Round(production.WetAcc*100) / 100
			production.DryAcc = math.Round(production.DryAcc*100) / 100

			productions = append(productions, production)
		}
		// 如果 cursor 在遍历时发生错误，打印错误并返回
		if err := cursor.Err(); err != nil {
			log.Println("Cursor error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 如果没有找到结果，返回提示信息
		if len(productions) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No matching productions found"})
			return
		}
		// 返回查询结果，200 状态码
		c.JSON(http.StatusOK, productions)
	}
}
