/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/26
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */
/*
读取是为了什么? 是为了找数据, 找一个key对应的value数据. 并不像写一样要把所有的内容写入
因此读取的接口是传递近一个 key, 然后对各种文件进行查找
对于具体的一个 sstable, 应该是根据他的文件名, 得到一个 os.File, 然后通过这个 os.File 来查找 key

可以建立一个从 文件名 -> os.File 的映射, 从而从LRU缓存中进行查找

读过程:
1. 读取 footer
2. 根据footer得到对应的索引位置
3. 读取index block, 并根据key找到其所在的data block以及其是属于第几个, 从而找到对应的meta block
4. 读取 meta block, 转化成 bloom filter, 如果可能存在该key, 那么找到对应的datab lock
5. 遍历data block, 首先利用 index(16key一间隔) 找到对应key的区间, 之后根据索引顺序遍历
6. 加载了的data block, meta block 以及 需要放入 LRU缓存中

PS. 读缓存部分见 cache.go 中的注释 (有三种读缓存)
*/

package gokv

import (
	"encoding/binary"
	"os"
	"time"
)

type sstReader struct {
	file *os.File
	key  string // 要查找的 key
}

func (r *sstReader) getTest() { // 对 sstReader进行测试的函数, debug
	stat, err := r.file.Stat()
	if err != nil {
		panic("get Stat failed")
	}
	size := stat.Size()
	footerSize := 24
	pFooter := r.getFooter(int(size)-footerSize, footerSize)
	pIndexBlock := r.getIndexBlock(pFooter.indexBlockIndex, int(size)-footerSize-pFooter.indexBlockIndex)
	pMetaindexBlock := r.getMetaindexBlock(pFooter.metaindexBlockIndex,
		pFooter.indexBlockIndex-pFooter.metaindexBlockIndex)
	pMetaBlock := r.getMetaBlock(pMetaindexBlock.offset, pMetaindexBlock.size)
	pDataBlock := r.getDataBlock(pIndexBlock.entries[0].offset, pIndexBlock.entries[0].size)
	_ = pIndexBlock
	_ = pMetaBlock
	_ = pDataBlock

	time.Sleep(1 * time.Second)
}

// 根据文件偏移获得一个 footer指针
// 首先 footer 的size是固定的, 所以根据文件的总size找到 footer 的偏移, 然后分别找到对应的offset
func (r *sstReader) getFooter(offset, len int) *footer {
	aFooter := footer{}
	content := ReadOffsetLen(r.file, offset, len)
	aFooter.metaindexBlockIndex = int(binary.LittleEndian.Uint64(content[:8]))
	aFooter.indexBlockIndex = int(binary.LittleEndian.Uint64(content[8:16]))
	aFooter.magic = int(binary.LittleEndian.Uint64(content[16:]))
	return &aFooter
}

// 根据文件偏移获得一个 metaindexBlock指针
func (r *sstReader) getMetaindexBlock(offset, len int) *metaindexBlock {
	aMetaindexBlock := metaindexBlock{}
	content := ReadOffsetLen(r.file, offset, len)
	aMetaindexBlock.offset = int(binary.LittleEndian.Uint64(content[:8]))
	aMetaindexBlock.count = int(binary.LittleEndian.Uint64(content[8:16]))
	aMetaindexBlock.size = int(binary.LittleEndian.Uint64(content[16:]))
	return &aMetaindexBlock
}

// 根据文件偏移获得一个 indexBlock指针
// 一个 indexBlock 的 entry 的长度是 不固定的, 它由
// keyLenByte, offsetByte, countByte, sizeByte, []byte(item.key))... 组成, 共四个字段
func (r *sstReader) getIndexBlock(offset, len int) *indexBlock {
	aIndexBlock := indexBlock{}
	content := ReadOffsetLen(r.file, offset, len)
	entry := indexEntry{}
	off := 0
	for off < len {
		keyLen := binary.LittleEndian.Uint64(content[:8])
		entry.offset = int(binary.LittleEndian.Uint64(content[8:16]))
		entry.count = int(binary.LittleEndian.Uint64(content[16:24]))
		entry.size = int(binary.LittleEndian.Uint64(content[24:32]))
		entry.key = string(content[32 : 32+keyLen])
		aIndexBlock.entries = append(aIndexBlock.entries, entry)
		entry = indexEntry{}
		off += 32 + int(keyLen)
		content = content[32+int(keyLen):]
	}
	return &aIndexBlock
}

// 根据文件偏移获得一个 metaBlock指针
// 不需要设置 metaBlock 或者 bloom filter 的count, 因为使用的时候用不到
func (r *sstReader) getMetaBlock(offset, len int) *metaBlock {
	// 这个最好配置在文件中, 参数跟写入部分一致
	aMetaBlock := newMetaBlock(2048)
	content := ReadOffsetLen(r.file, offset, len)
	aMetaBlock.bf.decode(content)
	return aMetaBlock
}

// 根据文件偏移获得一个 dataBlock指针
// 先读取 dataBlock indexKey, 再
func (r *sstReader) getDataBlock(offset, len int) *dataBlock {
	aDataBlock := dataBlock{}
	content := ReadOffsetLen(r.file, offset, len)
	indexKeyLen := int(binary.LittleEndian.Uint64(content[len-8:]))
	for i := len - 8 - 8*indexKeyLen; i < len-8; i += 8 {
		aDataBlock.indexKeys = append(aDataBlock.indexKeys, binary.LittleEndian.Uint64(content[i:i+8]))
	}
	aDataBlock.content = content[:len-8-8*indexKeyLen]
	// offset 有一些用处, 可以通过 offset 以及 indexKeys[i] 的差值来计算一个 key 在 datablock的content中的相对偏移
	aDataBlock.offset = offset
	// maxKey 和 count 在内存上都是无效的, 因此可以随便赋一个值
	aDataBlock.maxKey = ""
	aDataBlock.count = 0
	return &aDataBlock
}
