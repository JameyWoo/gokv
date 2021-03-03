/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/21
 * @Desc: 保存一些全局的常量, 后续改为配置
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package gokv

const (
	deleted string = "__deleted__"
)

type Op byte

const (
	PUT Op = iota // fix bug: 之前没写 iota, 没有iota不是枚举!
	DEL           // 1
	GET           // 2 在 server 中使用, 被redis请求
)
