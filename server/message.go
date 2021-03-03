/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/3
 * @Desc: main
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"encoding/binary"
	"github.com/Jameywoo/gokv"
)

// 从客户端传来的消息的格式
type Message struct {
	op    gokv.Op
	key   []byte
	value []byte
}

/*
将 []byte 的消息解码位 Message格式
key 和 value 的长度用 uint32编码, 4字节
编码为 [op(1), key_len(4), key, value_len(4), value]
*/
func (m *Message) parse(msg []byte) {
	m.op = gokv.Op(msg[0])
	klen := binary.LittleEndian.Uint32(msg[1:5])
	m.key = msg[5 : 5+klen]
	// 只有 PUT 的时候才需要获得 value, 其他时候只有 key
	if m.op == gokv.PUT {
		vlen := binary.LittleEndian.Uint32(msg[5+klen : 9+klen])
		m.value = msg[9+klen : 9+klen+vlen]
	}
}
