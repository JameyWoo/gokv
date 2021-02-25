/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/25
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestByteSlice(t *testing.T) {
	//var b []byte
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, 3)
	logrus.Info(b)
}

func TestEmptySlice(t *testing.T) {
	var s []int
	s = append(s, 1)
	logrus.Info(s)
}

func TestVarInt2(t *testing.T) {
	b := make([]byte, 8)
	binary.PutVarint(b, 128)
	logrus.Info(b)
}