package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type WeChatDetail struct {
	AppID      string `json:"AppID"`
	AppSecret  string `json:"AppSecret"`
	TemplateId string `json:"TemplateId"`
}
type Config struct {
	Module       string       `json:"Module"`
	TargetHosts  []string     `json:"Hosts"`
	Interval     int64        `json:"Interval"`
	WeChatDetail WeChatDetail `json:"WeChatDetail"`
	AlertLevel   string       `json:"AlertLevel"`
	Email        []string     `json:"Email"` //TODO
	Params       interface{}  `json:"Extra"`
}

type LocalConfig struct {
	RulerList []Config
}

func LoadConfig(filePth string) *LocalConfig {
	var config LocalConfig
	f, err := os.Open(filePth)
	if err != nil {
		panic(err)
	}

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
