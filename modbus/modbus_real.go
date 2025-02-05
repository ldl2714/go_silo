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
	"go.mongodb.org/mongo-driver/mongo/options"
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

// ReadPid 读取pid
func ReadPid(client *ModbusClient, db *mongo.Database) {
	//定义一个大数组来存储读取到的数据
	var values []interface{}
	// 读取 Modbus 线圈数据
	address := uint16(9000)
	length := uint16(318)

	results, err := client.ReadCoils(address, length)
	if err != nil {
		log.Printf("Failed to read coils from address %d: %v", address, err)
		return
	}

	// 定义要处理的字段和对应的起始地址偏移量
	fields := []struct {
		fieldName string
		offset    uint16
	}{
		{"PID_MAN", 0},
		{"PID_HALT", 1},
		{"PID_P", 2},
		{"PID_I", 3},
		{"PID_D", 4},
		{"PID_D_ON_X", 5},
		{"PID_QMAX", 6},
		{"PID_QMIN", 7},
	}

	// 将读取到的值按照顺序写入到 pid 表的 PID_MAN 和 PID_P 字段中
	for _, field := range fields {
		for i := 0; i < 40; i++ {
			coilAddress := address + field.offset + uint16(i*8)
			byteIndex := (coilAddress - address) / 8
			bitIndex := uint((coilAddress - address) % 8)
			value := (results[byteIndex] & (1 << bitIndex)) != 0

			// 将布尔值转换为 1 和 0
			intValue := 0
			if value {
				intValue = 1
			}

			// 更新数据库
			pidID := fmt.Sprintf("%d-%d", (i/4)+1, (i%4)+1)
			_, err := db.Collection("pid").UpdateOne(context.Background(), bson.M{"id": pidID}, bson.M{
				"$set": bson.M{
					field.fieldName: intValue,
				},
			})
			if err != nil {
				log.Printf("Failed to update %s for ID %s: %v", field.fieldName, pidID, err)
			}
		}
	}
	// 分段读取浮点数据
	startAddress := uint16(8000)
	endAddress := uint16(8960)
	maxLength := uint16(124)
	for addr := startAddress; addr <= endAddress; addr += maxLength {
		length := maxLength
		if addr+length > endAddress {
			length = endAddress - addr + 1
		}

		// log.Printf("Reading holding registers from address %d with length %d", addr, length)

		results, err := client.ReadHoldingRegisters(addr, length)
		if err != nil {
			log.Printf("Failed to read holding registers from address %d: %v", addr, err)
			return
		}

		// 将读取到的数据写入数据库
		for i := 0; i < len(results); i += 4 {
			if i+3 >= len(results) {
				break
			}
			left := binary.BigEndian.Uint16(results[i : i+2])
			right := binary.BigEndian.Uint16(results[i+2 : i+4])

			var value interface{}
			currentAddress := addr + uint16(i/2)

			// 检查是否是时间类型地址
			if (currentAddress >= 8008 && currentAddress <= 8008+22*39 && (currentAddress-8008)%22 == 0) ||
				(currentAddress >= 8010 && currentAddress <= 8010+22*39 && (currentAddress-8010)%22 == 0) ||
				(currentAddress >= 8012 && currentAddress <= 8012+22*39 && (currentAddress-8012)%22 == 0) {
				value = precision_conversion.TransformToTime(left, right)
			} else {
				value = precision_conversion.Transform16BitTo32FloatSmall(left, right)
			}
			values = append(values, value)
		}
	}
	// fmt.Println("All values:", values)
	// 从 values 数组中取出从地址 8000 开始，间隔 22 的 40 个值，并将它们写入数据库的 PID_SP 字段中
	for i := 0; i < 40; i++ {
		spIndex := i * 22 / 2         // 每个值占用 2 个字节，从地址 8000 开始
		biasIndex := (i*22 + 4) / 2   // 每个值占用 2 个字节，从地址 8004 开始
		gainIndex := (i*22 + 6) / 2   // 每个值占用 2 个字节，从地址 8006 开始
		tdIndex := (i*22 + 8) / 2     // 每个值占用 2 个字节，从地址 8008 开始
		tiIndex := (i*22 + 10) / 2    //每个值占用 2 个字节，从地址 8010 开始
		tdLagIndex := (i*22 + 12) / 2 //每个值占用 2 个字节，从地址 8012 开始
		ymaxIndex := (i*22 + 14) / 2  // 每个值占用 2 个字节，从地址 8014 开始
		yminIndex := (i*22 + 16) / 2  // 每个值占用 2 个字节，从地址 8016 开始
		ymanIndex := (i*22 + 18) / 2  // 每个值占用 2 个字节，从地址 8018 开始
		errIndex := (i*22 + 20) / 2   // 每个值占用 2 个字节，从地址 8020 开始
		yIndex := (i*2 + 880) / 2     // 每个值占用 2 个字节，从地址 8880 开始
		if spIndex >= len(values) || tdIndex >= len(values) {
			break
		}

		pidID := fmt.Sprintf("%d-%d", (i/4)+1, (i%4)+1)
		update := bson.M{
			"$set": bson.M{
				"PID_SP":     values[spIndex],
				"PID_BIAS":   values[biasIndex],
				"PID_GAIN":   values[gainIndex],
				"PID_TD":     values[tdIndex],
				"PID_TI":     values[tiIndex],
				"PID_TD_LAG": values[tdLagIndex],
				"PID_YMAX":   values[ymaxIndex],
				"PID_YMIN":   values[yminIndex],
				"PID_YMAN":   values[ymanIndex],
				"PID_ERR":    values[errIndex],
				"PID_Y":      values[yIndex],
			},
		}

		_, err := db.Collection("pid").UpdateOne(context.Background(), bson.M{"id": pidID}, update)
		if err != nil {
			log.Printf("Failed to update PID_SP and PID_TD for ID %s: %v", pidID, err)
		}
	}
}

func MaterialLevel(client *ModbusClient, db *mongo.Database) {
	startAddress := uint16(5576)
	interval := uint16(12)
	numData := 20
	var floatValues []float32

	for i := 0; i < numData; i++ {
		addr := startAddress + uint16(i)*interval
		length := uint16(4) // 每次读取 4 个字节

		results, err := client.ReadHoldingRegisters(addr, length)
		if err != nil {
			log.Printf("Failed to read holding registers from address %d: %v", addr, err)
			return
		}

		// log.Printf("Read results from address %d: %v", addr, results)

		if len(results) < 4 {
			log.Printf("Not enough data read from address %d", addr)
			continue
		}

		left := binary.BigEndian.Uint16(results[0:2])
		right := binary.BigEndian.Uint16(results[2:4])
		floatValue := precision_conversion.Transform16BitTo32FloatSmall(left, right)
		// log.Printf("Converted float value from left: %d, right: %d -> %f", left, right, floatValue)
		floatValues = append(floatValues, floatValue)
	}

	// log.Println("All float values:", floatValues)

	// 将数据写入 level 表中
	for i := 0; i < numData; i += 2 {
		id := (i / 2) + 1
		materialLevel1 := floatValues[i]
		materialLevel2 := floatValues[i+1]

		_, err := db.Collection("level").UpdateOne(
			context.Background(),
			bson.M{"id": id},
			bson.M{
				"$set": bson.M{
					"MaterialLevel1": materialLevel1,
					"MaterialLevel2": materialLevel2,
				},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Printf("Failed to update level data for ID %d: %v", id, err)
		}
	}
}
