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
