package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	DBHost    string `json:"dbHost"`
	DBPort    string `json:"dbPort"`
	DBUser    string `json:"dbUser"`
	DBPass    string `json:"dbPass"`
	DBName    string `json:"dbName"`
	SecretKey string `json:"secretKey"`
	LocalWeb  bool   `json:"localWeb"`
	WebPort   uint   `json:"webPort"`
}

var (
	currentPath string
)

func init() {
	currentPath, _ = os.Getwd()
	if isExists, _ := FileExists(currentPath + "/config.json"); !isExists {
		config := Config{
			DBHost:    "127.0.0.1",
			DBPort:    "27017",
			DBUser:    "admin",
			DBPass:    "admin",
			DBName:    "march",
			SecretKey: "123456",
			LocalWeb:  true,
			WebPort:   52080,
		}
		SaveConfig(config)
		log.Fatalln("未检测到配置文件，已写入默认配置，请修改后重新运行")
	}
}

func SaveConfig(c Config) {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = ioutil.WriteFile(currentPath+"/config.json", bytes, 0755)
	if err != nil {
		log.Fatalln("配置文件写入失败: " + err.Error())
	}
}

func LoadConfig() Config {
	bytes, err := ioutil.ReadFile(currentPath + "/config.json")
	if err != nil {
		log.Fatalln("配置文件读取失败: " + err.Error())
	}

	var config Config
	json.Unmarshal(bytes, &config)

	return config
}
