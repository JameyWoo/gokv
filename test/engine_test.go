package test

import (
	"fmt"
	"github.com/Jameywoo/TinyBase"
	"log"
	"testing"
	"unsafe"
)

/*
测试引擎基本接口: get, put, delete, scan
 */

func TestEnginePut(t *testing.T) {
	e := TinyBase.NewEngine()
	_ = e.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})
	fmt.Println(e)
}

func TestEngineGet(t *testing.T) {
	e := TinyBase.NewEngine()
	_ = e.Put(TinyBase.KeyValue{Key: "hello", Value: "world"})

	val, err := e.Get("hello")
	if err != nil {
		panic(err)
	}
	fmt.Printf("val = %s\n", val)

	val, err = e.Get("fuck")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("val = %s\n\n", val)
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
		log.Fatal(err)
	}
	fmt.Println(sr)

	sr, err = e.Scan("donet", "pqu")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sr)

	sr, err = e.Scan("", "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sr)
}

/*
测试memstore的存储容量与阈值
 */
func TestMaxMemSize(t *testing.T) {
	m := make(map[string]string)
	m["hello"] = "world"

	fmt.Println(unsafe.Sizeof(m))

	s := "hello"
	fmt.Println(unsafe.Sizeof(s))
	fmt.Println(len(s))
}