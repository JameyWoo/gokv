/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: TinyBase
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package TinyBase

import (
	"bufio"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type DB struct {
	eng *Engine
	wal *os.File

	dir string
	walPath string
}

func Open(dirPath string) (*DB, error) {
	// 生成一个文件夹
	// 包含 WAL, HFile目录
	if !Exists(dirPath) {
		os.Mkdir(dirPath, os.ModePerm)
	}
	var wal *os.File
	walPath := dirPath + "/wal.log"
	if !Exists(walPath) {
		wa, err := os.Create(walPath)
		if err != nil {
			logrus.Fatal(err)
		}
		wal = wa
	} else {
		wa, err := os.OpenFile(walPath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Fatal(err)
		}
		wal = wa
	}

	return &DB{eng: NewEngine(), wal: wal, dir: dirPath, walPath: walPath}, nil
}

func (db *DB) Close() {
	db.wal.Close()
}

func (db *DB) Get(key string) (string, error) {
	return db.eng.Get(key)
}

func (db *DB) Put(kv KeyValue) error {
	// TODO: 改进写日志的方式
	db.writeLogPut(kv)
	db.eng.Put(kv)
	// 这个阈值用常量 maxMemSize表示, maxMemSize定义在engine中, 后续改为可配置的量
	if db.eng.memSize >= maxMemSize {
		// 刷到磁盘
		db.flush()
		// TODO: 之后实现异步的flush操作
		db.eng.memStore = make(map[string]string)
		// ! 这个也要重置
		db.eng.memSize = 0
	}
	return nil
}

// 删除的元素的value用特殊的字符串来代替
func (db *DB) Delete(key string) error {
	db.writeLogDelete(key)
	return db.eng.Delete(key)
}

// 扫描一个区间的key, 得到key value的结果slice
// 如果value为deleted, 那么不添加
func (db *DB) Scan(startKey, endKey string) ([]KeyValue, error) {
	return db.eng.Scan(startKey, endKey)
}

// Put的wal, 后续需要将WAL单独写成一个分离的结构
func (db *DB) writeLogPut(kv KeyValue) error {
	wal := db.wal
	write := bufio.NewWriter(wal)
	_, err := write.WriteString((time.Now().String()) + ": put\n")
	if err != nil {
		logrus.Fatal(err)
	}
	_ = write.Flush()
	return err
}

// Delete的wal
func (db *DB) writeLogDelete(key string) error {
	wal := db.wal
	write := bufio.NewWriter(wal)
	_, err := write.WriteString((time.Now().String()) + ": delete\n")
	if err != nil {
		logrus.Fatal(err)
	}
	_ = write.Flush()
	return err
}

// 最简单的flush模式, 直接按顺序追加模式, 不
func (db *DB) flush() error {
	flushPath := db.dir + "/flush_files/"
	_, err := os.Stat(flushPath)
	if os.IsNotExist(err) {
		// create
		os.Mkdir(flushPath, os.ModePerm)
	}
	files, _ := ioutil.ReadDir(flushPath)  // 编号从0开始
	fileId := len(files)

	fileStr := ""
	for key, val := range db.eng.memStore {
		// 编码, [keyLen, key, valueLen, value]
		fileStr += strconv.Itoa(len(key)) + key + strconv.Itoa(len(val)) + val
	}
	// 创建一个新文件
	newFile, err := os.Create(flushPath + strconv.Itoa(fileId) + ".flush")
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	byteLen, err := newFile.WriteString(fileStr)
	if err != nil {
		logrus.Fatal(err)
		return err
	}

	if byteLen != len(fileStr) {
		err = errors.New("byteLen != len(fileStr)")
		logrus.Fatal(err)
		return err
	}
	return nil
}