/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2020/12/15
 * @Desc: gokv
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package gokv

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
	"time"
)

type DB struct {
	memDB    *MemDB
	wal      *os.File
	options  *Options
	manifest *Manifest

	dir     string
	walPath string
}

func (db *DB) Dir() string {
	return db.dir
}

func (db *DB) MemDB() *MemDB {
	return db.memDB
}

func Open(dirPath string, options *Options) (*DB, error) {
	// 生成一个文件夹
	// 包含 WAL, HFile目录
	if !Exists(dirPath) {
		os.Mkdir(dirPath, os.ModePerm)
	}
	var wal *os.File
	walPath := dirPath + "/wal.log"
	// 查看是否有日志文件, 如果没有则创建
	checkWal(walPath)
	wa, err := os.OpenFile(walPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	wal = wa

	// 在打开数据库的时候读取配置
	readConfig(options.ConfigPath)

	mani := &Manifest{}
	mani.dir = dirPath
	db := &DB{
		memDB:    NewEngine(),
		wal:      wal,
		dir:      dirPath,
		walPath:  walPath,
		options:  options,
		manifest: mani,
	}
	err = db.recoveryDB()
	if err != nil {
		logrus.Error(err)
	}

	return db, err
}

// 从预写日志中恢复数据, 如果wal中存在数据
func (db *DB) recoveryDB() error {
	bytes, err := ioutil.ReadFile(db.walPath)
	if err != nil {
		logrus.Error(err)
	}
	for len(bytes) != 0 {
		kv := KeyValue{}
		kv, bytes = KvDecode(bytes)
		err = db.put(kv)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return err
}

func checkWal(walPath string) {
	if !Exists(walPath) {
		wa, err := os.Create(walPath)
		if err != nil {
			logrus.Fatal(err)
		}
		wa.Close()
	}
}

// panic 的时候还是会执行 defer Close.
func (db *DB) Close() {
	//os.Exit(0)
	// 先落盘, 再清空预写日志
	// 刷到磁盘
	db.flush()
	db.wal.Close()
	// 清空 wal.log 文件内容, 用直接创建的方式
	// Create creates or truncates the named file. If the file already exists, it is truncated.
	wal, err := os.Create(db.walPath)
	if err != nil {
		logrus.Fatal(err)
	}
	wal.Close()
}

func (db *DB) Get(key string) (Value, bool) {
	value, err := db.memDB.Get(key)
	// 如果得到value
	if err == nil {
		return value, true
	}
	// 从磁盘上获取value
	return db.diskGet(key)
}

func (db *DB) Put(key, value string) error {
	// 时间戳, 以ms为单位. 但是如果在同一ms内对同一个值的不同操作的话, 应该怎么办?
	return db.put(KeyValue{Key: key,
		Val: Value{Value: value, Timestamp: time.Now().UnixNano() / 1e6, Op: PUT}})
}

func (db *DB) put(kv KeyValue) error {
	// TODO: 改进写日志的方式
	db.writeAheadLog(kv)
	db.memDB.Put(kv)
	// 这个阈值用常量 MaxMemSize表示, MaxMemSize定义在配置文件中, 后续改为可配置的量
	if db.memDB.memSize >= config.MaxMemSize {
		// 刷到磁盘
		db.flush()
		// TODO: 之后实现异步的flush操作
		db.memDB.memStore = NewSkipList()
		// ! 这个也要重置
		db.memDB.memSize = 0
		// 清空 wal.log 文件内容, 用直接创建的方式
		// Create creates or truncates the named file. If the file already exists, it is truncated.
		wal, err := os.Create(db.walPath)
		if err != nil {
			logrus.Fatal(err)
		}
		wal.Close()
		wal, err = os.OpenFile(db.walPath, os.O_WRONLY|os.O_APPEND, 0666)
		db.wal = wal
	}
	return nil
}

// 删除的元素的value用特殊的字符串来代替
func (db *DB) Delete(key string) error {
	delTime := time.Now().UnixNano() / 1e6
	db.writeAheadLog(KeyValue{Key: key, Val: Value{"", delTime, DEL}})
	return db.memDB.Delete(key, delTime)
}

// 暂时不需要区间扫描
// 扫描一个区间的key, 得到key value的结果slice
// 如果value为deleted, 那么不添加
//func (db *DB) Scan(startKey, endKey string) ([]KeyValue, error) {
//	return db.memDB.Scan(startKey, endKey)
//}

func (db *DB) writeAheadLog(kv KeyValue) error {
	write := bufio.NewWriter(db.wal)
	_, err := write.Write(kv.Encode())
	if err != nil {
		logrus.Error(err)
	}
	_ = write.Flush()
	return err
}

// 使用 sstable 实现的 flush
func (db *DB) flush() error {
	filename := GetTimeString() + ".sst"
	minKey := db.memDB.getMinKey()
	maxKey := db.memDB.getMaxKey()
	sst := NewSSTable(db.Dir(), filename, db.MemDB())
	filesize := sst.Write()

	// 初始的时候, level = 0, 给 levels设置初始值
	if db.manifest.level == 0 {
		db.manifest.levels = append(db.manifest.levels, make([]sstableMeta, 0))
		db.manifest.filesizes = append(db.manifest.filesizes, 0)
		db.manifest.level++
	}
	nsm := sstableMeta{
		dir:      db.Dir(),
		filename: filename,
		minKey:   minKey,
		maxKey:   maxKey,
		filesize: filesize,
	}
	// 将生成的 sstable 加入到 db.manifest 中
	db.manifest.levels[0] = append(db.manifest.levels[0], nsm)
	// 每次 flush 的时候就会触发 compaction. 检查 compaction条件
	db.judgeCompact()
	return nil
}

func (db *DB) diskGet(key string) (Value, bool) {
	// 根据 manifest的信息, 先从第0层依次读取sstable
	// 如果没有, 那么从其他层找
	if db.manifest.level == 0 { // 一个 sstable文件都没有
		return Value{}, false
	}
	v, find := db.findOnOneLevel(0, key)
	if find {
		return v, find
	}
	// 对剩下的每一层都检察一下
	for i := 1; i < db.manifest.level; i++ {
		v, find = db.findOnOneLevel(i, key)
		if find {
			return v, find
		}
	}
	return Value{}, false
}

func (db *DB) findOnOneLevel(level int, key string) (Value, bool) {
	for i := 0; i < len(db.manifest.levels[level]); i++ {
		if key <= db.manifest.levels[level][i].maxKey && key >= db.manifest.levels[level][i].minKey {
			ssr := sstReader{}
			defer ssr.close()
			ssr.open(db.dir + "/" + db.manifest.levels[level][i].filename)
			pv, find := ssr.FindKey(key)
			if find {
				return *pv, true
			}
		}
	}
	return Value{}, false
}

func (db *DB) MemIterator() *Iterator {
	return db.memDB.NewIterator()
}

// TODO: minor compaction 的重要性大于 major compaction, 所以当minor compaction需要执行的时候, 应当有机制能够暂停major compaction的执行
/*
判断是否进行compaction, 如果满足条件, 那么进行合并

有两种条件
1. level 0 的sstable文件数量超过4个
2. level 1 ~ n (设为i) 的sstable文件的总和超过该层限额 (10 ^ i) MB；

PS. 每次第二种条件都是在发生了第一种条件之后触发的

每一层 level的sstable文件的大小总和限制是程对数增长的. 第 i 层的文件综合不超过 (10 ^ i) MB；
1. 10MB
2. 100MB
3. 1000MB
4. 10000MB
5. ...

对于 > 0 层的sstable文件而言, 每个文件的区间都是没有交集的.
因此, 当第0层的文件

每次合并哪些文件? 见 https://leveldb-handbook.readthedocs.io/zh/latest/compaction.html
红星标注的为起始输入文件；
在level i层中，查找与起始输入文件有key重叠的文件，如图中红线所标注，最终构成level i层的输入文件；
利用level i层的输入文件，在level i+1层找寻有key重叠的文件，结果为绿线标注的文件，构成level i，i+1层的输入文件；
最后利用两层的输入文件，在不扩大level i+1输入文件的前提下，查找level i层的有key重叠的文件，结果为蓝线标准的文件，构成最终的输入文件；
*/
func (db *DB) judgeCompact() {
	cur := 0
	var more bool
	for {
		more = majorCompaction(db, cur)
		if !more {
			break
		}
		cur++
	}
}

// major compaction, 对 > 1 层的level进行合并; 是一个递归函数
func majorCompaction(db *DB, cur int) bool {
	zeroLen := len(db.manifest.levels[cur])
	nxt := cur + 1
	// TODO: 这个值 (4) 应当可配置
	if cur == 0 && zeroLen <= 4 { // 不需要进行压缩
		return false
	}
	// 先处理从0层到第1层的情况
	minKey := db.manifest.levels[cur][zeroLen-1].minKey
	maxKey := db.manifest.levels[cur][zeroLen-1].maxKey
	preToCompact := make([]sstableMeta, 0)
	// 新加入的 sstableMeta 一定是最后一个
	// 选择跟最新的sstableMeta有区间重合的
	for i := 0; i < zeroLen-1; i++ {
		// 排除两种非重合的情况
		if !((minKey >= db.manifest.levels[cur][i].maxKey) || (maxKey <= db.manifest.levels[cur][i].minKey)) {
			preToCompact = append(preToCompact, db.manifest.levels[cur][i])
		}
	}
	preToCompact = append(preToCompact, db.manifest.levels[cur][zeroLen-1])
	// 遍历第0层需要合并的sstable, 找到key的区间
	for i := 0; i < len(preToCompact); i++ {
		if preToCompact[i].minKey < minKey {
			minKey = preToCompact[i].minKey
		}
		if preToCompact[i].maxKey > maxKey {
			maxKey = preToCompact[i].maxKey
		}
	}
	// 找到第1层跟key重合的列表. 需要考虑当level 1 (level nxt)还没有文件的情况
	if db.manifest.level == nxt {
		// 需要扩展出一层
		db.manifest.levels = append(db.manifest.levels, make([]sstableMeta, 0))
		db.manifest.filesizes = append(db.manifest.filesizes, 0)
		db.manifest.level++
	} else {
		// 将 level 1 要合并的文件加入到 preToCompact 中
		for i := 0; i < len(db.manifest.levels[nxt]); i++ {
			if !((minKey >= db.manifest.levels[nxt][i].maxKey) || (maxKey <= db.manifest.levels[nxt][i].minKey)) {
				preToCompact = append(preToCompact, db.manifest.levels[nxt][i])
			}
		}
	}
	sstm := compact(preToCompact)
	sstm.minKey = minKey
	sstm.maxKey = maxKey

	// 得到了新的文件之后需要将其他的sstableMeta删除
	// l0 保存新的 sstableMeta结构
	l0 := make([]sstableMeta, 0)
	db.manifest.filesizes[cur] = 0
	for i := 0; i < len(db.manifest.levels[cur]); i++ {
		in := false
		for j := 0; j < len(preToCompact); j++ {
			if db.manifest.levels[cur][i] == preToCompact[j] {
				in = true
				break
			}
		}
		if !in {
			l0 = append(l0, db.manifest.levels[cur][i])
		}
	}
	db.manifest.levels[cur] = l0

	// 删除 l1的
	l1 := make([]sstableMeta, 0)
	db.manifest.filesizes[nxt] = 0
	for i := 0; i < len(db.manifest.levels[nxt]); i++ {
		in := false
		for j := 0; j < len(preToCompact); j++ {
			if db.manifest.levels[nxt][i] == preToCompact[j] {
				in = true
			}
		}
		if !in {
			// 累计filesize. filesize要先置零
			db.manifest.filesizes[nxt] += db.manifest.levels[nxt][i].filesize
			l1 = append(l1, db.manifest.levels[nxt][i])
		}
	}
	db.manifest.levels[nxt] = l1

	// 删除 sstableMeta 里的文件
	for i := 0; i < len(preToCompact); i++ {
		// ! fix bug: 想要删除文件但是发现删除不了. 说是文件被使用. 说明我使用这些文件compaction的时候并没有即使将他们close
		//logrus.Info(preToCompact[i].dir + "/" + preToCompact[i].filename)
		err := os.Remove(preToCompact[i].dir + "/" + preToCompact[i].filename)
		if err != nil {
			panic(err)
		}
	}

	db.manifest.filesizes[nxt] += sstm.filesize
	// 将新的 sstableMeta添加进去
	db.manifest.levels[nxt] = append(db.manifest.levels[nxt], *sstm)
	threshold := int(math.Pow10(nxt-1)) * 1024 * 1024
	//logrus.Info(db.manifest.filesizes[nxt])
	if db.manifest.filesizes[nxt] > threshold {
		return true
	}
	return false
}
