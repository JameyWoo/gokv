/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: test
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// 文件的打开与写入
// TODO: 目前compaction最大的问题在于
func TestDbOpen(t *testing.T) {
	db, err := gokv.Open("db5", &gokv.Options{ConfigPath: "../gokv.yaml"})
	if err != nil {
		logrus.Error(err)
	}
	db.Put("hello", "world")
	val, _ := db.Get("hello")
	logrus.Info("hello: ", val)
	db.Close()
}

// 测试flush的写入以及读取
func TestFlush(t *testing.T) {
	db, err := gokv.Open("db2", &gokv.Options{ConfigPath: "../gokv.yaml"})
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 1100; i++ {
		db.Put(strconv.Itoa(i)+"_key", strconv.Itoa(i)+"_value")
	}
	val, find := db.Get("100_key")
	if !find {
		logrus.Info("nothing!")
	} else {
		logrus.Info("val:", val)
	}

	val, find = db.Get("10000_key")
	if !find {
		logrus.Info("nothing!")
	} else {
		logrus.Info("val:", val)
	}
}

// 时间戳
func TestTimeChuo(t *testing.T) {
	logrus.Info(time.Now().UnixNano())
	logrus.Info(time.Now().Unix())
	logrus.Info("毫秒: ", time.Now().UnixNano()/1e6)
}

func BenchmarkDbOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db, err := gokv.Open("db0", &gokv.Options{ConfigPath: "../gokv.yaml"})
		if err != nil {
			logrus.Error(err)
		}
		db.Put("hello", "world")
		_, _ = db.Get("hello")
		db.Close()
	}
}

// 测试 flush 使用 sstable
func TestFlushSSTable(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	db, err := gokv.Open("db5", o)
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 200000; i++ {
		x := rand.Int() % 1000000
		db.Put(gokv.IntToStringWithZero8(x), gokv.IntToStringWithZero8(x))
	}
	time.Sleep(1 * time.Second)
}

// 测试 flush 使用 sstable
func TestDBReadSSTable(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	db, err := gokv.Open("db6", o)
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 20000; i++ {
		x := rand.Int() % 8000
		if i == 10000 {
			db.Put(gokv.IntToStringWithZero8(12345), "bingo")
		}
		db.Put(gokv.IntToStringWithZero8(x), gokv.IntToStringWithZero8(x))
	}
	db.Put(gokv.IntToStringWithZero8(666666), "v666666")
	logrus.Info("put over")
	v, find := db.Get(gokv.IntToStringWithZero8(666666))
	if find {
		logrus.Info(v.Value)
	} else {
		logrus.Info("nothing!")
	}
	time.Sleep(1 * time.Second)
}
