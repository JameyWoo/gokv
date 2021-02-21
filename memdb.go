package gokv

// the memDB MemDB
type MemDB struct {
	memStore *SkipList
	memSize int  // 记录mem存储的容量, put和del的时候进行计算
}

type KeyValue struct {
	Key       string
	Val Value
}

type Value struct {
	Value     string
	Timestamp int64
	Op        Op
}

// 将一个KeyValue结构编码
func (kv *KeyValue) Encode() []byte {
	// 首先获取一个字符串长度的可变长度编码
	keyLen := len(kv.Key)
	valueLen := len(kv.Val.Value)
	enc := make([]byte, 0)
	enc = append(enc, VarIntEncode(keyLen)...)
	enc = append(enc, []byte(kv.Key)...)
	enc = append(enc, VarIntEncode(valueLen)...)
	enc = append(enc, []byte(kv.Val.Value)...)

	// 以毫秒为单位的时间戳, 64位, 共8个字节.
	timestamp := kv.Val.Timestamp
	timeBytes := make([]byte, 8)
	x := 0
	for timestamp > 0 {
		timeBytes[x] = byte(timestamp % int64(256))
		x++
		timestamp /= 256
	}
	enc = append(enc, timeBytes...)

	enc = append(enc, byte(kv.Val.Op))

	return enc
}

// 将一个KeyValue结构解码
// 传入一个 []byte, 读取它的前面部分编码成一个KeyValue结构返回, 同时返回剩下的 []byte
func KvDecode(bytes []byte) (KeyValue, []byte) {
	kv := KeyValue{}
	keyLen, bytes := VarIntDecode(bytes)
	kv.Key = string(bytes[0: keyLen])
	valueLen, bytes := VarIntDecode(bytes[keyLen:])
	kv.Val.Value = string(bytes[0: valueLen])

	timeBytes := bytes[valueLen: valueLen + 8]
	kv.Val.Timestamp = 0
	bias := 1
	for i := 0; i < len(timeBytes); i++ {
		kv.Val.Timestamp += int64(int(timeBytes[i]) * bias)
		bias *= 256
	}
	kv.Val.Op = Op(bytes[valueLen + 8])
	return kv, bytes[valueLen + 9:]
}

func VarIntEncode(x int) []byte {
	bytes := make([]byte, 0)
	for x >= 0x80 {
		bytes = append(bytes, byte(x % 0x80))
		x /= 0x80
	}
	bytes = append(bytes, byte(x + 0x80))
	return bytes
}

// 传入一个[]byte, 读取前面的int, 然后返回剩下的[]byte
func VarIntDecode(bytes []byte) (int, []byte) {
	i := 0
	val := 0
	bias := 1
	for bytes[i] < 0x80 {
		val += bias * int(bytes[i])
		bias *= 0x80
		i++
	}
	val += bias * int(bytes[i] - 0x80)
	return val, bytes[i + 1:]
}

func NewEngine() *MemDB {
	return &MemDB{memStore: NewSkipList(), memSize: 0}
}

func (e *MemDB) Get(key string) (Value, error) {
	m, ok := e.memStore.Get(key)
	if !ok {
		return Value{}, GetEmptyError
	}
	return m.Val, nil
}

// 这个Put只是在内存上的操作, 不涉及磁盘的操作.
func (e *MemDB) Put(kv KeyValue) error {
	e.memStore.Put(KeyValue{Key: kv.Key, Val: kv.Val})
	// TODO: 需要调整一下, 加上8字节的时间戳和1字节的Op
	e.memSize += len(kv.Key) + len(kv.Val.Value) + 8 + 1
	return nil
}

// 删除的元素的value用特殊的字符串来代替
func (e *MemDB) Delete(key string, delTime int64) error {
	kv, ok := e.memStore.Get(key)
	val := kv.Val
	if ok {
		// 这里调整大小只需要调整字符串的大小
		e.memSize = e.memSize - len(val.Value) + len(deleted)
		e.memStore.Put(KeyValue{Key: key, Val: Value{Value: deleted, Timestamp: delTime, Op: DEL}})
	} else {
		err := e.Put(KeyValue{
			Key: key,
			Val: Value{
				Value: deleted,
				Timestamp: delTime,
				Op: DEL,
			},
		})
		return err
	}
	return nil
}

// 先注释掉, 这部分暂时不需要
// 扫描一个区间的key, 得到key value的结果slice
// 如果value为deleted, 那么不添加
//func (e *MemDB) Scan(startKey, endKey string) ([]KeyValue, error) {
//	keys := make([]string, e.memStore.Len())
//	i := 0
//	for k := range e.memStore {
//		keys[i] = k
//		i += 1
//	}
//	// 排序
//	sort.Slice(keys, func(i, j int) bool {
//		return keys[i] < keys[j]
//	})
//	kvs := make([]KeyValue, 0)
//	for _, k := range keys {
//		if k >= startKey && k <= endKey {
//			Value := e.memStore[k]
//			if Value.Value == deleted {  // 如果已删除
//				continue
//			}
//			kvs = append(kvs, KeyValue{LruKey: k, Val: Value})
//		}
//		if k > endKey {
//			break
//		}
//	}
//	return kvs, nil
//}