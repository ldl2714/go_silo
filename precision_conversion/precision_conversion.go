package precision_conversion

import (
	"encoding/binary"
	"math"
)

// Transform16BitTo32FloatSmall 将两个 16 位整数转换为一个 32 位浮点数
func Transform16BitTo32FloatSmall(left uint16, right uint16) float32 {
	// 创建一个 4 字节的缓冲区
	var buffer [4]byte

	// 使用小端字节序将两个 16 位值写入缓冲区
	binary.LittleEndian.PutUint16(buffer[0:], left)  // left 在低位
	binary.LittleEndian.PutUint16(buffer[2:], right) // right 在高位

	// 从缓冲区读取 32 位浮点数
	return math.Float32frombits(binary.LittleEndian.Uint32(buffer[:]))
}

// 读取转换
func TransformToTime(lowWord uint16, highWord uint16) uint32 {
	// 合并高低字
	combinedValue := uint32(lowWord)

	return combinedValue // 返回以秒为单位的时间
}

// 写入转换
func TransToTimeWrite(data uint32) [2]uint16 {
	var buffer [4]byte
	binary.LittleEndian.PutUint32(buffer[:], data)

	lowWord := binary.LittleEndian.Uint16(buffer[0:2])  // 低 16 位
	highWord := binary.LittleEndian.Uint16(buffer[2:4]) // 高 16 位

	return [2]uint16{lowWord, highWord}
}

func Transform32FloatTo16BitSmall(value float32) [2]uint16 {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], math.Float32bits(value))

	left := binary.BigEndian.Uint16(buffer[0:2])
	right := binary.BigEndian.Uint16(buffer[2:4])

	return [2]uint16{right, left}
}
