/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/2
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

// 迭代器
type Iterator struct {
	tail *SkipListNode
	now  *SkipListNode
}

// 还要实现一个迭代器, 用来从小到达迭代所有的KeyValue元素
func (sl *SkipList) NewIterator() *Iterator {
	it := &Iterator{
		tail: sl.tail,
		now:  sl.head,
	}
	return it
}

// 获取下一个元素, 两个返回值
// 第一个为 KeyValue结果, 第二个为 ok表示是否存在元素
func (it *Iterator) Next() (KeyValue, bool) {
	if it.now.pointers[0] == it.tail {
		return KeyValue{}, false
	}
	kv := KeyValue{Key: it.now.pointers[0].key, Val: it.now.pointers[0].value}
	it.now = it.now.pointers[0]
	return kv, true
}
