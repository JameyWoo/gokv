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
	"testing"
)

func TestReadFooter(t *testing.T) {
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
