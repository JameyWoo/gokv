/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/4/7
 * @Desc: main
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPutGet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushAll()
	client.Set("hello", "world", 0)
	kvs := make(map[string]string)
	addNKVMap(client, 0, 3000, kvs)
	time.Sleep(1 * time.Second)
	for key, val := range kvs {
		res := client.Get(key)
		v, _ := res.Result()
		//logrus.Infof("Info: val: %s, v: %s", val, v)
		if val != v {
			logrus.Warnf("val: %s, v: %s", val, v)
		}
	}
}

// 向redis中添加N个key-value
func addNKVMap(client *redis.Client, start, end int, kvs map[string]string) {
	for i := start; i <= end; i++ {
		rand.Seed(time.Now().UnixNano())
		value := "00000000000000000000011111111111111111111111"
		kvs["key_"+fmt.Sprintf("%08d", i)] = "v_" + value
		client.Set("key_"+fmt.Sprintf("%08d", i), "v_"+value, 0)
		if i%1000 == 0 {
			fmt.Printf("number: %d\n", i)
		}
	}
}

// 对普通 list类型的ziplist的测试
func TestListZiplist(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.RPush("la", "a", "b", "123", "shit", "fuck")
	strs, err := client.LRange("la", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logrus.Info(strs)
}

func TestSetPutListZiplist(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushAll()
	client.Set("hello", "world", 0)
	kvs := make(map[string]string)
	listZiplistNKVMap(client, 0, 10000, kvs)
	time.Sleep(1 * time.Second)
	for key, val := range kvs {
		res, err := client.LRange(key, 0, 0).Result()
		logrus.Infof("key: %s, val: %s", key, res)
		if err != nil {
			logrus.Panic(err)
		}
		//logrus.Infof("Info: val: %s, v: %s", val, v)
		if val != res[0] {
			logrus.Warnf("val: %s, v: %s", val, res[0])
		}
	}
}

func listZiplistNKVMap(client *redis.Client, start, end int, kvs map[string]string) {
	for i := start; i <= end; i++ {
		key := "list_ziplist_" + fmt.Sprintf("%06d", i)
		kvs[key] = "hello, world"
		client.RPush(key, "hello, world")
		if i%1000 == 0 {
			fmt.Printf("number: %d\n", i)
		}
	}
}

func TestHashZiplist(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	client.FlushAll()
	key := "fuckit"
	client.HSet(key, "hello", "world")
	te, _ := client.Type(key).Result()
	logrus.Infof("Type: %s", te)
	res, err := client.HGet(key, "hello").Result()
	if err != nil {
		panic(err)
	}
	logrus.Infof("key: %s, val: %s", key, res)
}

func TestTypeHashZiplist(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	//client.FlushAll()
	start := 0
	end := 0
	for i := start; i <= end; i++ {
		key := "hash_ziplist_" + fmt.Sprintf("%06d", i)
		client.HSet(key, "hello", "world")
		logrus.Infof("key: %s", key)
		res, _ := client.Type(key).Result()
		logrus.Infof("type: %s", res)
	}
}

func TestSetPutHashZiplist(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		// Addr:     "www.firego.cn:6379",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushAll()
	client.Set("hello", "world", 0)
	kvs := make(map[string]string)
	hashZiplistNKVMap(client, 0, 10000, kvs)
	time.Sleep(1 * time.Second)
	for key, val := range kvs {
		res, err := client.HGet(key, "hello").Result()
		logrus.Infof("key: %s, val: %s", key, res)
		if err != nil {
			logrus.Panic(err)
		}
		if val != res {
			logrus.Warnf("val: %s, v: %s", val, res)
		}
	}
}

func hashZiplistNKVMap(client *redis.Client, start, end int, kvs map[string]string) {
	for i := start; i <= end; i++ {
		key := "hash_ziplist_" + fmt.Sprintf("%06d", i)
		kvs[key] = "world"
		client.HSet(key, "hello", "world")
		client.HSet(key, "fuck", "you")
		logrus.Infof("key: %s", key)
		if i%1000 == 0 {
			fmt.Printf("number: %d\n", i)
		}
	}
}

// 测试大量的key被set, 用多个协程
func TestScaleKeySet(t *testing.T) {
	scale := 10000
	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := redis.NewClient(&redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
			})
			addNKV(client, scale*(i-1), scale*i)
		}(i)
	}
	wg.Wait()
	logrus.Info("over!")
}

func TestScaleKeyGet(t *testing.T) {
	scale := 100000
	value := "v_00000000000000000000011111111111111111111111"
	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := redis.NewClient(&redis.Options{
				// Addr:     "www.firego.cn:6379",
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
			})
			for j := scale * (i - 1); j < scale*i; j++ {
				val, err := client.Get("key_" + fmt.Sprintf("%08d", i)).Result()
				if err != nil {
					panic(err)
				}
				if val != value {
					logrus.Infof("key: %s", "key_"+fmt.Sprintf("%08d", i))
					logrus.Warnf("value: %s, val: %s", value, val)
				}
			}
		}(i)
	}
	wg.Wait()
	logrus.Info("over!")
}

// 对redis进行一个性能测试
func TestRedisPerformance(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	start := time.Now()
	for i := 0; i < 100000; i++ {
		client.Set(fmt.Sprintf("k_%d", i), "value", 0)
	}
	end := time.Now()
	logrus.Info(end.Sub(start))
}
