package main

import (
	"time"

	"github.com/516108736/QKCMonitor/wechatClient"

	"github.com/516108736/QKCMonitor/common"
)

type Instance struct {
	config    common.Config
	checker   ModuleInterface
	weChatClt *wechatClient.WechatClient
}

func NewInstance(config common.Config) (*Instance, error) {
	r := &Instance{
		config: config,
	}
	r.weChatClt = wechatClient.NewWeChatClient(config.WeChatDetail)
	var err error
	r.checker, err = NewRuler(config)
	return r, err
}

func (r *Instance) loop() {
	checkTicker := time.NewTicker(time.Duration(r.config.Interval) * time.Second)
	defer checkTicker.Stop()

	//summaryTicker := time.NewTicker(12 * 60 * time.Minute)
	summaryTicker := time.NewTicker(60 * time.Second)
	defer summaryTicker.Stop()
	for {
		select {
		case <-checkTicker.C:
			errList := r.checker.Check()
			for _, err := range errList {
				r.weChatClt.SendMsg(r.makeErrMsg(err))
			}
		case <-summaryTicker.C:
			errList := r.checker.Summary()
			r.weChatClt.SendMsg(r.makeSummaryMsg(errList))

		}

	}
}

func (r *Instance) makeErrMsg(err string) string {
	return "故障时间:" + time.Now().Format("2006-01-02 15:04:05") + "     故障级别:" + r.config.AlertLevel + "\n出错模型:" + r.config.Module + "\n" + err
}

func (r *Instance) makeSummaryMsg(d []string) string {
	ans := "当前各节点状态展示,没有故障!!!!!!!!!  时间" + time.Now().Format("2006-01-02 15:04:05") + "\n"
	for _, v := range d {
		ans = ans + v
		ans = ans + "\n"
	}
	return ans
}

func main() {
	config := common.LoadConfig("./config.json")
	for _, v := range config.RulerList {
		r, err := NewInstance(v)
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
