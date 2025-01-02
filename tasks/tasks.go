package tasks

import (
	"context"
	"fmt"
	"log"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PidModel struct {
	ID     string  `bson:"id"`
	PID_PV float64 `bson:"pid_pv"`
}

// 读取 pid 表中的 pv 值并进行加和 为vol，然后更新 belt 表
func UpdateBeltVol(db *mongo.Database) {
	for i := 1; i <= 10; i++ {
		id1 := fmt.Sprintf("%d-1", i)
		id2 := fmt.Sprintf("%d-2", i)
		id3 := fmt.Sprintf("%d-3", i)
		id4 := fmt.Sprintf("%d-4", i)

		var pid1, pid2, pid3, pid4 PidModel

		err := db.Collection("pid").FindOne(context.Background(), bson.M{"id": id1}).Decode(&pid1)
		if err != nil {
			log.Printf("Failed to find pid for ID %s: %v", id1, err)
			continue
		}

		err = db.Collection("pid").FindOne(context.Background(), bson.M{"id": id2}).Decode(&pid2)
		if err != nil {
			log.Printf("Failed to find pid for ID %s: %v", id2, err)
			continue
		}

		err = db.Collection("pid").FindOne(context.Background(), bson.M{"id": id3}).Decode(&pid3)
		if err != nil {
			log.Printf("Failed to find pid for ID %s: %v", id3, err)
			continue
		}

		err = db.Collection("pid").FindOne(context.Background(), bson.M{"id": id4}).Decode(&pid4)
		if err != nil {
			log.Printf("Failed to find pid for ID %s: %v", id4, err)
			continue
		}

		// 计算加和
		vol1 := pid1.PID_PV + pid2.PID_PV
		vol2 := pid3.PID_PV + pid4.PID_PV

		// 更新 belt 表
		beltID1 := fmt.Sprintf("%d-1", i)
		beltID2 := fmt.Sprintf("%d-2", i)

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": beltID1}, bson.M{
			"$set": bson.M{
				"vol": math.Round(float64(vol1)*100) / 100,
			},
		})
		if err != nil {
			log.Printf("Failed to update belt for ID %s: %v", beltID1, err)
		}

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": beltID2}, bson.M{
			"$set": bson.M{
				"vol": math.Round(float64(vol2)*100) / 100,
			},
		})
		if err != nil {
			log.Printf("Failed to update belt for ID %s: %v", beltID2, err)
		}
	}
}
