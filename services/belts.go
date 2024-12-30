package services

import (
	"context"
	"go_silo/models"
	"go_silo/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 获取 皮带信息
func GetBelt(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		collection := db.Collection("belt")
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.Background())

		var belts []models.BeltModel
		if err = cursor.All(context.Background(), &belts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 查询 material 表中的数据
		materialCollection := db.Collection("material")
		for i, belt := range belts {
			var material models.MaterialModel
			err := materialCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: belt.MaterialId}}).Decode(&material)
			if err != nil {
				log.Printf("Error finding material for belt %v: %v", belt.ID, err)
				continue
			}
			// 将 material 数据替换 belt 表中的 materialId 字段
			belts[i].MaterialId = material.ID
			belts[i].MaterialName = material.MaterialName
			belts[i].MaxWater = material.MaxWater
			belts[i].MinWater = material.MinWater
			belts[i].Water = material.Water

			// 更新 belt 表中的数据
			update := bson.M{
				"$set": bson.M{
					"materialName": material.MaterialName,
					"maxWater":     material.MaxWater,
					"minWater":     material.MinWater,
					"water":        material.Water,
				},
			}
			_, err = collection.UpdateOne(context.Background(), bson.M{"id": belt.ID}, update)
			if err != nil {
				log.Printf("Error updating belt for ID %s: %v", belt.ID, err)
			}
		}

		c.JSON(http.StatusOK, belts)
	}
}

// 修改 皮带配比
func UpdateBeltRatio(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Ratios []models.BeltModel `json:"ratios"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := db.Collection("belt")
		staticCollection := db.Collection("static")
		eventCollection := db.Collection("event")

		// 获取 static 表中的数据
		var staticData models.StaticModel
		err := staticCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: 1}}).Decode(&staticData)
		if err != nil {
			log.Printf("Error finding static data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		for _, ratioUpdate := range requestBody.Ratios {
			// 验证 ID 和 Ratio 是否有效
			if ratioUpdate.ID == "" || ratioUpdate.Ratio < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID or Ratio"})
				return
			}

			// 设置查询条件
			filter := bson.D{{Key: "id", Value: ratioUpdate.ID}}

			// 获取 belt 数据
			var belt models.BeltModel
			err := collection.FindOne(context.Background(), filter).Decode(&belt)
			if err != nil {
				log.Printf("Error finding belt: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// 获取 material 数据
			var material models.MaterialModel
			materialCollection := db.Collection("material")
			err = materialCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: ratioUpdate.MaterialId}}).Decode(&material)
			if err != nil {
				log.Printf("Error finding material: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// 计算 specVol
			var flowRate float64
			if belt.Parent == 1 {
				flowRate = float64(staticData.SetFlowRate1)
			} else if belt.Parent == 2 {
				flowRate = float64(staticData.SetFlowRate2)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent value"})
				return
			}

			specVol := flowRate * (float64(ratioUpdate.Ratio) / 100) / (1 - material.Water/100)
			specVol = float64(int(specVol*100)) / 100 // 保留两位小数

			// 设置更新操作
			update := bson.M{
				"$set": bson.M{
					"ratio":      ratioUpdate.Ratio,
					"materialId": ratioUpdate.MaterialId,
					"specVol":    specVol,
				},
			}

			_, err = collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Printf("Error updating belt ratio for ID %s: %v", ratioUpdate.ID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// 记录事件
			event := models.EventModel{
				Date:  time.Now(),
				Shift: utils.GetShift(),
				Event: "修改" + ratioUpdate.ID + "皮带，配比由" + strconv.Itoa(int(belt.Ratio)) + "%，改为" + strconv.Itoa(int(ratioUpdate.Ratio)) + "%",
			}
			_, err = eventCollection.InsertOne(context.Background(), event)
			if err != nil {
				log.Printf("Error inserting event: %v", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Belt ratios and specVol updated successfully"})
	}
}

// 修改 皮带物料
func UpdateBeltMaterialId(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody models.BeltModel
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证 ID 和 MaterialId 是否有效
		if requestBody.ID == "" || requestBody.MaterialId < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID or MaterialId"})
			return
		}

		collection := db.Collection("belt")
		materialCollection := db.Collection("material")
		eventCollection := db.Collection("event")

		// 获取旧的 belt 数据
		var oldBelt models.BeltModel
		err := collection.FindOne(context.Background(), bson.D{{Key: "id", Value: requestBody.ID}}).Decode(&oldBelt)
		if err != nil {
			log.Printf("Error finding old belt: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 获取旧的 material 数据
		var oldMaterial models.MaterialModel
		err = materialCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: oldBelt.MaterialId}}).Decode(&oldMaterial)
		if err != nil {
			log.Printf("Error finding old material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 获取新的 material 数据
		var newMaterial models.MaterialModel
		err = materialCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: requestBody.MaterialId}}).Decode(&newMaterial)
		if err != nil {
			log.Printf("Error finding new material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 设置查询条件
		filter := bson.D{{Key: "id", Value: requestBody.ID}}

		// 设置更新操作
		update := bson.M{
			"$set": bson.M{
				"materialId":   requestBody.MaterialId,
				"materialName": newMaterial.MaterialName,
				"maxWater":     newMaterial.MaxWater,
				"minWater":     newMaterial.MinWater,
				"water":        newMaterial.Water,
			},
		}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Printf("Error updating belt materialId for ID %s: %v", requestBody.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 记录事件
		event := models.EventModel{
			Date:  time.Now(),
			Shift: utils.GetShift(),
			Event: "修改" + requestBody.ID + "皮带物料，由" + oldMaterial.MaterialName + "，改为" + newMaterial.MaterialName,
		}
		_, err = eventCollection.InsertOne(context.Background(), event)
		if err != nil {
			log.Printf("Error inserting event: %v", err)
		}

		// 同步更新相关的 belt 文档
		relatedIDs := getRelatedIDs(requestBody.ID)
		for _, relatedID := range relatedIDs {
			filter = bson.D{{Key: "id", Value: relatedID}}
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				log.Printf("Error updating related belt materialId for ID %s: %v", relatedID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Belt materialId and related data updated successfully"})
	}
}

// 修改 皮带物料 -- 获取相关的 belt ID
func getRelatedIDs(id string) []string {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return nil
	}

	switch parts[1] {
	case "1":
		return []string{parts[0] + "-2"}
	case "2":
		parentInt, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Error converting parts[0] to int: %v", err)
			return nil
		}
		nextParent := strconv.Itoa((parentInt % 2) + 1)
		return []string{nextParent + "-2"}
	default:
		return nil
	}
}
