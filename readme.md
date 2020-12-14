# Tiny Base

A persistent, LSM tree structured key value database engine implemented by go language.

## Devlopment log

## v0.1.0

### target

先实现框架, 接口. 对每个接口的具体实现不做要求, 后续再改进. 

### todo
- [ ] 接口: Get, Put, Delete, Scan
- [ ] 布隆过滤器算法版 (暂不集成)
- [ ] map版 Memstore
- [ ] 基本SSTable结构 (追加, 无序)
- [ ] LSM结构