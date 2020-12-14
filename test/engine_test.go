package test

import (
	"github.com/Jameywoo/TinyBase"
	"github.com/sirupsen/logrus"
	"testing"
	"unsafe"
)

/*
测试引擎基本接口: get, put, delete, scan
 */

func TestEnginePut(t *testing.T) {
	e := TinyBase.NewEngine()
	_ = e.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})
	logrus.Info(e)
}

func TestEngineGet(t *testing.T) {
	e := TinyBase.NewEngine()
	_ = e.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})

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
	e := TinyBase.NewEngine()
	_ = e.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})
	_ = e.Put(TinyBase.KeyValue{Key: "fuck", Value: "you"})
	_ = e.Put(TinyBase.KeyValue{Key: "do", Value: "it"})
	_ = e.Put(TinyBase.KeyValue{Key: "left", Value: "right"})
	_ = e.Put(TinyBase.KeyValue{Key: "shutdown", Value: "away"})

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