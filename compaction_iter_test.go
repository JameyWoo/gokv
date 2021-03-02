/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/2
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

/*
对 compaction时的 sstable文件的迭代器进行测试, 该迭代器会一步一步地得到sstable的每一个key-value
*/
func TestCompactionIter(t *testing.T) {
	file, err := os.Open("test/compaction_test/iter_test.sst")
	if err != nil {
		panic(err)
	}
	si := newSSTableIter(file)
	for {
		kv, more := si.Next()
		if !more {
			break
		}
		// debug时用, 现在不需要了
		// 将每个 key 都打印出来
		//if kv.Key == "973_key" {
		//	time.Sleep(1 * time.Second)
		//}
		logrus.Info(kv.Key)
	}
}
