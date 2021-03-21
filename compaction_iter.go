/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/2
 * @Desc: compaction部分的迭代器实现
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"encoding/binary"
	"os"
)

// 一个 sstable文件遍历的抽象对象
// 这个抽象对外隐藏了 datablock等细节, 提供一个同一接口 Next. 这个接口不断地返回一个键值完整的 KeyValue结构
type sstableIter struct {
	r       *sstReader
	entries []indexEntry
	current int    // 当前是第current个datablock
	offset  int    // 在 sstable中的偏移, 用来找具体的key
	count   int    // 当前datablock的计数, 从0开始
	content []byte // 当前datablock的字节序列. 由于顺序遍历, 所以不需要索引信息
}

/*
构造函数, 包含一个 sstableIter 的初始化, 读取sstable的元数据
步骤:
1. 读取 footer, 从而读取 indexblock的索引
2. 读取 indexblock, 并且将索引保存到 sstableIter 结构中
3. 初始化加载第一个key所属的datablock到 content []byte
*/
func newSSTableIter(file *os.File) *sstableIter {
	si := sstableIter{}
	r := sstReader{
		file: file,
	}
	si.r = &r
	stat, err := r.file.Stat()
	if err != nil {
		panic("get Stat failed")
	}
	size := stat.Size()
	footerSize := 24
	pFooter := r.getFooter(int(size)-footerSize, footerSize)
	pIndexBlock := r.getIndexBlock(pFooter.indexBlockIndex, int(size)-footerSize-pFooter.indexBlockIndex)
	si.entries = pIndexBlock.entries
	si.current = 0
	si.offset = 0
	// 获取第一个datablock
	si.content = r.getDataBlock(pIndexBlock.entries[0].offset, pIndexBlock.entries[0].size).content
	return &si
}

func (si *sstableIter) close() {
	si.r.close()
}

/*
获得当前 sstableIter 的下一个 KeyValue结构, 迭代器向前
第二个参数是 是否结束
每次都需要判断当前key是否是当前datablock的最后一个key, 如果是, 那么需要更换 datablock
*/
func (si *sstableIter) Next() (KeyValue, bool) {
	if si.count == si.entries[si.current].count {
		// 当计数达到当前datablock的值时, 那么更换下一个datablock
		si.current++
		if si.current == len(si.entries) { // 已经没有datablock了
			return KeyValue{}, false
		}
		// fix bug: 下面两个 赋值 原先都是 bug, 测试了好一会发现问题出在这里
		// ! 这里 si.count 并不需要归零! 因为我entries中的count是累加的, 而不是每个datablock单独的. 所以不需要亲临
		//si.count = 0
		// ! 这里取偏移又是从 0 开始了, 所以不能直接使用 si.offset, 因为它是在sstable中全局的偏移
		si.offset = 0
		si.content = si.r.getDataBlock(si.entries[si.current].offset, si.entries[si.current].size).content
	}
	// 读取一个 key-value, 并计算偏移
	keyLenOffset := si.offset
	keyLen := int(binary.LittleEndian.Uint64(si.content[keyLenOffset : keyLenOffset+8]))
	valueLen := int(binary.LittleEndian.Uint64(si.content[keyLenOffset+8 : keyLenOffset+16]))
	key := string(si.content[keyLenOffset+16 : keyLenOffset+16+keyLen])
	value := Value{}
	value.Timestamp = int64(binary.LittleEndian.Uint64(
		si.content[keyLenOffset+keyLen+16 : keyLenOffset+keyLen+24]))
	value.Op = Op(si.content[keyLenOffset+keyLen+24])
	value.Value = string(si.content[keyLenOffset+keyLen+25 : keyLenOffset+keyLen+25+valueLen])
	// step 向前一小步
	si.count++
	si.offset += 25 + keyLen + valueLen
	return KeyValue{Key: key, Val: value}, true
}
