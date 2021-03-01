/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/21
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package test

import (
	"github.com/Jameywoo/gokv"
	"testing"
)

func TestLru(t *testing.T) {
	lru := gokv.NewLru(4)
	lru.Insert(gokv.LruKey{Key: 1}, gokv.LruValue{Value: 4})
	lru.Insert(gokv.LruKey{Key: 2}, gokv.LruValue{Value: 4})
	lru.Insert(gokv.LruKey{Key: 3}, gokv.LruValue{Value: 4})
	lru.Insert(gokv.LruKey{Key: 1}, gokv.LruValue{Value: 4})

	lru.Print()

	lru.Insert(gokv.LruKey{Key: 5}, gokv.LruValue{Value: 4})

	lru.Print()

	lru.Insert(gokv.LruKey{Key: 6}, gokv.LruValue{Value: 4})
	lru.Insert(gokv.LruKey{Key: 1}, gokv.LruValue{Value: 4})

	lru.Print()

	lru.Insert(gokv.LruKey{Key: 3}, gokv.LruValue{Value: 4})

	lru.Print()
}
