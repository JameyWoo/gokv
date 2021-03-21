/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/25
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"crypto/sha1"
	"encoding/binary"
)

/****************************************
*
* 以下是标识sstable各个block部分的内存结构:
	<beginning_of_file>
    [data block 1]
    [data block 2]
    ...
    [data block N]
    [meta block 1]
    ...
    [meta block K]
    [metaindex block]
    [index block]
    [Footer]        (fixed size; starts at file_size - sizeof(Footer))
    <end_of_file>
*
****************************************/

// 一个 dataBlock结构
type dataBlock struct {
	content   []byte
	maxKey    string
	offset    int
	count     int
	indexKeys []uint64 // 内部索引键的位置. 每16个写入一次, 查找时通过这个位置找到对应的key
}

// 将一个 key-value 结构传进来, 一个一个地写入到当前 dataBlock 的 []byte结构中
// 为什么不批量写呢? 因为需要控制 一个dataBlock 的大小
// 每添加一个 key-value, 就需要更新一个 maxKey
func (db *dataBlock) putKV(kv KeyValue) {
	if db.count%16 == 0 {
		db.indexKeys = append(db.indexKeys, uint64(db.offset+len(db.content)))
	}
	// 将单个key编码
	// [key_len:8, value_len:8, key, timestamp:8, op:1, value]
	keyLenByte := make([]byte, 8)
	valueLenByte := make([]byte, 8)
	timestampByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyLenByte, uint64(len(kv.Key)))
	binary.LittleEndian.PutUint64(valueLenByte, uint64(len(kv.Val.Value)))
	binary.LittleEndian.PutUint64(timestampByte, uint64(kv.Val.Timestamp))

	db.content = append(db.content,
		BytesCombine(keyLenByte, valueLenByte, []byte(kv.Key), timestampByte,
			[]byte{byte(kv.Val.Op)}, []byte(kv.Val.Value))...)
	db.count++
	if kv.Key > db.maxKey {
		db.maxKey = kv.Key
	}
}

// 返回当前 datablock 的size, 以 byte为单位
func (db *dataBlock) size() int {
	return len(db.content)
}

// 返回 datablock 的 []byte编码
func (db *dataBlock) encode() []byte {
	var allIndexByte []byte
	indexByte := make([]byte, 8)
	db.content = append(db.content)
	for i := 0; i < len(db.indexKeys); i++ {
		binary.LittleEndian.PutUint64(indexByte, db.indexKeys[i])
		allIndexByte = append(allIndexByte, indexByte...)
	}
	indexCountByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexCountByte, uint64(len(db.indexKeys)))
	return BytesCombine(db.content, allIndexByte, indexCountByte)
}

// 过滤器结构
type metaBlock struct {
	bf    *BloomFilter
	count int
}

// 参数 cnt 代表元素的数量, 默认 2k及 2048
func newMetaBlock(cnt int) *metaBlock {
	return &metaBlock{bf: NewBloomFilterWithArgs(cnt, 3, 3)}
}

// 添加一个key到过滤器中
func (mb *metaBlock) add(key string) {
	mb.bf.Put(key)
	mb.count++
}

// 返回已添加的大小
func (mb *metaBlock) size() int {
	return mb.count
}

// 将一个内存的布隆过滤器编码, 然后
func (mb *metaBlock) encode() []byte {
	return mb.bf.encode()
}

// 过滤器索引结构, 因为每个过滤器的大小结构都是一样的, 所以保存所有的元信息, 不需要单独保存
type metaindexBlock struct {
	offset int // meta block 的起点
	count  int // 所有的meta block的数量
	size   int // 每个 meta block 的大小
}

func (mib *metaindexBlock) set(off, cou, siz int) {
	mib.offset = off
	mib.count = cou
	mib.size = siz
}

func (mib *metaindexBlock) encode() []byte {
	offsetByte := make([]byte, 8)
	countByte := make([]byte, 8)
	sizeByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(offsetByte, uint64(mib.offset))
	binary.LittleEndian.PutUint64(countByte, uint64(mib.count))
	binary.LittleEndian.PutUint64(sizeByte, uint64(mib.size))
	return BytesCombine(offsetByte, countByte, sizeByte)
}

// datablock 索引结构
type indexBlock struct {
	entries []indexEntry
}

// 一个索引的条目结构
type indexEntry struct {
	key    string
	offset int
	count  int
	size   int
}

// 添加一个条目的方法 sst.dataBlock.maxKey, offset, globalCount, len(content)
func (ib *indexBlock) add(key string, offset, globalCount, size int) {
	ib.entries = append(ib.entries, indexEntry{key: key, offset: offset, count: globalCount, size: size})
}

// 将内存的值编码成 字节数组 []byte, 以便写入到文件
func (ib *indexBlock) encode() []byte {
	// 因为每个key的大小不同, 所以需要将key的长度编码之后, 因此一个条目需要保存五个字段
	// key_len 8字节
	// offset 8字节 datablock起点的偏移
	// count 8字节 key的数量
	// size 8字节 block的块字节数
	content := make([]byte, 0)
	for _, item := range ib.entries {
		keyLenByte := make([]byte, 8)
		offsetByte := make([]byte, 8)
		countByte := make([]byte, 8)
		sizeByte := make([]byte, 8)
		binary.LittleEndian.PutUint64(keyLenByte, uint64(len(item.key)))
		binary.LittleEndian.PutUint64(offsetByte, uint64(item.offset))
		binary.LittleEndian.PutUint64(countByte, uint64(item.count))
		binary.LittleEndian.PutUint64(sizeByte, uint64(item.size))
		content = append(content, BytesCombine(keyLenByte, offsetByte, countByte, sizeByte, []byte(item.key))...)
	}
	return content
}

// 文件的结尾, 保留整个文件的元信息
type footer struct {
	metaindexBlockIndex int
	indexBlockIndex     int
	magic               int // 一个特殊字符串的编码
}

// 两个索引都使用 8字节的编码来编码, magic也是使用8字节的编码进行编码
func (ft *footer) encode() []byte {
	bMetaindexBlockIndex := make([]byte, 8)
	bIndexBlockIndex := make([]byte, 8)
	bMagic := make([]byte, 8)
	binary.LittleEndian.PutUint64(bMetaindexBlockIndex, uint64(ft.metaindexBlockIndex))
	binary.LittleEndian.PutUint64(bIndexBlockIndex, uint64(ft.indexBlockIndex))
	shash := sha1.New()
	bMagic = shash.Sum([]byte("wujiahao"))[:8]
	return append(bMetaindexBlockIndex, append(bIndexBlockIndex, bMagic...)...)
}
