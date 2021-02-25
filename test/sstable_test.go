/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/21
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"strconv"
	"testing"
)

// 对sstable的写入进行测试
func TestSstableWrite(t *testing.T) {
	db, err := gokv.Open("db3")
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 1100; i++ {
		db.Put(strconv.Itoa(i)+"_key", strconv.Itoa(i)+"_value")
	}
	iter := db.MemIterator()
	for {
		kv, ok := iter.Next()
		if !ok {
			break
		}
		logrus.Info(kv)
	}
}
