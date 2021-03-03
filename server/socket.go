/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/3
 * @Desc: main
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 从客户端读取字节流
// 如何判断远程client是否断开连接?
func ReceiveBytesFromConn(conn net.Conn) ([]byte, error) {
	// 先读取 4 字节, 作为长度
	lengthByte := make([]byte, 4)
	conn.Read(lengthByte)                                 // 忽略错误
	length := int(binary.LittleEndian.Uint32(lengthByte)) // 小端序读取
	// TODO: 注意这里如果是传输比较大的文件的话, 是否需要拆分成小的段?
	inputByte := make([]byte, length)   // 输入命令
	length, err := conn.Read(inputByte) // 忽略错误
	return inputByte, err
}

// 向服务器发送命令
// 先计算数据长度, 然后拼接
func SendBytesToConn(conn net.Conn, inputStrByte []byte) {
	length := len(inputStrByte)
	preSend := BytesCombine(IntToBytes(length), inputStrByte)
	_, err := conn.Write(preSend)
	if err != nil {
		panic(err)
	}
}

// 合并两个 []byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
