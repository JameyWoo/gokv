/**
 * @Author: JameyWoo
 * @Email: 2622075127wjh@gmail.com
 * @Date: 2021/1/21
 * @Desc: gokv
 * @Copyright (c) 2020, JameyWoo. All rights reserved.
 */

package gokv

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var config *conf

func init() {
	var c conf
	config = c.getConf()
	//fmt.Println(config.MaxMemSize)
}

//profile variables
type conf struct {
	MaxMemSize int `yaml:"MaxMemSize"`  // 注意一定要是大写
}

func (c *conf) getConf() *conf {
	// TODO: 诸葛这个配置文件的位置需要根据test或者发布版本进行修改
	yamlFile, err := ioutil.ReadFile("../gokv.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}