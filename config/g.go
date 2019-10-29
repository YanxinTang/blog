package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigStruct struct {
	Site struct {
		Name string `json:"name"`
	}
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
	Mysql struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mysql"`
}

var Config ConfigStruct

func init() {
	data, err := ioutil.ReadFile("./config/config.local.json")
	if err != nil {
		log.Fatal("Load config error %v:", err)
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		log.Fatal("Read config error %v:", err)
	}
}
