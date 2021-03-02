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
	file, err := os.Open("test/compaction_test/read_test.sst")
	if err != nil {
		panic(err)
	}
	sstR.file = file
	sstR.key = "990_key"
	value, ok := sstR.FindKey()
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
	file, err := os.Open("test/compaction_test/read_test.sst")
	if err != nil {
		panic(err)
	}
	sstR.file = file
	for i := 800; i < 1200; i++ {
		sstR.key = strconv.Itoa(i) + "_key"
		value, ok := sstR.FindKey()
		if ok {
			logrus.Info("value: ", value)
		} else {
			logrus.Info("find nothing")
		}
	}

	time.Sleep(1 * time.Second)
}
