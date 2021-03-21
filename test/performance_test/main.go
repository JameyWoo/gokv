/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/20
 * @Desc: performance_test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"fmt"
	"github.com/Jameywoo/gokv"
	"math/rand"
	"strconv"
)

func main() {
	db, err := gokv.Open("../db7", &gokv.Options{ConfigPath: "./gokv.yaml"})
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		//time.Sleep(1 * time.Millisecond)
		_ = db.Put(strconv.Itoa(rand.Int()%10000000), "val")
		i++
		if i%1000 == 0 {
			fmt.Println(i)
		}
	}
}
