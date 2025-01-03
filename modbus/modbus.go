package modbus

import (
	"fmt"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

type ModbusClient struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
	address string
}

func NewModbusClient(address string) *ModbusClient {
	handler := modbus.NewTCPClientHandler(address)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1

	client := modbus.NewClient(handler)

	return &ModbusClient{
		handler: handler,
		client:  client,
		address: address,
	}
}

func (mc *ModbusClient) Connect() error {
	err := mc.handler.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	return nil
}

func (mc *ModbusClient) Close() {
	mc.handler.Close()
}

func (mc *ModbusClient) AutoReconnect() {
	for {
		err := mc.Connect()
		if err == nil {
			fmt.Println("Connected to Modbus server")
			break
		}
		fmt.Println("Failed to connect, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (mc *ModbusClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	results, err := mc.client.ReadCoils(address, quantity)
	if err != nil {
		mc.AutoReconnect()
		return mc.client.ReadCoils(address, quantity)
	}
	return results, nil
}

func (mc *ModbusClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	results, err := mc.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		mc.AutoReconnect()
		return mc.client.ReadHoldingRegisters(address, quantity)
	}
	return results, nil
}

// 添加写入线圈的方法
func (mc *ModbusClient) WriteCoil(address uint16, value bool) error {
	var coilValue uint16
	if value {
		coilValue = 0xFF00
	} else {
		coilValue = 0x0000
	}
	_, err := mc.client.WriteSingleCoil(address, coilValue)
	if err != nil {
		mc.AutoReconnect()
		_, err = mc.client.WriteSingleCoil(address, coilValue)
	}
	return err
}

// 添加写入寄存器的方法
func (mc *ModbusClient) WriteRegisters(address uint16, values [2]uint16) error {
	fmt.Println(address, values)
	_, err := mc.client.WriteMultipleRegisters(address, 2, []byte{
		byte(values[0] >> 8), byte(values[0] & 0xFF),
		byte(values[1] >> 8), byte(values[1] & 0xFF),
	})
	if err != nil {
		mc.AutoReconnect()
		_, err = mc.client.WriteMultipleRegisters(address, 2, []byte{
			byte(values[0] >> 8), byte(values[0] & 0xFF),
			byte(values[1] >> 8), byte(values[1] & 0xFF),
		})
	}
	return err
}
func Modbus() *ModbusClient {
	var address = "192.168.2.149:502" // 替换为你的 PLC 地址
	client := NewModbusClient(address)
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to PLC: %v", err)
	} else {
		fmt.Println("PLC连接成功")
	}
	return client
}
