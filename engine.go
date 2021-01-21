package TinyBase

import (
	"sort"
)

// the memory Engine
type Engine struct {
	memStore map[string]string
	memSize int  // 记录mem存储的容量, put和del的时候进行计算
}

type KeyValue struct {
	Key string
	Value string
}

func NewEngine() *Engine {
	return &Engine{memStore: make(map[string]string), memSize: 0}
}

func (e *Engine) Get(key string) (string, error) {
	m, ok := e.memStore[key]
	if !ok {
		return "", GetEmptyError
	}
	return m, nil
}

// 这个Put只是在内存上的操作, 不涉及磁盘的操作.
func (e *Engine) Put(kv KeyValue) error {
	e.memStore[kv.Key] = kv.Value
	e.memSize += len(kv.Key) + len(kv.Value)
	return nil
}

// 删除的元素的value用特殊的字符串来代替
func (e *Engine) Delete(key string) error {
	val, ok := e.memStore[key]
	if ok {
		e.memSize = e.memSize - len(val) + len(deleted)
		e.memStore[key] = deleted
	} else {
		err := e.Put(KeyValue{key, deleted})
		return err
	}
	return nil
}

// 扫描一个区间的key, 得到key value的结果slice
// 如果value为deleted, 那么不添加
func (e *Engine) Scan(startKey, endKey string) ([]KeyValue, error) {
	keys := make([]string, len(e.memStore))
	i := 0
	for k, _ := range e.memStore {
		keys[i] = k
		i += 1
	}
	// 排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	kvs := make([]KeyValue, 0)
	for _, k := range keys {
		if k >= startKey && k <= endKey {
			value := e.memStore[k]
			if value == deleted {  // 如果已删除
				continue
			}
			kvs = append(kvs, KeyValue{Key: k, Value: value})
		}
		if k > endKey {
			break
		}
	}
	return kvs, nil
}