/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/25
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTimeString(t *testing.T) {
	logrus.Info(strconv.FormatInt(time.Now().UnixNano(), 10))

	logrus.Info(strconv.FormatInt(10, 2))

	a, b := os.Create("nihao ")
	_, _ = a, b
}
