/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/2
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestCompare(t *testing.T) {
	if "1_key" > "10_key" {
		logrus.Info("yes")
	}
}

// 实现一个可比较接口
type LruKey interface {
	Equal(key LruKey) bool // 两接口是否相等
}

type LruK struct {
	key string
}

func (k LruK) Equal(key LruKey) bool {
	other := key.(LruK)
	return other.key == k.key
}

func TestKeyCompareInterface(t *testing.T) {
	a := LruK{key: "hello"}
	b := LruK{key: "hello"}
	c := LruK{key: "world"}
	logrus.Info(a.Equal(b))
	logrus.Info(a.Equal(c))
}
