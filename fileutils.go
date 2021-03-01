/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: gokv
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package gokv

import (
	"io"
	"os"
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

// 读取一个文件指定偏移之后的指定字节数并返回 []byte
func ReadOffsetLen(f *os.File, offset, len int) []byte {
	res := make([]byte, 0)
	buf := make([]byte, 1024)
	count := 0
	for count < len {
		size, err := f.ReadAt(buf, int64(offset+count))
		if err != nil && err != io.EOF { // 读取到文件结尾时会出现 EOF错误
			panic("ReadOffsetLen failed!")
		}
		count += size
		res = append(res, buf...)
	}
	// 如果读多了, 那么直接截取
	return res[:len]
}
