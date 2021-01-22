/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/22
 * @Desc: 跳跃表, 用来做内存上的有序结构
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import "math/rand"

// 一个跳跃表的节点
// level 从0开始
type SkipListNode struct {
	key string
	value Value
	level int
	pointers []*SkipListNode
}

// 一个完整的跳跃表结构
type SkipList struct {
	maxLevel int  // 最大的层级
	head, tail *SkipListNode  // 跳跃表的头和尾
}

// 这个maxLevel有一个默认值.
// 这个值应该是可以配置的
func NewSkipList() *SkipList {
	sl := &SkipList{maxLevel: 5, head: &SkipListNode{key: "head"}, tail: &SkipListNode{key: "tail"}}
	// 初始化 head和tail的指针
	for i := 0; i <= sl.maxLevel; i++ {
		sl.head.pointers = append(sl.head.pointers, sl.tail)
	}
	return sl
}

type pNode struct {
	node *SkipListNode
	l int
}

// 插入一个KeyValue
func (sl *SkipList) Put(kv KeyValue) {
	key := kv.Key
	p := sl.head
	level := sl.maxLevel
	// 路径上的node
	pathNodes := make([]pNode, 0)
	for p != sl.tail {
		// 当key比当前节点的下一个节点的值小, 或者下一个节点是tail即到了结尾, 那么指针向下走
		if p.pointers[level] == sl.tail || key < p.pointers[level].key {
			// 插入时: 当指针向下的时候, 需要保存该指针, 用来在后面进行插入
			pathNodes = append(pathNodes, pNode{node: p, l: level})

			// 当 level为0的时候, 就插入
			if level == 0 {
				h := sl.randHeight()  // 获取随机高度, 从0开始
				last := len(pathNodes) - 1
				now := &SkipListNode{key: kv.Key, value: kv.Val, level: h, pointers: make([]*SkipListNode, 0)}
				// TODO: 考虑h超过maxLevel的情况
				for i := last; i >= last - h; i-- {
					node := pathNodes[i].node
					l := pathNodes[i].l
					next := node.pointers[l]
					node.pointers[l] = now
					now.pointers = append(now.pointers, next)
				}
				// 插入完之后就可以返回了
				return
			}
			// 指针下移
			level--
		} else if key > p.pointers[level].key {  // 当key比当前的节点的下一个节点的值大, 那么指针向右
			p = p.pointers[level]
		} else {  // key 和 下一节点相等, 可以返回结果了
			// 这种情况是 key 已经存在在跳跃表中. 那么只需要修改其Value
			p.pointers[level].value = kv.Val
			return
		}
	}
}

// 只有Put方法, 没有Delete方法. 因为引擎的Del会变成put

// 获取一个keyValue
func (sl *SkipList) Get(key string) (KeyValue, bool) {
	p := sl.head
	level := sl.maxLevel
	for p != sl.tail {
		// 当key比当前节点的下一个节点的值小, 或者下一个节点是tail即到了结尾, 那么指针向下走
		if p.pointers[level] == sl.tail || key < p.pointers[level].key {
			// 当 level为0的时候, 说明没有找到这个值
			if level == 0 {
				return KeyValue{}, false
			}
			// 指针下移
			level--
		} else if key > p.pointers[level].key {  // 当key比当前的节点的下一个节点的值大, 那么指针向右
			p = p.pointers[level]
		} else {  // key 和 下一节点相等, 可以返回结果了
			return KeyValue{Key: key, Val: p.pointers[level].value}, true
		}
	}
	return KeyValue{}, false
}

// 找到第一个大于等于当前key的节点
func (sl *SkipList) FindGE(key string) KeyValue {
	p := sl.head
	level := sl.maxLevel
	for p != sl.tail {
		// 当key比当前节点的下一个节点的值小, 或者下一个节点是tail即到了结尾, 那么指针向下走
		if p.pointers[level] == sl.tail || key < p.pointers[level].key {
			// 当 level为0的时候, 说明没有找到这个值.
			// 下一个值比他更大
			if level == 0 {
				return KeyValue{Key: p.pointers[level].key, Val: p.pointers[level].value}
			}
			// 指针下移
			level--
		} else if key > p.pointers[level].key {  // 当key比当前的节点的下一个节点的值大, 那么指针向右
			p = p.pointers[level]
		} else {  // key 和 下一节点相等, 可以返回结果了
			return KeyValue{Key: p.pointers[level].key, Val: p.pointers[level].value}
		}
	}
	return KeyValue{}
}

// 获取一个随机的高度值
func (sl *SkipList) randHeight() int {
	// bias 可以调整. 这个值越大, skiplist越稀疏
	bias := 2
	h := 0
	for rand.Int() % bias == 0 {
		h += 1
	}
	return h
}

// 迭代器
type Iterator struct {
	tail *SkipListNode
	now *SkipListNode
}

// 还要实现一个迭代器, 用来从小到达迭代所有的KeyValue元素
func (sl *SkipList) NewIterator() *Iterator {
	it := &Iterator{
		tail: sl.tail,
		now: sl.head,
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