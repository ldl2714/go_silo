package modbus

import (
	"encoding/binary"
	"fmt"
	"go_silo/precision_conversion"
	"log"
)

func ReadRegisters(client *ModbusClient) {
	// 调用 ReadHoldingRegisters 方法读取保持寄存器的数据
	addressToRead := uint16(8000) // 替换为你要读取的地址
	quantity := uint16(10)        // 替换为你要读取的数量
	results, err := client.ReadHoldingRegisters(addressToRead, quantity)
	if err != nil {
		log.Printf("Failed to read holding registers: %v", err)
		return
	}
	// 确保读取到的结果长度足够
	if len(results) < 4 {
		log.Printf("Not enough data read from holding registers")
		return
	}

	// 打印读取到的原始数据
	fmt.Printf("Holding Registers: %v\n", results)

	// 将读取到的两个 16 位整数转换为 32 位浮点数，并拼成一个数组
	var floatValues []float32
	for i := 0; i < len(results); i += 4 {
		if i+4 <= len(results) {
			left := binary.BigEndian.Uint16(results[i : i+2])
			right := binary.BigEndian.Uint16(results[i+2 : i+4])
			floatValue := precision_conversion.Transform16BitTo32FloatSmall(left, right)
			floatValues = append(floatValues, floatValue)
		}
	}

	// 打印转换后的浮点数数组
	fmt.Printf("Converted Float Values: %v\n", floatValues)
}
