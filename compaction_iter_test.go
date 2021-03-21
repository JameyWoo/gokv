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
	//file, err := os.Open("test/compaction_test/iter_test.sst")
	//file, err := os.Open("test/compaction_test/test1.sst")
	//file, err := os.Open("test/compaction_test/test2.sst")
	file, err := os.Open("test/compaction_test/1614850309913033500.sst")
	//file, err := os.Open("test/db6/1614770311758931700.sst")
	if err != nil {
		panic(err)
	}
	si := newSSTableIter(file)
	for {
		kv, more := si.Next()
		if !more {
			break
		}
		logrus.Info(kv.Key)
	}
}
