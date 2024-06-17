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
	} `yaml:"init"`
	Log struct {
		Path   string `yaml:"path"`
		Level  string `yaml:"level"`
		Log_nu string `yaml:"log_Nu"`
	} `yaml:"log"`
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
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("无法解析YAML文件：%v", err)
	}
	ConfMap = make(map[string]interface{})
	ConfMap["Init.Port"] = config.Init.Port
	ConfMap["log_Nu"] = config.Log.Log_nu
	ConfMap["save_dir"] = config.Init.SavePath
	// 打印配置项的值
	confjson, _ := json.Marshal(ConfMap)
	fmt.Println("conf:", string(confjson))
}
