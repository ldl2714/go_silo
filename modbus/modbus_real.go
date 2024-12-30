package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"

	"go_silo/precision_conversion"
	"go_silo/profiles"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReadVol 读取 Vol 地址的值并将值存入 pid 表中

func ReadVol(client *ModbusClient, db *mongo.Database) {
	for _, profile := range profiles.VolProfiles {
		// 读取 Modbus 数据
		// log.Printf("Reading Modbus data from address: %d", profile.Vol)
		results, err := client.ReadHoldingRegisters(uint16(profile.Vol), 2)
		if err != nil {
			log.Printf("Failed to read holding registers from address %d: %v", profile.Vol, err)
			continue
		}

		// 打印读取到的原始数据
		// fmt.Printf("results: %v\n", results)

		// 确保读取到的结果长度足够
		if len(results) < 4 {
			log.Printf("Not enough data read from holding registers at address %d", profile.Vol)
			continue
		}

		// 将读取到的两个 16 位整数转换为 32 位浮点数
		left := binary.BigEndian.Uint16(results[0:2])
		right := binary.BigEndian.Uint16(results[2:4])
		floatValue := precision_conversion.Transform16BitTo32FloatSmall(left, right)

		// 保留两位有效数字
		floatValueStr := fmt.Sprintf("%.2f", floatValue)
		parsedFloatValue, err := strconv.ParseFloat(floatValueStr, 64)
		floatValue = float32(parsedFloatValue)
		if err != nil {
			log.Printf("Failed to parse float value: %v", err)
			continue
		}

		// 打印转换后的浮点数
		// fmt.Printf("floatValue: %.2f\n", floatValue)

		// 更新数据库
		_, err = db.Collection("pid").UpdateOne(context.Background(), bson.M{"id": profile.ID}, bson.M{
			"$set": bson.M{
				"PID_PV": math.Round(float64(floatValue)*100) / 100,
			},
		})
		if err != nil {
			log.Printf("Failed to update pid for ID %s: %v", profile.ID, err)
		}
	}
}

// func ReadPid(client *ModbusClient) {
// 	// 调用 ReadHoldingRegisters 方法读取保持寄存器的数据
// 	addressToRead := uint16(8000) // 替换为你要读取的地址
// 	quantity := uint16(10)        // 替换为你要读取的数量
// 	results, err := client.ReadHoldingRegisters(addressToRead, quantity)
// 	if err != nil {
// 		log.Printf("Failed to read holding registers: %v", err)
// 		return
// 	}
// 	// 确保读取到的结果长度足够
// 	if len(results) < 4 {
// 		log.Printf("Not enough data read from holding registers")
// 		return
// 	}

// 	// 打印读取到的原始数据
// 	fmt.Printf("Holding Registers: %v\n", results)

// 	// 将读取到的两个 16 位整数转换为 32 位浮点数，并拼成一个数组
// 	var floatValues []float32
// 	for i := 0; i < len(results); i += 4 {
// 		if i+4 <= len(results) {
// 			left := binary.BigEndian.Uint16(results[i : i+2])
// 			right := binary.BigEndian.Uint16(results[i+2 : i+4])
// 			floatValue := precision_conversion.Transform16BitTo32FloatSmall(left, right)
// 			floatValues = append(floatValues, floatValue)
// 		}
// 	}

// 	// 打印转换后的浮点数数组
// 	fmt.Printf("Converted Float Values: %v\n", floatValues)
// }
