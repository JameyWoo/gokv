/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/21
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
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

func TestPerformance(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	rand.Seed(int64(time.Now().Nanosecond()))
	str := "db" + strconv.Itoa(rand.Int()%10000)
	logrus.Info(str)
	db, err := gokv.Open(str, o)
	value := "hello"
	for i := 0; i < 10; i++ {
		value += value
	}
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		_ = db.Put(gokv.IntToStringWithZero8(rand.Int()%10000000), value)
		i++
		if i%1000 == 0 {
			logrus.Info(i)
		}
	}
}

// 测试大数据及key长度不等的情况的写入, 旨在fix bug
// 这个函数测试的是key顺序递增的情况, 会导致没有compaction, 所以没有出现错误. 于是可以认定错误是compaction阶段
func TestBigNotEqualLen(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	db, err := gokv.Open("./db1", o)
	value := "hello"
	for i := 0; i < 30; i++ {
		value += "_world"
	}
	logrus.Info(value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100000; i++ {
		_ = db.Put(strconv.Itoa(i), value)
	}
	time.Sleep(1 * time.Second)
}

// 测试大数据情况下频繁compaction的例子
func TestBigCompaction(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	db, err := gokv.Open("./db2", o)
	value := "hello"
	for i := 0; i < 30; i++ {
		value += "_world"
	}
	logrus.Info(value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100000; i++ {
		_ = db.Put(strconv.Itoa(rand.Int()%10000), value)
	}
	time.Sleep(1 * time.Second)
}
