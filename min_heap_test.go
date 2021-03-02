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

func TestMinHeap(t *testing.T) {
	mh := newMinHeap()
	mh.push(KeyValue{Key: "4", Val: Value{}}, 1)
	mh.push(KeyValue{Key: "1", Val: Value{}}, 1)
	mh.push(KeyValue{Key: "3", Val: Value{}}, 1)
	mh.push(KeyValue{Key: "4", Val: Value{}}, 1)
	mh.push(KeyValue{Key: "6", Val: Value{}}, 1)
	mh.push(KeyValue{Key: "0", Val: Value{}}, 1)

	for {
		item, have := mh.getMin()
		if !have {
			break
		}
		mh.pop()
		logrus.Info(item.kv.Key)
	}
}
