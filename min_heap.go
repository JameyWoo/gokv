/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/1
 * @Desc: 最小堆
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

// TODO: 需要测试一下这个多路合并的效果

package gokv

// heap 里的一个元素
type heapItem struct {
	kv    KeyValue
	index int // 索引, 用来判断属于哪个 sstableIter
}

// 小于比较
func (hi *heapItem) Less(hi2 *heapItem) bool {
	if hi.kv.Key != hi2.kv.Key {
		return hi.kv.Key < hi2.kv.Key
	}
	return hi.kv.Val.Timestamp < hi2.kv.Val.Timestamp
}

// 专为 KeyValue结构设计的最小堆
type minHeap struct {
	size  int
	items []heapItem
}

// 构造函数, 不给参数, size自动适配
func newMinHeap() *minHeap {
	mh := minHeap{}
	mh.size = 0
	mh.items = make([]heapItem, 1) // 将位置为0的kv空出来, 方便计算子节点 (左: 2*x, 右: 2*x+1)
	return &mh
}

// 获得最小的 KeyValue
// 这个比较不是简单地比较key, 如果key相同, 那么timestamp越小的kv越小
func (h *minHeap) getMin() (heapItem, bool) {
	if h.size == 0 {
		return heapItem{}, false
	}
	return h.items[1], true
}

// 添加元素, 上浮, 直到不满足条件时停止
func (h *minHeap) push(kv KeyValue, index int) {
	h.size++
	h.items = append(h.items, heapItem{kv: kv, index: index})
	idx := h.size
	for idx > 1 {
		if h.items[idx].Less(&h.items[idx/2]) {
			tmp := h.items[idx]
			h.items[idx] = h.items[idx/2]
			h.items[idx/2] = tmp
			idx /= 2
		} else {
			break
		}
	}
}

// 弹出顶端元素. 将最后一个元素放到顶端, 然后让它下沉
func (h *minHeap) pop() {
	if h.size == 0 {
		return
	}
	h.items[1] = h.items[h.size]
	h.items = h.items[:len(h.items)-1]
	h.size--
	// 开始下沉
	idx := 1
	next := -1
	for idx < h.size {
		if idx*2 > h.size {
			break
		} else if idx*2+1 > h.size {
			// 父节点只需要和左子节点比较
			next = idx * 2
		} else { // 两个子节点都在, 选出一个更小的
			if h.items[idx*2].Less(&h.items[idx*2+1]) {
				next = idx * 2
			} else {
				next = idx*2 + 1
			}
		}
		if h.items[next].Less(&h.items[idx]) {
			// 需要下沉
			tmp := h.items[idx]
			h.items[idx] = h.items[next]
			h.items[next] = tmp
			idx = next
		} else {
			// 不需要下沉了, 直接退出
			break
		}
	}
}
