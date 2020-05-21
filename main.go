package main

import (
	"time"

	"github.com/516108736/QKCMainnetMonitor/wechatClient"

	"github.com/516108736/QKCMainnetMonitor/common"
)

type Role struct {
	config    common.Config
	checker   RulerI
	WeChatClt *wechatClient.WechatClient
}

func NewRole(config common.Config) (*Role, error) {
	r := &Role{
		config: config,
	}
	r.WeChatClt = wechatClient.NewWeChatClient(config.WeChatDetail)
	var err error
	r.checker, err = NewRuler(config)
	return r, err
}

func (r *Role) loop() {
	ticker := time.NewTicker(time.Duration(r.config.Interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			errList := r.checker.Check()
			for _, err := range errList {
				r.WeChatClt.SendMsg(r.makeErrMsg(err))
			}
		}

	}
}

func (r *Role) makeErrMsg(err string) string {
	return "故障时间:" + time.Now().Format("2006-01-02 15:04:05") + "     故障级别:" + r.config.AlertLevel + "\n出错模型:" + r.config.Module + "\n" + err
}

func main() {
	config := common.LoadConfig("./config.json")
	for _, v := range config.RulerList {
		r, err := NewRole(v)
		if err != nil {
			panic(err)
		}
		if err := r.checker.PreCheck(); err != nil {
			panic(err)
		}
		go r.loop()
	}
	select {}
}
