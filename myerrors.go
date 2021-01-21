/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/21
 * @Desc: 定义一些错误
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package TinyBase

import "errors"

var GetEmptyError error

func init() {
	GetEmptyError = errors.New("GetEmptyError: no such element")
}
