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

记录每个层级的sstable文件的名称
*/
type Manifest struct {
}
