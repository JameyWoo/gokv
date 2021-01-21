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
	"strconv"
	"testing"
	"time"
)

// 文件的打开与写入
func TestDbOpen(t *testing.T) {
	db, err := gokv.Open("db0")
	if err != nil {
		logrus.Error(err)
	}
	db.Put(gokv.KeyValue{Key: "hello", Value: "world"})
	val, _ := db.Get("hello")
	logrus.Info("hello: ", val)
	db.Close()
}

// 测试flush的写入以及读取
func TestFlush(t *testing.T) {
	db, err := gokv.Open("db1")
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 11000; i++ {
		db.Put(gokv.KeyValue{Key: strconv.Itoa(i) + "_key", Value: strconv.Itoa(i) + "_value"})
	}
	val, err := db.Get("100_key")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("val:", val)

	val, err = db.Get("10000_key")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("val:", val)
}

// 时间戳
func TestTimeChuo(t *testing.T) {
	logrus.Info(time.Now().UnixNano())
	logrus.Info(time.Now().Unix())
	logrus.Info(time.Now().UnixNano() / 1000)
}

func BenchmarkDbOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db, err := gokv.Open("db0")
		if err != nil {
			logrus.Error(err)
		}
		db.Put(gokv.KeyValue{Key: "hello", Value: "world"})
		_, _ = db.Get("hello")
		db.Close()
	}
}