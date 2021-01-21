/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/21
 * @Desc: 保存一些全局的常量, 后续改为配置
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package TinyBase

const (
	// memStore所占内存的阈值, 如果到达了该阈值则将其持久化. 暂定 1024B = 1KB
	maxMemSize int = 1 << 15
	deleted string = "__deleted__"
)
