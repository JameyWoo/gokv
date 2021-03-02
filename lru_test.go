/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/21
 * @Desc: test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"testing"
)

func TestLru(t *testing.T) {
	lru := NewLru(4)
	lru.Insert(LruKey{Key: 1}, LruValue{Value: 4})
	lru.Insert(LruKey{Key: 2}, LruValue{Value: 4})
	lru.Insert(LruKey{Key: 3}, LruValue{Value: 4})
	lru.Insert(LruKey{Key: 1}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(LruKey{Key: 5}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(LruKey{Key: 6}, LruValue{Value: 4})
	lru.Insert(LruKey{Key: 1}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(LruKey{Key: 3}, LruValue{Value: 4})

	lru.Print()
}
