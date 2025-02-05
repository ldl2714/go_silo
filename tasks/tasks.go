package tasks

import (
	"context"
	"fmt"
	"go_silo/models"
	"log"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
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

// 每秒更新 belt 表中的diff rate wetAcc dryAcc 和 running 字段
func ProductionAryRunTime(db *mongo.Database) {
	for i := 1; i <= 10; i++ {
		id1 := fmt.Sprintf("%d-1", i)
		id2 := fmt.Sprintf("%d-2", i)
		var belt1, belt2 models.BeltModel

		err := db.Collection("belt").FindOne(context.Background(), bson.M{"id": id1}).Decode(&belt1)
		if err != nil {
			log.Printf("Failed to find belt for ID %s: %v", id1, err)
			continue
		}

		err = db.Collection("belt").FindOne(context.Background(), bson.M{"id": id2}).Decode(&belt2)
		if err != nil {
			log.Printf("Failed to find belt for ID %s: %v", id2, err)
			continue
		}

		if belt1.Vol > 0 {
			belt1.Seconds++
			belt1.WetAcc += belt1.Vol / 3600
			belt1.DryAcc = belt1.WetAcc * (1 - belt1.Water/100)
			belt1.Diff = (belt1.Vol - belt1.SpecVol) / belt1.SpecVol
		}

		if belt2.Vol > 0 {
			belt2.Seconds++
			belt2.WetAcc += belt2.Vol / 3600
			belt2.DryAcc = belt2.WetAcc * (1 - belt2.Water/100)
			belt2.Diff = (belt2.Vol - belt2.SpecVol) / belt2.SpecVol
		}

		if math.Abs(belt1.Diff) <= 0.05 {
			belt1.Numerator++
			belt1.Denominator++
		} else {
			belt1.Denominator++
		}

		belt1.Rate = belt1.Numerator / belt1.Denominator
		if math.Abs(belt2.Diff) <= 0.05 {
			belt2.Numerator++
			belt2.Denominator++
		} else {
			belt2.Denominator++
		}
		belt2.Rate = belt2.Numerator / belt2.Denominator

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": id1}, bson.M{
			"$set": bson.M{
				"seconds":     belt1.Seconds,
				"running":     belt1.Seconds / 60,
				"wetAcc":      math.Round(belt1.WetAcc*100) / 100,
				"dryAcc":      math.Round(belt1.DryAcc*100) / 100,
				"diff":        math.Round(belt1.Diff * 100),
				"numerator":   belt1.Numerator,
				"denominator": belt1.Denominator,
				"rate":        math.Round(belt1.Rate * 100),
			},
		})
		if err != nil {
			log.Printf("Failed to update belt for ID %s: %v", id1, err)
		}

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": id2}, bson.M{
			"$set": bson.M{
				"seconds":     belt2.Seconds,
				"running":     belt2.Seconds / 60,
				"wetAcc":      math.Round(belt2.WetAcc*100) / 100,
				"dryAcc":      math.Round(belt2.DryAcc*100) / 100,
				"diff":        math.Round(belt2.Diff * 100),
				"numerator":   belt2.Numerator,
				"denominator": belt2.Denominator,
				"rate":        math.Round(belt2.Rate * 100),
			},
		})
		if err != nil {
			log.Printf("Failed to update belt for ID %s: %v", id2, err)
		}
	}
}

// 定时存储 ProductionTable 班产报表

