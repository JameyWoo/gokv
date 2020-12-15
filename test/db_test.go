/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: test
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/Jameywoo/TinyBase"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

// 文件的打开与写入
func TestDbOpen(t *testing.T) {
	db, err := TinyBase.Open("db0")
	if err != nil {
		logrus.Error(err)
	}
	db.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})
	val, _ := db.Get("hello")
	logrus.Info(val)
}

// 时间戳
func TestTimeChuo(t *testing.T) {
	logrus.Info(time.Now().UnixNano())
}