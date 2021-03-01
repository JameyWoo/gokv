/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/21
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"crypto/sha1"
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"strconv"
	"testing"
	"time"
)

// 对sstable的写入进行测试
func TestSstableWrite(t *testing.T) {
	o := &gokv.Options{ConfigPath: "../gokv.yaml"}
	db, err := gokv.Open("db4", o)
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()
	for i := 0; i < 1100; i++ {
		db.Put(strconv.Itoa(i)+"_key", strconv.Itoa(i)+"_value")
	}
	sst := gokv.NewSSTable(db.Dir(), "test.sst", db.MemDB())
	sst.Write()
	time.Sleep(3 * time.Second)
}

// 读取测试
func TestSstableRead(t *testing.T) {

}

func TestSha1(t *testing.T) {
	n := sha1.New()
	res := n.Sum([]byte("wujiahao"))
	logrus.Info(res[20:])
	time.Sleep(3 * time.Second)
}
