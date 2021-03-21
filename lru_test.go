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

type myKey struct {
	Key int
}

type LruValue struct {
	Value int
}

//func (k myKey) Equal(k2 LruKey) bool {
//	other := k2.(myKey)
//	return other.Key == k.Key
//}

func TestLru(t *testing.T) {
	lru := NewLru(4)

	lru.Insert(myKey{Key: 1}, LruValue{Value: 4})
	lru.Insert(myKey{Key: 2}, LruValue{Value: 4})
	lru.Insert(myKey{Key: 3}, LruValue{Value: 4})
	lru.Insert(myKey{Key: 1}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(myKey{Key: 5}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(myKey{Key: 6}, LruValue{Value: 4})
	lru.Insert(myKey{Key: 1}, LruValue{Value: 4})

	lru.Print()

	lru.Insert(myKey{Key: 3}, LruValue{Value: 4})

	lru.Print()
}
