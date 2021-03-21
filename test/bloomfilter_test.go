/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: test
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := gokv.NewBloomFilter()
	bf.Put("hello")
	bf.Put("world")
	bf.Put("shit")
	bf.Put("abcd")
	bf.Put("event")
	bf.Put("event1")
	logrus.Info(bf.MayContain("event"))
	logrus.Info(bf.MayContain("nihao"))
}
