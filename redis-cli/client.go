/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/26
 * @Desc: redis_cli
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

/*
一个连接redis服务器的客户端. 为了方便测试
*/

func main() {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	// client.FlushAll()
	client.Set("hello", "world", 0)
	addNKV(client, 0, 3000)
}

// 向redis中添加N个key-value
func addNKV(client *redis.Client, start, end int) {
	for i := start; i <= end; i++ {
		rand.Seed(time.Now().UnixNano())
		// value := "01"
		// for i := 0; i <= 5+rand.Int()%6; i++ {
		// 	value += value
		// }
		value := "00000000000000000000011111111111111111111111"
		client.Set("key_"+fmt.Sprintf("%08d", i), "v_"+value, 0)
		if i%1000 == 0 {
			fmt.Printf("number: %d\n", i)
		}
	}
}
