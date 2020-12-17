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
	"github.com/sirupsen/logrus"
	"os"
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
			logrus.Error(err)
		}
		wal = wa
	} else {
		wa, err := os.OpenFile(walPath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Error(err)
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
	return db.eng.Put(kv)
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

// Put的wal
func (db *DB) writeLogPut(kv KeyValue) error {
	wal := db.wal
	write := bufio.NewWriter(wal)
	_, err := write.WriteString((time.Now().String()) + ": put\n")
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
	}
	_ = write.Flush()
	return err
}