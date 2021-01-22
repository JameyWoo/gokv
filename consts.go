/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/21
 * @Desc: 保存一些全局的常量, 后续改为配置
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package gokv

const (
	deleted string = "__deleted__"
)

type Op byte

const (
	PUT Op = 0
	DEL
)