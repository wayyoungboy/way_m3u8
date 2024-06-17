package conf

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Init struct {
		Port     int    `yaml:"port"`
		SavePath string `yaml:"save_dir"`
		WorkMax  int    `yaml:"work_max"`
	} `yaml:"init"`
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
		LogNu string `yaml:"log_Nu"`
	} `yaml:"log"`
}

func NewConfig() *Config {
	c := new(Config)
	c.Init.Port = 2045
	c.Init.SavePath = "../../m3u8"
	c.Init.WorkMax = 1
	c.Log.Level = "debug"
	c.Log.Path = "./log"
	c.Log.LogNu = "10"
	return c
}

var ConfMap map[string]interface{}

func ConfInit() {
	// 读取YAML配置文件内容
	yamlFile, err := os.ReadFile("./conf.yaml")
	if err != nil {
		log.Fatalf("无法读取YAML文件：%v", err)
		return
	}

	// 解析YAML配置文件
	config := NewConfig()
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("无法解析YAML文件：%v", err)
	}
	ConfMap = make(map[string]interface{})
	ConfMap["Init.Port"] = config.Init.Port
	ConfMap["log_Nu"] = config.Log.LogNu
	ConfMap["save_dir"] = config.Init.SavePath
	ConfMap["work_max"] = config.Init.WorkMax
	// 打印配置项的值
	confjson, _ := json.Marshal(ConfMap)
	fmt.Println("conf:", string(confjson))
}
