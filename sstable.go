/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/2/20
 * @Desc: sstable 的接口.
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

/**
定义 sstable 的格式.

每个 sstable文件以后缀 .sst 表示, sstable的单位为块 block


 */


// sstable, 每个sstable结构对应一个sstable文件.
// sstable文件根据传入的内存memtable将
type SSTable struct {
	name string
	writer SstWriter
	
}


// sstable 的写类
type SstWriter struct {

}


// sstable 的读类
type SstReader struct {

}