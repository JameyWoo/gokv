/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/20
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import "github.com/sirupsen/logrus"

// Lru 的 Key, 用来标识一个LRU节点
type LruKey struct {
	Key int // 正式接入中需要进行修改. 最好是写成接口.
}

// Lru 的 Value, 比如key是 fd, Value 则是它的元数据
type LruValue struct {
	Value int
}

// Lru 双向链表的节点. 链表头部为最近被访问的, 尾部为应该剔除的
type LruNode struct {
	k    LruKey   // 键
	v    LruValue // 值
	prev *LruNode // 双向链表
	next *LruNode
}

func (k *LruKey) Compare(k2 *LruKey) bool {
	return false
}

// LRU 结构
type Lru struct {
	size int      // LRU容纳的元素数量
	dl   *LruNode // 双向链表
	m    map[LruKey]*LruNode
}

func (l *Lru) Insert(lk LruKey, lv LruValue) {
	// 先判断 lk 是否在l.m中
	node, ok := l.m[lk]
	if !ok { // 如果不在hash表里面, 则直接插入新节点
		// 如果链表为空
		if l.dl == nil {
			l.dl = &LruNode{k: lk, v: lv}
			l.dl.prev = l.dl
			l.dl.next = l.dl

			l.m[lk] = l.dl
		} else { // 如果不为空
			tmp := &LruNode{k: lk, v: lv}
			tmp.prev = l.dl.prev
			tmp.next = l.dl
			l.dl.prev.next = tmp
			l.dl.prev = tmp

			l.dl = tmp
			l.m[lk] = tmp
		}
		// 如果插入新节点之后, size超了, 那么将最旧的节点替换掉
		if len(l.m) > l.size {
			delKey := l.dl.prev.k
			l.dl.prev.prev.next = l.dl
			l.dl.prev = l.dl.prev.prev
			delete(l.m, delKey) // 更新 hash表
		}
	} else { // 否则修改旧节点, 将旧节点提前
		node.prev.next = node.next
		node.next.prev = node.prev

		node.prev = l.dl.prev
		node.next = l.dl
		l.dl.prev.next = node
		l.dl.prev = node
		l.dl = node
	}
}

func (l *Lru) Get(lk LruKey) (LruValue, bool) {
	node, ok := l.m[lk]
	if ok {
		return node.v, ok
	}
	return LruValue{}, false
}

func (l *Lru) Print() {
	var values []interface{}
	iter := l.dl
	for iter.next != l.dl {
		values = append(values, iter.k)
		iter = iter.next
	}
	values = append(values, iter.k)
	logrus.Info(values)
}

func NewLru(size int) *Lru {
	return &Lru{size: size, m: make(map[LruKey]*LruNode, 0)}
}
