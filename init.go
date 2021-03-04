/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/4
 * @Desc: gokv
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

import "errors"

const (
	deleted string = "__deleted__"
)

// 操作类型及常量
type Op byte

const (
	PUT Op = iota // fix bug: 之前没写 iota, 没有iota不是枚举!
	DEL           // 1
	GET           // 2 在 server 中使用, 被redis请求
)

var GetEmptyError error

// 错误变量初始化
func errorInit() {
	GetEmptyError = errors.New("GetEmptyError: no such element")
}

// 一个统一的初始化函数, 调用其他初始化变量
func init() {
	cacheInit() // 缓存初始化
	logInit()   // 日志初始化
	errorInit() // 错误初始化, 虽然没怎么用到这些错误
}
