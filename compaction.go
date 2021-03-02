/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/26
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

/*
// ! 数据压缩, 就是将两个(或更多)sstable合并成一个
策略? 根据 https://leveldb-handbook.readthedocs.io/zh/latest/compaction.html 的压缩策略

接口:
	输入: 两(多)个 sstableMeta (传递一个 []sstableMeta进来, 合并的时候遍历这个结构)
	输出: 下一level的一个sstableMeta结构
	过程: 读取 sstableMeta中的文件内容, 使用合并排序的方式将两(多)个文件合并成一个

注意:
1. 待合并的文件可能有多个
2. 如何实现在合并的过程中只有部分内容在内存上, 大部分内容都不在内存上.
	每次合并的时候, 取每个文件的最新的datablock, 按顺序将全局最小的 key-value加入到新的datablock中
	在这个过程中, 内存上只需要维持 k + 1 个datablock, 因此占用内存较少, 是可以合理实现的.

特点:
1. 合并过程中, 主要关注 datablock以及indexblock(用来索引datablock), 其他如 metablock, metaindexblock等都不需要了
*/

package gokv

import (
	"os"
	"strconv"
	"time"
)

// 压缩, 将几个 sstable 压缩成一个
func compact(sstMetas []sstableMeta) *sstableMeta {
	n := len(sstMetas)
	iters := make([]*sstableIter, 0, n)
	for _, sm := range sstMetas {
		file, err := os.Open(sm.dir + "/" + sm.filename)
		if err != nil {
			panic(err)
		}
		iters = append(iters, newSSTableIter(file))
	}
	// 应该实现一个 min heap, 每次得到最小的keyvalue
	mh := newMinHeap()
	for i := 0; i < n; i++ {
		kv, have := iters[i].Next()
		if have { // 每次都需要判断是否
			mh.push(kv, i)
		}
	}
	// 新的 nsm
	nsm := sstableMeta{dir: sstMetas[0].dir, filename: strconv.FormatInt(time.Now().UnixNano(), 10) + ".sst"}
	// 创建一个sstable文件
	sst := NewSSTable(nsm.dir, nsm.filename, nil)
	sst.open()
	// 类似 sstable.Write() 实现的内容. 这里本来应该写成一个接口的, 因为代码很多都是重复的
	metaB := newMetaBlock(2048) // 初始化布隆过滤器, 使用构造函数
	offset := 0                 // 全局的偏移
	var content []byte
	globalCount := 0 // 全局的count, 表示当前的 datablock index 的key的数量
	isFirstKey := true
	var lastKey KeyValue

	for {
		hi, more := mh.getMin() // 一个元素都没有了
		if !more {
			break
		}
		mh.pop()
		kvNext, have := iters[hi.index].Next()
		if have {
			mh.push(kvNext, hi.index)
		}
		// TODO: 需要增加一个缓冲, 从而当key相同的时候能够特殊处理(只保留一个key)
		if isFirstKey {
			isFirstKey = false
			lastKey = hi.kv
			// 第一个key的时候不添加, 设置为 lastKey然后跳过
			continue
		}
		if hi.kv.Key == lastKey.Key {
			// 如果当前 key 跟 上一个 key相同, 那么 赋值然后跳过
			lastKey = hi.kv
			continue
		}
		kv := lastKey
		lastKey = hi.kv
		//logrus.Info(kv.Key, kv.Val.Op)
		// 处理 kv, 将 kv 添加到新的 sstable中. 内容类似 sstable.Write()
		sstAddKeyValue(sst, metaB, content, offset, globalCount, kv)
	}
	// 剩下一个lastKey需要添加
	sstAddKeyValue(sst, metaB, content, offset, globalCount, lastKey)
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

	return &nsm
}

func sstAddKeyValue(sst *SSTable, metaB *metaBlock, content []byte, offset, globalCount int, kv KeyValue) {
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
