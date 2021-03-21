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
	"testing"
)

func TestCompact(t *testing.T) {
	sms := make([]sstableMeta, 0)
	sms = append(sms, sstableMeta{dir: "test/compaction_test", filename: "test1.sst"})
	sms = append(sms, sstableMeta{dir: "test/compaction_test", filename: "test2.sst"})
	res := compact(sms)
	logrus.Info(res.dir + "/" + res.filename)
}
