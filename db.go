/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: gokv
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package gokv

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type DB struct {
	memory *Engine
	wal    *os.File

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
	// 查看是否有日志文件, 如果没有则创建
	checkWal(walPath)
	wa, err := os.OpenFile(walPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	wal = wa

	return &DB{memory: NewEngine(), wal: wal, dir: dirPath, walPath: walPath}, nil
}

func checkWal(walPath string) {
	if !Exists(walPath) {
		wa, err := os.Create(walPath)
		if err != nil {
			logrus.Fatal(err)
		}
		wa.Close()
	}
}

func (db *DB) Close() {
	db.wal.Close()
}

func (db *DB) Get(key string) (Value, error) {
	value, err := db.memory.Get(key)
	// 如果没有得到value
	if err == nil {
		return value, err
	}
	// 从磁盘上获取value
	return db.diskGet(key)
}

func (db *DB) Put(key, value string) error {
	return db.put(KeyValue{Key: key,
		Val: Value{Value: value, Timestamp: time.Now().UnixNano() / 1e6, Op: SET}})
}

func (db *DB) put(kv KeyValue) error {
	// TODO: 改进写日志的方式
	db.writeLogPut(kv)
	db.memory.Put(kv)
	// 这个阈值用常量 MaxMemSize表示, MaxMemSize定义在配置文件中, 后续改为可配置的量
	if db.memory.memSize >= config.MaxMemSize {
		//logrus.Info("flush")
		// 刷到磁盘
		db.flush()
		// TODO: 之后实现异步的flush操作
		db.memory.memStore = make(map[string]Value)
		// ! 这个也要重置
		db.memory.memSize = 0
		// 清空 wal.log 文件内容, 用直接创建的方式
		// Create creates or truncates the named file. If the file already exists, it is truncated.
		wal, err := os.Create(db.walPath)
		if err != nil {
			logrus.Fatal(err)
		}
		wal.Close()
		wal, err = os.OpenFile(db.walPath, os.O_WRONLY|os.O_APPEND, 0666)
		db.wal = wal
	}
	return nil
}

// 删除的元素的value用特殊的字符串来代替
func (db *DB) Delete(key string) error {
	db.writeLogDelete(key)
	return db.memory.Delete(key)
}

// 扫描一个区间的key, 得到key value的结果slice
// 如果value为deleted, 那么不添加
func (db *DB) Scan(startKey, endKey string) ([]KeyValue, error) {
	return db.memory.Scan(startKey, endKey)
}

// Put的wal, 后续需要将WAL单独写成一个分离的结构
func (db *DB) writeLogPut(kv KeyValue) error {
	wal := db.wal
	write := bufio.NewWriter(wal)
	_, err := write.WriteString(fmt.Sprintf("%s: put {key: %s, value: %s}\n",
		time.Now().String(), kv.Key, kv.Val))
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

	fileBytes := make([]byte, 0)

	// 有序地flush
	keys := make([]string, 0, len(db.memory.memStore))
	for key, _ := range db.memory.memStore {
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// 编码, [varint_key, key, varint_value, value, Timestamp, Op]
		kv := KeyValue{Key: key, Val: db.memory.memStore[key]}
		//logrus.Info(kv.Key)
		kvBytes := kv.Encode()
		fileBytes = append(fileBytes, kvBytes...)
	}
	// 创建一个新文件, 前面补零
	newFile, err := os.Create(flushPath + fmt.Sprintf("%010d", fileId) + ".flush")
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	byteLen, err := newFile.Write(fileBytes)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	_ = newFile.Sync()

	if byteLen != len(fileBytes) {
		err = errors.New("byteLen != len(fileStr)")
		logrus.Fatal(err)
		return err
	}
	return nil
}

func (db *DB) diskGet(key string) (Value, error) {
	// 从磁盘上获取目录及文件, 然后一个一个读取
	flushPath := db.dir + "/flush_files/"
	_, err := os.Stat(flushPath)
	if os.IsNotExist(err) {
		// create
		os.Mkdir(flushPath, os.ModePerm)
	}
	files, _ := ioutil.ReadDir(flushPath)  // 编号从0开始
	for ii := 0; ii < len(files); ii++ {
		bytes, err := ioutil.ReadFile(flushPath + files[ii].Name())
		if err != nil {
			logrus.Error(err)
		}
		//logrus.Info("file: ", files[ii].Name())
		// 解码
		for len(bytes) != 0 {
			kv := KeyValue{}
			kv, bytes = KvDecode(bytes)
			//logrus.Infof("key: %s, val: %s", kv.Key, kv.Val.Value)
			if key == kv.Key {
				return kv.Val, nil
			}
		}
	}
	return Value{}, GetEmptyError
}
