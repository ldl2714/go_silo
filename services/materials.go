package services

import (
	"context"
	"go_silo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 获取物料信息
func GetMaterial(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 设置一个空的查询条件 (相当于 find({}) )
		filter := bson.D{}
		// 获取集合
		collection := db.Collection("material")

		// 查询数据
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Println("Error finding materials:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx) //确保查询结束时关闭 cursor

		// 存储结果
		var materials []models.MaterialModel
		for cursor.Next(ctx) {
			// 解码单个 material 到结构体
			var material models.MaterialModel
			if err := cursor.Decode(&material); err != nil {
				log.Println("Error decoding material:", err)
				continue // 如果解码失败，跳过当前记录
			}
			materials = append(materials, material)
		}
		// 如果 cursor 在遍历时发生错误，打印错误并返回
		if err := cursor.Err(); err != nil {
			log.Println("Cursor error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 如果没有找到结果，返回提示信息
		if len(materials) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No static models found"})
			return
		}
		// 返回查询结果，200 状态码
		c.JSON(http.StatusOK, materials)
	}
}

// 根据id 修改文档
func UpdateMaterial(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求中的数据并绑定到 material 结构体
		var material models.MaterialModel
		if err := c.ShouldBindJSON(&material); err != nil {
			// 如果请求体解析错误，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证 ID 是否有效
		if material.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// 设置查询条件，查找文档
		filter := bson.D{{Key: "id", Value: material.ID}}

		// 获取现有文档
		var existingMaterial models.MaterialModel
		collection := db.Collection("material")
		err := collection.FindOne(c, filter).Decode(&existingMaterial)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// 如果没有找到文档，返回 404 错误
				c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
				return
			}
			// 处理其他错误，返回 500 错误
			log.Printf("Error finding material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 只更新非空字段，避免不必要的修改
		updateFields := bson.D{}
		if material.MaterialName != "" && material.MaterialName != existingMaterial.MaterialName {
			updateFields = append(updateFields, bson.E{Key: "materialName", Value: material.MaterialName})
		}
		if material.MaxWater != 0 && material.MaxWater != existingMaterial.MaxWater {
			updateFields = append(updateFields, bson.E{Key: "maxWater", Value: material.MaxWater})
		}
		if material.MinWater != 0 && material.MinWater != existingMaterial.MinWater {
			updateFields = append(updateFields, bson.E{Key: "minWater", Value: material.MinWater})
		}
		if material.Water != 0 && material.Water != existingMaterial.Water {
			updateFields = append(updateFields, bson.E{Key: "water", Value: material.Water})
		}
		// if material.MaxRatio != 0 && material.MaxRatio != existingMaterial.MaxRatio {
		// 	updateFields = append(updateFields, bson.E{Key: "maxRatio", Value: material.MaxRatio})
		// }
		// if material.MinRatio != 0 && material.MinRatio != existingMaterial.MinRatio {
		// 	updateFields = append(updateFields, bson.E{Key: "minRatio", Value: material.MinRatio})
		// }

		// 如果没有字段需要更新，返回提示信息
		if len(updateFields) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No updates made"})
			return
		}

		// 设置更新操作
		update := bson.D{
			{Key: "$set", Value: updateFields},
		}

		// 执行更新操作
		_, err = collection.UpdateOne(c, filter, update)
		if err != nil {
			log.Printf("Error updating material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "Material updated successfully"})
	}
}

// 新增 文档
func AddMaterial(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var material map[string]interface{}

		// 绑定请求的 JSON 数据
		if err := c.ShouldBindJSON(&material); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查是否传递了必需字段
		requiredFields := []string{"materialName", "maxWater", "minWater", "water", "maxRatio", "minRatio"}
		for _, field := range requiredFields {
			if _, ok := material[field]; !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
				return
			}
		}

		collection := db.Collection("material")

		// 获取当前最大id
		var result map[string]interface{}
		opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
		err := collection.FindOne(context.Background(), bson.D{}, opts).Decode(&result)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("Error finding max id:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 设置新文档的id为最大id加1
		var newID int32
		if result != nil {
			newID = result["id"].(int32) + 1
		} else {
			newID = 1
		}
		material["id"] = newID

		// 插入新文档
		_, err = collection.InsertOne(context.Background(), material)
		if err != nil {
			log.Println("Error inserting material:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Material added successfully"})
	}
}

// 根据id 删除文档
func DeleteMaterial(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求体中获取 id
		var requestBody struct {
			ID int `json:"id"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 检查 id 是否有效
		if requestBody.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// 设置查询条件
		filter := bson.D{{Key: "id", Value: requestBody.ID}}

		collection := db.Collection("material")

		// 删除文档
		result, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Println("Error deleting material:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Material deleted successfully"})
	}
}
