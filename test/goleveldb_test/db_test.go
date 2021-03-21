/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/3/3
 * @Desc: goleveldb_test
 * @Copyright (c) 2021, JameyWoo. All rights reserved.
 */

package goleveldb_test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	//"math/rand"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)
import "github.com/syndtr/goleveldb/leveldb"

func TestGoLeveldb(t *testing.T) {
	db, err := leveldb.OpenFile("./db0", nil)
	if err != nil {
		panic(err)
	}
	db.Put([]byte("hello"), []byte("world"), nil)
	value := "hello_"
	for i := 0; i < 4; i++ {
		value += value
	}
	logrus.Info(value)
	defer db.Close()
	for k := 0; k < 1000; k++ {
		go func(k int) {
			for i := 0; i < 10000000; i++ {
				//x := rand.Int() % 100000000
				x := i
				db.Put([]byte(fmt.Sprintf("%04d", k)+fmt.Sprintf("%08d", x)), []byte(value), nil)
				if i%1000000 == 0 {
					logrus.Info(fmt.Sprintf("%04d", k) + "_" + fmt.Sprintf("%08d", x))
				}
			}
		}(k)
	}
	panic(http.ListenAndServe(":8080", nil))
	time.Sleep(1 * time.Second)
}
