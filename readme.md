# Tiny Base

A persistent, LSM tree structured key value database engine implemented by go language.

## Dev log

## v0.1.0

### target

先实现框架, 接口. 对每个接口的具体实现不做要求, 后续再改进. 

### todo
- [x] 接口: Get, Put, Delete, Scan
- [x] 布隆过滤器算法版 (暂不集成)
- [x] map版 Memstore
- [ ] 基本SSTable结构 (追加, 无序, 递增文件)
- [ ] LSM结构
- [ ] 实现key-value接口, 包括时间戳
- [ ] 实现WAL存储 (暂不实现故障恢复功能)


## 其他待实现的feature

- [ ] SkipList版 Memstore
- [ ] 支持并发的 SkipList
- [ ] 有序SSTable结构
- [ ] 集成Bloom Filter
- [ ] 实现基于WAL的故障恢复
- [ ] 完整SSTable结构(带索引, 暂不分块)
- [ ] 分块的SSTable结构
- [ ] 内存的blockCache缓存
- [ ] 多种WAL策略实现, 可选与默认
- [ ] Compaction基本实现(按时间线与文件大小, hbase的基本策略)
- [ ] Compaction多种策略的实现
- [ ] SSTable的前缀压缩