func ProductionTable(db *mongo.Database) {
	for i := 1; i <= 10; i++ {
		id1 := fmt.Sprintf("%d-1", i)
		id2 := fmt.Sprintf("%d-2", i)
		var belt1, belt2 models.BeltModel

		err := db.Collection("belt").FindOne(context.Background(), bson.M{"id": id1}).Decode(&belt1)
		if err != nil {
			log.Printf("Failed to find belt for ID %s: %v", id1, err)
			continue
		}

		err = db.Collection("belt").FindOne(context.Background(), bson.M{"id": id2}).Decode(&belt2)
		if err != nil {
			log.Printf("Failed to find belt for ID %s: %v", id2, err)
			continue
		}

		var static models.StaticModel
		err = db.Collection("static").FindOne(context.Background(), bson.M{}).Decode(&static)
		if err != nil {
			log.Printf("Failed to find static data: %v", err)
			continue
		}

		date := time.Now().Format("2006-1-2")

		production1 := models.ProductionModel{
			ID:           belt1.ID,
			MaterialName: belt1.MaterialName,
			Parent:       belt1.Parent,
			Rate:         belt1.Rate,
			Running:      belt1.Running,
			SpecVol:      belt1.SpecVol,
			Vol:          belt1.Vol,
			Water:        belt1.Water,
			WetAcc:       belt1.WetAcc,
			DryAcc:       belt1.DryAcc,
			Diff:         belt1.Diff,
			Ratio:        belt1.Ratio,
			Shift:        static.Shift,
			Team:         static.Team,
			Date:         date,
		}

		production2 := models.ProductionModel{
			ID:           belt2.ID,
			MaterialName: belt2.MaterialName,
			Parent:       belt2.Parent,
			Rate:         belt2.Rate,
			Running:      belt2.Running,
			SpecVol:      belt2.SpecVol,
			Vol:          belt2.Vol,
			Water:        belt2.Water,
			WetAcc:       belt2.WetAcc,
			DryAcc:       belt2.DryAcc,
			Diff:         belt2.Diff,
			Ratio:        belt2.Ratio,
			Shift:        static.Shift,
			Team:         static.Team,
			Date:         date,
		}

		_, err = db.Collection("production").InsertOne(context.Background(), production1)
		if err != nil {
			log.Printf("Failed to insert production for ID %s: %v", id1, err)
		}

		_, err = db.Collection("production").InsertOne(context.Background(), production2)
		if err != nil {
			log.Printf("Failed to insert production for ID %s: %v", id2, err)
		}

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": id1}, bson.M{
			"$set": bson.M{
				"wetAcc":      0,
				"dryAcc":      0,
				"rate":        0,
				"running":     0,
				"seconds":     0,
				"numerator":   0,
				"denominator": 0,
			},
		})
		if err != nil {
			log.Printf("Failed to reset belt for ID %s: %v", id1, err)
		}

		_, err = db.Collection("belt").UpdateOne(context.Background(), bson.M{"id": id2}, bson.M{
			"$set": bson.M{
				"wetAcc":      0,
				"dryAcc":      0,
				"rate":        0,
				"running":     0,
				"seconds":     0,
				"numerator":   0,
				"denominator": 0,
			},
		})
		if err != nil {
			log.Printf("Failed to reset belt for ID %s: %v", id2, err)
		}
	}
}

// 将数缓存到redis中  每秒更新一次
func PostNewDateByTrend(db *mongo.Database, redisClient *redis.Client) {
	for parentId := 1; parentId <= 2; parentId++ {
		for i := 1; i <= 10; i++ {
			id := fmt.Sprintf("%d-%d", i, parentId)
			formattedDateTime := time.Now().Format("2006-01-02")
			cacheKeyVol := fmt.Sprintf("%s:%s:vol", formattedDateTime, id)
			cacheKeySpecVol := fmt.Sprintf("%s:%s:specVol", formattedDateTime, id)

			cachedVolValue, err := redisClient.Get(context.Background(), cacheKeyVol).Result()
			if err != nil && err != redis.Nil {
				log.Printf("Failed to get cached vol value for key %s: %v", cacheKeyVol, err)
				continue
			}

			cachedSpecVolValue, err := redisClient.Get(context.Background(), cacheKeySpecVol).Result()
			if err != nil && err != redis.Nil {
				log.Printf("Failed to get cached specVol value for key %s: %v", cacheKeySpecVol, err)
				continue
			}

			var belt models.BeltModel
			err = db.Collection("belt").FindOne(context.Background(), bson.M{"id": id}).Decode(&belt)
			if err != nil {
				log.Printf("Failed to find belt for ID %s: %v", id, err)
				continue
			}

			vol := math.Round(belt.Vol*100) / 100
			specVol := math.Round(belt.SpecVol*100) / 100

			if cachedVolValue == "" {
				err = redisClient.Set(context.Background(), cacheKeyVol, fmt.Sprintf("%.2f", vol), 48*time.Hour).Err()
				if err != nil {
					log.Printf("Failed to set cached vol value for key %s: %v", cacheKeyVol, err)
				}
				err = redisClient.Set(context.Background(), cacheKeySpecVol, fmt.Sprintf("%.2f", specVol), 48*time.Hour).Err()
				if err != nil {
					log.Printf("Failed to set cached specVol value for key %s: %v", cacheKeySpecVol, err)
				}
			} else {
				newVolValue := fmt.Sprintf("%s,%.2f", cachedVolValue, vol)
				newSpecVolValue := fmt.Sprintf("%s,%.2f", cachedSpecVolValue, specVol)
				err = redisClient.Set(context.Background(), cacheKeyVol, newVolValue, 48*time.Hour).Err()
				if err != nil {
					log.Printf("Failed to update cached vol value for key %s: %v", cacheKeyVol, err)
				}
				err = redisClient.Set(context.Background(), cacheKeySpecVol, newSpecVolValue, 48*time.Hour).Err()
				if err != nil {
					log.Printf("Failed to update cached specVol value for key %s: %v", cacheKeySpecVol, err)
				}
			}
		}
	}
}

