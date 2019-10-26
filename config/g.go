package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigStruct struct {
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
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
