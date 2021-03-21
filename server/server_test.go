/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/3
 * @Desc: main
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestOutput(t *testing.T) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, 258)
	logrus.Info("-" + string(bs) + "-")
	logrus.Info(bs[0])
	logrus.Info(bs[1])
	logrus.Info(bs[2])
	logrus.Info(bs[3])
}
