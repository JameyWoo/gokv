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
}

func (r *sstReader) open(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	r.file = file
}

func (r *sstReader) close() {
	r.file.Close()
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

/*
通过已知的 key, 找到其对应的 Value值
步骤:
1. 读取 footer, 然后获得 metaindex block 和 index block 的索引位置
2. 读取 index block, 然后得到key对应与哪个data block
3. 根据累加的count值得到该key属于哪个bloom filter, 然后使用该 bloom filter进行过滤

PS. 需要加入缓存结构
*/
func (r *sstReader) FindKey(key string) (*Value, bool) {
	stat, err := r.file.Stat()
	if err != nil {
		panic("get Stat failed")
	}
	size := stat.Size()
	footerSize := 24
	pFooter := r.getFooter(int(size)-footerSize, footerSize)
	pIndexBlock := r.getIndexBlock(pFooter.indexBlockIndex, int(size)-footerSize-pFooter.indexBlockIndex)
	if key > pIndexBlock.entries[len(pIndexBlock.entries)-1].key {
		return nil, false
	}
	// 搜索 pIndexBlock, 找到key对应的datablock的位置. key 是从小到大的
	cntL, cntR := 0, 0
	blockOffset := -1
	blockLen := -1
	for i := 0; i < len(pIndexBlock.entries); i++ {
		// ! fix bug: 之前这里是 += ; 导致计数过多! 于是debug很久. 发现之后想起之前有一个bug也是由这个引起的, 哎!
		cntR = pIndexBlock.entries[i].count    // 排位的上限
		if key <= pIndexBlock.entries[i].key { // index 里保存的应该是最大值
			blockOffset = pIndexBlock.entries[i].offset
			blockLen = pIndexBlock.entries[i].size
			break // fix bug: 之前忘记了 break, 导致总是查找到最后面的那个block!
		}
		cntL = pIndexBlock.entries[i].count // 排位的下限
	}
	// 如果 offset == -1, 说明 entries 中没有条目(当然, 这种情况基本不会出现)
	if blockOffset == -1 {
		return nil, false
	}
	// 根据 (cntL, cntR] 找到对应的 metaBlock
	pMetaindexBlock := r.getMetaindexBlock(pFooter.metaindexBlockIndex,
		pFooter.indexBlockIndex-pFooter.metaindexBlockIndex)
	metaL := (cntL + 1) / 2048
	metaR := cntR / 2048
	if metaL == metaR { // 范围是一个 metaBlock
		// ! debug: 这里获取 bloom filter的时候, 出现了 len=512 的情况, 导致 content超标. pMetaindexBlock.size的设置有问题
		pMetaBlock := r.getMetaBlock(pMetaindexBlock.offset+metaL*pMetaindexBlock.size, pMetaindexBlock.size)
		if !pMetaBlock.bf.MayContain(key) {
			return nil, false
		}
	} else { // 范围是两个 metaBlock
		pMetaBlockL := r.getMetaBlock(pMetaindexBlock.offset+metaL*pMetaindexBlock.size, pMetaindexBlock.size)
		pMetaBlockR := r.getMetaBlock(pMetaindexBlock.offset+metaR*pMetaindexBlock.size, pMetaindexBlock.size)
		if !pMetaBlockL.bf.MayContain(key) && !pMetaBlockR.bf.MayContain(key) {
			return nil, false
		}
	}
	// 从以 blockOffset 为偏移的 datablock中查找 key
	pDataBlock := r.getDataBlock(blockOffset, blockLen)
	return findKeyFromDataBlock(key, pDataBlock)
}

/*
从一个 dataBlock 中找到 key
步骤:
1. 由于datablock已经解析成了内存数据结构, 因此第一步查找索引, 找到其对应的范围
2. 找到范围后, 顺序遍历这个范围区间, 如果没有找到, 那么这个key不存在在datablock中
*/

func findKeyFromDataBlock(key string, pDataBlock *dataBlock) (*Value, bool) {
	startKey, _ := getKeyByOffset(0, pDataBlock)
	if startKey > key {
		return nil, false
	}
	startOffset := int(pDataBlock.indexKeys[len(pDataBlock.indexKeys)-1]) - pDataBlock.offset
	for i := 0; i < len(pDataBlock.indexKeys); i++ {
		// 相对偏移
		indexKeyOffset := int(pDataBlock.indexKeys[i]) - pDataBlock.offset
		indexKey, _ := getKeyByOffset(indexKeyOffset, pDataBlock)
		if indexKey > key { // 上一个 index区间
			startOffset = int(pDataBlock.indexKeys[i-1]) - pDataBlock.offset
			break
		}
	}
	// 从 startOffset开始依次遍历 datablock 上的 key-value (根据key判断)
	iterOffset := startOffset
	for i := 0; i < 16; i++ { // 最多查找16个key
		// 需要判断offset是否超过datablock的content的len, 因为可能是最后一个
		if iterOffset >= len(pDataBlock.content) {
			return nil, false
		}
		iterKey, step := getKeyByOffset(iterOffset, pDataBlock)
		if iterKey > key { // 遍历到了一个更大的key, 说明待查找的key不在这个block里面
			return nil, false
		}
		if iterKey == key {
			return getValueByOffset(iterOffset, pDataBlock), true
		}
		// 下一个键的偏移
		iterOffset += step
	}
	return nil, false
}

// 在一个 dataBlock中, 通过一个 offset得到一个 key值(string)
// 返回值有两个, 一个是 key值, 第二个是当前key-value 的长度(字节数)
func getKeyByOffset(offset int, pDataBlock *dataBlock) (string, int) {
	keyLenOffset := offset
	keyLen := int(binary.LittleEndian.Uint64(pDataBlock.content[keyLenOffset : keyLenOffset+8]))
	valueLen := int(binary.LittleEndian.Uint64(pDataBlock.content[keyLenOffset+8 : keyLenOffset+16]))
	return string(pDataBlock.content[keyLenOffset+16 : keyLenOffset+16+keyLen]), 25 + keyLen + valueLen
}

// 在一个 dataBlock中, 通过一个 offset得到一个 value值(Value)
func getValueByOffset(offset int, pDataBlock *dataBlock) *Value {
	keyLenOffset := offset
	keyLen := int(binary.LittleEndian.Uint64(pDataBlock.content[keyLenOffset : keyLenOffset+8]))
	valueLen := int(binary.LittleEndian.Uint64(pDataBlock.content[keyLenOffset+8 : keyLenOffset+16]))
	value := Value{}
	value.Timestamp = int64(binary.LittleEndian.Uint64(
		pDataBlock.content[keyLenOffset+keyLen+16 : keyLenOffset+keyLen+24]))
	value.Op = Op(pDataBlock.content[keyLenOffset+keyLen+24])
	value.Value = string(pDataBlock.content[keyLenOffset+keyLen+25 : keyLenOffset+keyLen+25+valueLen])
	return &value
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
// keyLenByte, offsetByte, countByte, sizeByte, []byte(sstableIter.key))... 组成, 共四个字段
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
