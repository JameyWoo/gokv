/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/21
 * @Desc: manifest. 记录了版本信息, 以及sstable的元数据信息
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

/*
一个内存上的数据结构, 记录当前所有的sstable文件信息.
可持久化到磁盘上面

记录每个层级的sstable文件的名称(路径)

有多个level, 默认level为0, 每个level保存了这个level上的文件元数据

文件元数据包括:
1. 该 sstable 的 key范围, 即 最小key和最大key, 因为除了 第0个level, 其他所有同level的文件的区间都是不相交的
2. 该 sstable 的 文件名. 每个文件的文件名都是不同的, 这将他们唯一标识. 从而在使用缓存的时候, 能够通过 文件名+offset 找到对应的缓存对象

在同一个level中(除level 0), 所有的 sstableMeta 按照他们的key来排序
*/
type Manifest struct {
	dir       string          // 保存 sstable文件的目录
	levels    [][]sstableMeta // 二维, 有多层, 每一层都有多个 sstableMeta
	filesizes []int           // 每层的文件size数; level 0不需要维护层数, 从level 1 开始维护
	level     int             // 当前的最大层数
}

// 一个 sstable 的元数据
type sstableMeta struct {
	dir            string // sstable文件的目录
	filename       string // sstable的文件名
	maxKey, minKey string // 最大最小值
	filesize       int    // 文件大小(字节), 每次生成了新的文件之后写入这个值保存
}
