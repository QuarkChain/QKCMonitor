package blockTime

import (
	"fmt"
	"time"

	"github.com/516108736/QKCMonitor/qkcClient"

	"github.com/516108736/QKCMonitor/common"
)

type BlockTime struct {
	param      param
	QkcClients map[string]*qkcClient.Client
	LastHeight map[string]*blockDetail
}

type blockDetail struct {
	height uint64
	time   int64
}

type param struct {
	Interval int64
	MaxBlock uint64
}

func (b *BlockTime) SetExtraParams(data interface{}) {
	res := param{}
	res.Interval = int64(data.(map[string]interface{})["Interval"].(float64))
	res.MaxBlock = uint64(data.(map[string]interface{})["MaxBlock"].(float64))
	b.param = res
}

func New(config common.Config) (*BlockTime, error) {
	b := &BlockTime{
		QkcClients: make(map[string]*qkcClient.Client),
		LastHeight: make(map[string]*blockDetail),
	}

	for _, host := range config.TargetHosts {
		b.QkcClients[host] = qkcClient.NewClient(host)
		b.LastHeight[host] = &blockDetail{}
	}

	b.SetExtraParams(config.Params)
	return b, nil
}

func (b *BlockTime) makeError(host string, errMsg error) string {
	return fmt.Sprintf("host:%v \nerrMsg:%v", host, errMsg)
}

func (b *BlockTime) Check() []string {
	res := make([]string, 0)
	for host, client := range b.QkcClients {
		ts := time.Now().Unix()
		height, err := client.GetRootBlockHeight()
		if err != nil {
			res = append(res, b.makeError(host, err))
		}
		ls := b.LastHeight[host]
		if ts-ls.time >= b.param.Interval {
			if height-ls.height <= b.param.MaxBlock {
				res = append(res, b.makeError(host, fmt.Errorf(
					"上次检查点%d 上次检查区块高度%d 本次检查点%d 本次检查区块高度%d  时间间隔%d root块增加%d<=%d",
					ls.time, ls.height, ts, height, ts-ls.time, height-ls.height, b.param.MaxBlock)))
			}
		}
		fmt.Println("BlockTime check end", "ip", host, time.Now().Format("2006-01-02 15:04:05"), "lastHeight", ls.height, "latest Height", height)
		b.LastHeight[host].height = height
		b.LastHeight[host].time = ts
	}
	return res
}

func (b *BlockTime) PreCheck() error {
	for host, client := range b.QkcClients {
		_, err := client.GetRootBlockHeight()
		if err != nil {
			return fmt.Errorf("host %v GetRootBlockHeight err %v", host, err.Error())
		}
	}
	fmt.Println("BlockTime PreCheck end", len(b.QkcClients))
	return nil
}
