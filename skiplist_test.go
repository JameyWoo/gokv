/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/22
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// 测试跳跃表的插入
func TestSkipList(t *testing.T) {
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(rand.Int() % 10)
		logrus.Infof("put: %s", key)
		list.Put(KeyValue{Key: key,
			Val: Value{Value: key + "_v", Timestamp: time.Now().UnixNano() / 1e6, Op: PUT}})
	}
	for i := 0; i < 10; i++ {
		kv, ok := list.Get(strconv.Itoa(i))
		if ok {
			logrus.Info(kv)
		} else {
			logrus.Infof("%d: no kv!", i)
		}
	}
	for i := 0; i < 10; i++ {
		kv := list.FindGE(strconv.Itoa(i))
		logrus.Infof("k: %s", kv.Key)
	}
}

// 测试跳跃表的迭代器, 使其可以从小到达得到元素
// 迭代器是为了在将内存的内容写入到磁盘而建立的
func TestIteration(t *testing.T) {
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(rand.Int() % 10)
		logrus.Infof("put: %s", key)
		list.Put(KeyValue{Key: key,
			Val: Value{Value: key + "_v", Timestamp: time.Now().UnixNano() / 1e6, Op: PUT}})
	}

	it := list.NewIterator()
	for {
		kv, ok := it.Next()
		if !ok {
			break
		}
		logrus.Infof("key: %s, value: %s, timestamp: %d", kv.Key, kv.Val.Value, kv.Val.Timestamp)
	}
}
