/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/25
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

// 检查文件是否存在
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

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 将多个 []byte 合并成一个
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 获取用字符串表示的时间
func GetTimeString() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// 转化为字符串, 不够的位补零
func IntToStringWithZero8(x int) string {
	return fmt.Sprintf("%08d", x)
}
