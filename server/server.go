/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/3
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

/*
提供一种服务
*/

package main

import (
	"fmt"
	"github.com/Jameywoo/gokv"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
)

func main() {
	//open()
	// 写一个web服务, 监听端口, 然后不断地读取服务
	port := "5379"
	if len(os.Args) > 2 {
		fmt.Printf("Usage : %s <port>\n", os.Args[0])
		os.Exit(1)
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	}
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// 处理一个连接
func handleConn(conn net.Conn) {
	defer conn.Close()
	toClient := make(chan []byte) // outgoing client messages
	go clientWriter(conn, toClient)
	// 当连接第一次打开的时候, 会传递数据库的目录进来供 DB.Open() 使用
	dir, err := ReceiveBytesFromConn(conn)
	if err != nil {
		logrus.Error(err)
		return
	}
	db, err := gokv.Open(string(dir), &gokv.Options{ConfigPath: "./gokv.yaml"})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("conn")
	toClient <- []byte(dir)
	for {
		// 从客户端读取数据
		fromClient, err := ReceiveBytesFromConn(conn)
		if err != nil { // 如果读取出现问题, 那么直接退出
			logrus.Error(err)
			return
		}
		//logrus.Info(string(fromClient))
		//toClient <- fromClient
		// redis需要有能够将各种数据结构序列化成字节数组的函数. 并且有反序列化的方法
		// server 需要支持三种命令, GET SET DEL. 使用switch对三种情况分别处理
		msg := Message{}
		// 解析消息
		msg.parse(fromClient)
		// 判断消息类型并执行
		switch msg.op {
		case gokv.PUT:
			err := db.Put(string(msg.key), string(msg.value))
			// 应当要有错误处理通知客户端key写失败了, 需要重新写
			if err != nil {
				toClient <- []byte("PUT_ERROR")
			} else {
				toClient <- []byte("PUT_SUCCESS")
			}
		case gokv.DEL:
			err := db.Delete(string(msg.key))
			// 应当要有错误处理通知客户端key删除失败了, 需要重新删除
			if err != nil {
				toClient <- []byte("DEL_ERROR")
			} else {
				toClient <- []byte("DEL_SUCCESS")
			}
		case gokv.GET:
			value, get := db.Get(string(msg.key))
			// 成功get则向客户端发送value. 失败则通知失败
			if get {
				toClient <- []byte(value.Value)
				_ = value
			} else {
				toClient <- []byte("GET_ERROR")
			}
		}
	}
}

// 向客户端写入数据的协程
// ! ch循环阻塞, 知道ch被传入了数据, 这个协程不会直接终止
func clientWriter(conn net.Conn, ch <-chan []byte) {
	// * range 遍历, 当 ch 为空的时候, 这个语句会阻塞. 当 ch 得到了值, 他又会醒过来
	for msg := range ch {
		msgByte := msg
		msgByteNew := BytesCombine(IntToBytes(len(msgByte)), msgByte)
		conn.Write(msgByteNew)
	}
}
