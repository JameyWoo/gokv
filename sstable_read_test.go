/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/1
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestReadBlocks(t *testing.T) {
	sstR := sstReader{}
	file, err := os.Open("test/read_test.sst")
	if err != nil {
		panic(err)
	}
	sstR.file = file
	sstR.getTest()
}

func TestTest(t *testing.T) {
	logrus.Info("hello, gokv")
}

func TestFindKey(t *testing.T) {
	sstR := sstReader{}
	defer sstR.close()
	sstR.open("test/db6/1614770312880889600.sst")
	value, ok := sstR.FindKey(IntToStringWithZero8(12345))
	sstR.close()
	if ok {
		logrus.Info("value: ", value)
	} else {
		logrus.Info("find nothing")
	}
	time.Sleep(1 * time.Second)
}

// 能够find每个 key, 这说明这里的查找和读取逻辑是正确的
func TestFindAll(t *testing.T) {
	sstR := sstReader{}
	file, err := os.Open("test/db5/1616305839294176700.sst")
	if err != nil {
		panic(err)
	}
	sstR.file = file
	for i := 800; i < 1200; i++ {
		value, ok := sstR.FindKey(strconv.Itoa(i) + "_key")
		if ok {
			logrus.Info("value: ", value)
		} else {
			logrus.Info("find nothing")
		}
	}

	time.Sleep(1 * time.Second)
}

func getOnce() {
	sstR := sstReader{}
	defer sstR.close()
	sstR.open("test/db6/1614770311758931700.sst")
	value, ok := sstR.FindKey(IntToStringWithZero8(1277))
	sstR.close()
	if ok {
		_ = value
		//logrus.Info("value: ", value)
	} else {
		logrus.Info("find nothing")
	}
}

// 读缓存的测试. 连续read同一个文件两次, 看看第二次的时候是否会经过缓存; 通过查找一个key来体现
func TestReadDataBlockCache(t *testing.T) {
	getOnce()
	// 第二次
	getOnce()
}

/*
对缓存性能做测试;
开启缓存时: BenchmarkGetOnce-4   	   10000	    104201 ns/op
关闭缓存时: BenchmarkGetOnce-4   	    4220	    266825 ns/op

这啥呀!!! 为什么关闭缓存还快几倍???

原因: 应该是因为测试范围很小, 操作系统(以及go标准库)就给我缓存了这个文件, 所以很快. 而维护缓存还需要一定的时间.

结论: 这个测试不科学!!!
*/
func BenchmarkGetOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getOnce()
	}
}
