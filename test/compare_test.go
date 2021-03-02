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
