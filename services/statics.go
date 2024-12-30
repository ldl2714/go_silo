package services

import (
	"go_silo/models"
	"go_silo/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 查询所有文档
func GetStatic(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 gin.Context 传递的上下文
		ctx := c.Request.Context()
		// 查询条件
		filter := bson.D{}
		// 获取集合
		collection := db.Collection("static")

		// 查询数据
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Printf("Error finding statics: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer cursor.Close(ctx) // 确保查询结束时关闭 cursor
		var statics []models.StaticModel

		for cursor.Next(ctx) {
			// 解码单个 static 到结构体
			var static models.StaticModel
			if err := cursor.Decode(&static); err != nil {
				log.Printf("Error decoding static: %v, skipping document", err)
				continue // 如果解码失败，跳过当前记录
			}
			// 获取当前班次并更新 shift 字段
			shift := utils.GetShift()
			static.Shift = shift

			// 更新 static 表中的 shift 字段
			update := bson.M{
				"$set": bson.M{
					"shift": shift,
				},
			}
			_, err = collection.UpdateOne(ctx, bson.M{"id": static.ID}, update)
			if err != nil {
				log.Printf("Error updating static shift for ID %d: %v", static.ID, err)
			}

			statics = append(statics, static)
		}

		//如果 cursor 在遍历时发生错误，打印错误并返回
		if err := cursor.Err(); err != nil {
			log.Printf("Cursor error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 如果没有找到结果，返回提示信息
		if len(statics) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No static models found"})
			return
		}
		// 返回查询结果，200 状态码
		c.JSON(http.StatusOK, statics)
	}
}

// 根据 ID 修改文档
func UpdateStatic(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求中的数据并绑定到 static 结构体
		var static models.StaticModel
		if err := c.ShouldBindJSON(&static); err != nil {
			// 如果请求体解析错误，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证 ID 是否有效
		if static.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// 设置查询条件，查找文档
		filter := bson.D{{Key: "id", Value: static.ID}}

		// 获取现有文档
		var existingStatic models.StaticModel
		collection := db.Collection("static")
		err := collection.FindOne(c, filter).Decode(&existingStatic)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// 如果没有找到文档，返回 404 错误
				c.JSON(http.StatusNotFound, gin.H{"error": "static not found"})
				return
			}
			// 处理其他错误，返回 500 错误
			log.Printf("Error finding static: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 只更新非空字段，避免不必要的修改
		updateFields := bson.M{}
		if static.SetFlowRate1 >= 0 && static.SetFlowRate1 != existingStatic.SetFlowRate1 {
			updateFields["setFlowRate1"] = static.SetFlowRate1
		}
		if static.SetFlowRate2 >= 0 && static.SetFlowRate2 != existingStatic.SetFlowRate2 {
			updateFields["setFlowRate2"] = static.SetFlowRate2
		}
		if static.Team != "" && static.Team != existingStatic.Team {
			updateFields["team"] = static.Team
		}

		// 如果没有字段需要更新，返回提示信息
		if len(updateFields) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No updates made"})
			return
		}

		// 设置更新操作
		update := bson.M{
			"$set": updateFields,
		}

		// 执行更新操作
		_, err = collection.UpdateOne(c, filter, update)
		if err != nil {
			log.Printf("Error updating static: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "static updated successfully"})
	}
}
