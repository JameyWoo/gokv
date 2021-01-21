package test

import (
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
	"unsafe"
)

/*
测试引擎基本接口: get, put, delete, scan
 */

func TestEnginePut(t *testing.T) {
	e := gokv.NewEngine()
	_ = e.Put(gokv.KeyValue{Key: "hello", Value: "world"})
	logrus.Info(e)
}

func TestEngineGet(t *testing.T) {
	e := gokv.NewEngine()
	_ = e.Put(gokv.KeyValue{Key: "hello", Value: "world"})

	val, err := e.Get("hello")
	if err != nil {
		panic(err)
	}
	logrus.Infof("val = %s", val)

	val, err = e.Get("fuck")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("val = %s", val)
}

func TestEngineDelete(t *testing.T) {
	e := gokv.NewEngine()
	_ = e.Put(gokv.KeyValue{Key: "hello", Value: "world"})
	_ = e.Put(gokv.KeyValue{Key: "fuck", Value: "you"})
	_ = e.Put(gokv.KeyValue{Key: "do", Value: "it"})
	_ = e.Put(gokv.KeyValue{Key: "left", Value: "right"})
	_ = e.Put(gokv.KeyValue{Key: "shutdown", Value: "away"})

	sr, err := e.Scan("", "z")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(sr)

	sr, err = e.Scan("donet", "pqu")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(sr)

	sr, err = e.Scan("", "foo")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(sr)
}

/*
测试memstore的存储容量与阈值
 */
func TestMaxMemSize(t *testing.T) {
	m := make(map[string]string)
	m["hello"] = "world"

	logrus.Info(unsafe.Sizeof(m))

	s := "hello"
	logrus.Info(unsafe.Sizeof(s))
	logrus.Info(len(s))
}

// 测试可变长度varint的编码和解码
func TestVarInt(t *testing.T) {
	//bytes := gokv.VarIntEncode(666)
	// 127, 128, 0, 12345678910
	bytes := gokv.VarIntEncode(12345678910)
	logrus.Info(len(bytes))

	logrus.Info(bytes)

	value, b := gokv.VarIntDecode(bytes)
	logrus.Info("value: ", value)
	logrus.Info("b: ", b)
}

// 测试 KeyValue 的编码和解码
func TestKeyValue(t *testing.T) {
	kv := gokv.KeyValue{
		Key: "hello",
		Value: "world",
		Timestamp: time.Now().UnixNano() / 1e6,
		Op: gokv.DEL,
	}
	logrus.Info(kv)
	bytes := kv.Encode()
	nkv, bytes := gokv.KvDecode(bytes)
	logrus.Info(nkv)
}