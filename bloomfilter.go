package gokv

import "github.com/howeyc/crc16"

// 所有的哈希函数
var fs = [...]func(data []byte) uint16{crc16.ChecksumCCITT, crc16.ChecksumCCITTFalse,
	crc16.ChecksumIBM, crc16.ChecksumMBus, crc16.ChecksumSCSI}

type BloomFilter struct {
	m     int // 布隆过滤器的长度（如比特数组的大小）
	k     int // 哈希的次数
	h     int // 使用的hash函数的数量, 默认3个
	cnt   int // 已经过滤元素的数量
	array []bool
}

// 默认的构造函数
func NewBloomFilter() *BloomFilter {
	return &BloomFilter{m: 1024, k: 3, h: 3, array: make([]bool, 1024, 1024)}
}

// 带有参数的自定义的构造函数
func NewBloomFilterWithArgs(m, k, h int) *BloomFilter {
	return &BloomFilter{m: m, k: k, h: h, array: make([]bool, m, m)}
}

func (bf *BloomFilter) Put(s string) {
	idx := bf.GetHashIndex(s)

	for _, i := range idx {
		bf.array[i] = true
	}

	bf.cnt += 1

	// 无法扩容... 因为没有保存字符串. 所以需要预先就预测好空间
	// TODO: 处理扩容逻辑, 保证精确度
	//if float32(bf.m) * math.Ln2 * math.Ln2 / float32(bf.cnt) < 0.8 {
	//	bf.expand()
	//}
}

func (bf *BloomFilter) MayContain(s string) bool {
	idx := bf.GetHashIndex(s)

	bit := true
	for _, i := range idx {
		bit = bit && bf.array[i]
	}
	return bit
}

func (bf BloomFilter) GetHashIndex(s string) []int {
	funcs := fs[:bf.h]

	idx := make([]int, bf.h)
	for i, fu := range funcs {
		idx[i] = int(fu([]byte(s))) % bf.m
	}
	//log.Println(idx)
	return idx
}

// 将 [1, 0, 0, 0, 1, 1, 0, 1,    1, 1, 1, 1, 0, 0, 0, 0]
// 编码成 10110001, 00001111
func (bf *BloomFilter) encode() []byte {
	content := make([]byte, 0, len(bf.array)/8)
	for i := 0; i < bf.m; i += 8 {
		var b byte = 0
		for j := i; j < i+8; j++ {
			if bf.array[j] {
				b |= byte(0x1 << j)
			}
		}
		content = append(content, b)
	}
	return content
}

// 解码
func (bf *BloomFilter) decode(content []byte) {

}
