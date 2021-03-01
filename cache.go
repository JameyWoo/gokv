/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/1
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

/*
缓存分类
1. 文件缓存, 通过对文件名到 *os.File 的hash得到. 得到了缓存之后/或者直接open文件之后 存储到sstReader中
	文件缓存会缓存什么?
	1) 一个 *os.File结构
	2) 一个 footer 对象 (因为比较小)
	3) 一个 metaindexBlock 结构
    或者这几个结构可以综合一下, 组成一个新的结构. 这样的文件缓存可以减少大量的对文件的读取次数
2. 块缓存, 通过 文件名+偏移 到 对应的内存结构

所以, 我建立多种缓存. 分别为
1. 对 打开文件及其元数据 的缓存
2. 对 meta block 的缓存 (即 bloom filter)
3. 对 data block 的缓存

这三种缓存分别由三个缓存对象来实现, 在不同的地方调用.
缓存不能无限增长, 因此对每个结构我都需要限制他们的最大的大小. 例如, 默认限制大小为 8MB
*/

package gokv
