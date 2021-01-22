# gokv

A persistent, LSM tree structured key value database engine implemented by go language.

## Dev log

## v0.1.0

### target

先实现框架, 接口. 对每个接口的具体实现不做要求, 后续再改进. 

### todo
- [x] 接口: Get, Put, Delete, Scan
- [x] 布隆过滤器算法版 (暂不集成)
- [x] map版 Memstore
- [x] 基本SSTable结构 (追加, 无序, 递增文件序号)
- [x] LSM-Tree结构
- [x] 实现WAL存储 (暂不实现故障恢复功能)


## v0.2.0

### target 

针对某些点进行改善, 实现更多的功能和完整性

### todo

- [x] 读取一个配置文件
- [x] 有序SSTable结构
- [x] 使用标准使用接口(如参考leveldb)改造项目调用方式
- [x] 实现完整的key-value结构, 包括时间戳
- [x] 实现基于WAL的故障恢复(启动时从内存恢复)
- [x] 设计实现varint可变长度编码
- [x] SkipList版 Memstore


## 全部待实现的feature

- [x] 接口: Get, Put, Delete, Scan
- [x] 布隆过滤器算法版 (暂不集成)
- [x] map版 Memstore
- [x] 基本SSTable结构 (追加, 无序, 递增文件序号)
- [x] 基本LSM-Tree结构
- [x] 实现WAL存储 (暂不实现故障恢复功能)
- [x] 实现完整的key-value结构, 包括时间戳
- [x] 读取一个配置文件
- [x] SkipList版 Memstore
- [ ] 支持并发的 SkipList
- [x] 有序SSTable结构
- [ ] 集成Bloom Filter
- [x] 实现基于WAL的故障恢复(内存)
- [ ] 完整SSTable结构(带索引, 暂不分块)
- [ ] 分块的SSTable结构
- [ ] 内存的blockCache缓存
- [ ] 多种WAL策略实现, 可选与默认
- [ ] Compaction基本实现(按时间线与文件大小, hbase的基本策略)
- [ ] Compaction多种策略的实现
- [ ] SSTable的前缀压缩
- [x] 设计实现varint可变长度编码
- [ ] 考虑增量设计(偏移增量等)
- [ ] 考虑键值分离的结构, 类似boltDB
- [x] 使用标准使用接口(如参考leveldb)改造项目调用方式
- [ ] 实现可变长度整形 varint 从而更好地压缩 (binary中有实现)
- [ ] 实现字符串压缩算法
- [ ] 分块读文件, 分块处理