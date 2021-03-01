/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/20
 * @Desc: sstable 的接口.
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

/**
定义 sstable 的格式. 文档看 'doc/sstable格式.md'

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

每个 sstable文件以后缀 .sst 表示, sstable的单位为块 block

写sstable的方法:
1. 利用memdb的迭代器, 顺序遍历 key-value, 同时写入到 datablock
2. 对于datablock, 有一个上限大小, 当 当前的datablock达到这个大小的时候, 就将datablock内容写入到磁盘
	我的datablock的实现不压缩前缀, 全都直接保留原始的值.
	同时, 默认每16个key记录一个块内索引用来快速查找. (对于redis的数据来说, 一般key不会多长, 但是value可能很长. )
3. 在写 datablock的同时, 还需要记录一些元数据以及过滤器数据.
4. 布隆过滤器是按照数据量来写的. 默认每 2k 个key-value写到一个布隆过滤器上. 过滤器有很多, 每个过滤器是一个metablock.
	metaindex block 需要记录每个metaindex block的偏移.
	由于所有的布隆过滤器长度都可以一样, 所以其实所有的metablock可以合并成一个metablock来设计, 可以直接通过长度计算来索引.
	一个 metablock <= 2k. 那么一个metablock的key是大于一个datablock的. 所以知道了datablock的范围就可以确定该key所在的metablock的范围
5. index block 保存了每一个datablock中的最大key, 以及对应的datablock偏移, 和该datablock中的key的数量. 还有这个block的字节数(大小)
	在查找的时候, 先根据key的大小找到其对应的datablock(通过线性遍历). 以及累加的key顺序值(从而找到对应的布隆过滤器)
6.
*/

// sstable, 每个sstable结构对应一个sstable文件.
// sstable文件根据传入的内存memtable将
type SSTable struct {
	dir      string
	filename string
	memdb    *MemDB

	tmpFilename string
	writer      *sstWriter

	dataBlock      dataBlock      // 内存上同时只会存在一个datablock
	metaBlock      []*metaBlock   // 可能有多个布隆过滤器
	metaindexBlock metaindexBlock // 一个 过滤器索引
	indexBlock     indexBlock     // 一个索引块
	footer         footer         // 一个 footer
}

// 构造函数
func NewSSTable(dir, filename string, memdb *MemDB) *SSTable {
	return &SSTable{dir: dir, filename: filename, memdb: memdb, writer: &sstWriter{}}
}

// 创建 sstable 文件并打开
func (sst *SSTable) open() {
	// 取一个临时的名字, 在 close 的时候改名
	sst.tmpFilename = strconv.FormatInt(time.Now().UnixNano(), 10) + ".sst.tmp"
	file, err := os.Create(sst.dir + "/" + sst.tmpFilename)
	sst.writer.file = file
	if err != nil {
		logrus.Panic("sstable open failed")
	}
}

// 关闭文件
func (sst *SSTable) close() {
	// 修改名称
	sst.writer.file.Close()
	os.Rename(sst.dir+"/"+sst.tmpFilename, sst.dir+"/"+sst.filename)
}

// sstable的写方法. 传递一个 memdb进来, 然后将其内容写入到 sstable文件, 最后将memdb删除
func (sst *SSTable) Write() {
	// 先初始化文件描述符, 打开文件
	sst.open()

	// 内存key-value的迭代器
	iter := sst.memdb.NewIterator()
	metaB := newMetaBlock(2048) // 初始化布隆过滤器, 使用构造函数
	offset := 0                 // 全局的偏移
	var content []byte
	globalCount := 0 // 全局的count, 表示当前的 datablock index 的key的数量

	for {
		kv, ok := iter.Next()
		if !ok {
			break
		}
		sst.dataBlock.putKV(kv)
		// 过滤器添加key
		metaB.add(kv.Key)
		// ! 考虑剩下的 dataBlock内容
		if sst.dataBlock.size() > 4096 { // 一个阈值, 要配置
			content = sst.dataBlock.encode()
			offset += len(content)
			// 同时将 dataBlock 的信息写入到 indexBlock 中
			globalCount += sst.dataBlock.count
			// fix bug: 这里的offset应该减去大小, offset是一个block的起点
			sst.indexBlock.add(sst.dataBlock.maxKey, offset-len(content), globalCount, len(content))
			// 将这个dataBlock 的值写入到sstable
			sst.writer.write(content)

			// 重置 dataBlock
			sst.dataBlock = dataBlock{offset: offset}
		}
		// ! 考虑剩下的布隆过滤器内容
		if metaB.size() == 2048 { // 更换下一个布隆过滤器
			sst.metaBlock = append(sst.metaBlock, metaB)
			metaB = newMetaBlock(2048)
		}
	}
	// 还可能剩下一些dataBlock
	if sst.dataBlock.count > 0 {
		content = sst.dataBlock.encode()
		offset += len(content)
		// 同时将 dataBlock 的信息写入到 indexBlock 中
		globalCount += sst.dataBlock.count
		sst.indexBlock.add(sst.dataBlock.maxKey, offset-len(content), globalCount, len(content))
		// 将这个dataBlock 的值写入到sstable
		sst.writer.write(content)

		// 重置 dataBlock
		sst.dataBlock = dataBlock{offset: offset}
	}
	// 还可能剩下一些布隆过滤器内容
	if metaB.count > 0 {
		sst.metaBlock = append(sst.metaBlock, metaB)
	}

	// 向文件中写入 metaBlock
	metaBlockOffset := offset
	for i := 0; i < len(sst.metaBlock); i++ {
		content = sst.metaBlock[i].encode()
		sst.writer.write(content)
		offset += len(content)
	}
	sst.metaindexBlock.set(metaBlockOffset, offset-metaBlockOffset, len(sst.metaBlock))

	// 向文件中写入 metaindexBlock
	content = sst.metaindexBlock.encode()
	sst.footer.metaindexBlockIndex = offset
	offset += len(content)
	sst.writer.write(content)

	// 向文件中写入 indexBlock
	content = sst.indexBlock.encode()
	sst.footer.indexBlockIndex = offset
	offset += len(content)
	sst.writer.write(content)

	// 向文件中写入 footer
	content = sst.footer.encode()
	sst.writer.write(content)

	// 重命名文件, 并且将文件关闭
	sst.close()
}

// sstable 的写类
type sstWriter struct {
	file *os.File
}

// 将内存写入到磁盘
func (w *sstWriter) write(content []byte) {
	w.file.Write(content)
}