// 将缓存数据存入mongo数据库中  每天的 00:02 执行任务
func SaveHistoryVol(db *mongo.Database, redisClient *redis.Client) {
	dateKey := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	beltIds := []string{
		"1-1", "1-2",
		"2-1", "2-2",
		"3-1", "3-2",
		"4-1", "4-2",
		"5-1", "5-2",
		"6-1", "6-2",
		"7-1", "7-2",
		"8-1", "8-2",
		"9-1", "9-2",
		"10-1", "10-2",
	}

	var trendDocuments []interface{}
	for _, beltId := range beltIds {
		volKey := fmt.Sprintf("%s:%s:vol", dateKey, beltId)
		specVolKey := fmt.Sprintf("%s:%s:specVol", dateKey, beltId)

		vol, err := redisClient.Get(context.Background(), volKey).Result()
		if err != nil && err != redis.Nil {
			log.Printf("Failed to get cached vol value for key %s: %v", volKey, err)
			continue
		}

		specVol, err := redisClient.Get(context.Background(), specVolKey).Result()
		if err != nil && err != redis.Nil {
			log.Printf("Failed to get cached specVol value for key %s: %v", specVolKey, err)
			continue
		}

		trendDocuments = append(trendDocuments, models.TrendModel{
			Date:    dateKey,
			ID:      beltId,
			Vol:     vol,
			SpecVol: specVol,
		})
	}

	if len(trendDocuments) > 0 {
		_, err := db.Collection("trend").InsertMany(context.Background(), trendDocuments)
		if err != nil {
			log.Printf("Failed to insert trend documents: %v", err)
		}
	}
}

// 每秒更新一次 alert表中的数据 用于实时报警
func AlertTask(db *mongo.Database) {
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 4; j++ {
			id := fmt.Sprintf("%d-%d", i, j)
			var pid models.PidModel
			var alert models.AlertModel

			err := db.Collection("pid").FindOne(context.Background(), bson.M{"id": id}).Decode(&pid)
			if err != nil {
				log.Printf("Failed to find pid for ID %s: %v", id, err)
				continue
			}

			err = db.Collection("alert").FindOne(context.Background(), bson.M{"id": id}).Decode(&alert)
			if err != nil {
				log.Printf("Failed to find alert for ID %s: %v", id, err)
				continue
			}

			diff := ((pid.PID_PV - pid.PID_SP) / pid.PID_SP) * 100

			if !pid.PID_MAN && pid.PID_PV > 0 && math.Abs(diff) > 5 {
				alert.DiffSeconds++
			} else {
				alert.DiffSeconds = 0
			}

			if alert.DiffSeconds < 30 {
				alert.DiffAlert = 0
			} else if alert.DiffSeconds >= 30 && alert.DiffSeconds < 90 {
				alert.DiffAlert = 1
			} else if alert.DiffSeconds >= 90 {
				alert.DiffAlert = 2
			}

			alert.Diff = math.Round(diff*100) / 100

			_, err = db.Collection("alert").UpdateOne(context.Background(), bson.M{"id": id}, bson.M{
				"$set": bson.M{
					"diff":        alert.Diff,
					"diffSeconds": alert.DiffSeconds,
					"diffAlert":   alert.DiffAlert,
				},
			})
			if err != nil {
				log.Printf("Failed to update alert for ID %s: %v", id, err)
			}
		}
	}
}
