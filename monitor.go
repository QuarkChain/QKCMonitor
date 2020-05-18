package main

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

const (
	sleepInterval = 10
	maxBlock      = 1
)

type SendDetail struct {
	fList       []gjson.Result
	accessToken string
	time        int64
}

type ToolManager struct {
	IPList        []string //"http://127.0.0.1:38391"
	MapClient     map[string]*Client
	MapLastHeight map[string]uint64
	WeChatDetail  *SendDetail
}

func NewToolManager(ipList []string) *ToolManager {
	m := &ToolManager{
		IPList:        ipList,
		MapClient:     make(map[string]*Client),
		MapLastHeight: make(map[string]uint64),
		WeChatDetail:  &SendDetail{},
	}
	for _, v := range ipList {
		m.MapClient[v] = NewClient(v)
	}

	m.UpdateFList()
	return m
}

func (t *ToolManager) checkFListUpdate() {
	if time.Now().Unix()-t.WeChatDetail.time > 30*60 {
		t.UpdateFList()
	}
}

func (t *ToolManager) UpdateFList() {
	t.WeChatDetail.time = time.Now().Unix()
	accessToken, _, err := fetchAccessToken()
	if err != nil {
		panic(err)
	}
	t.WeChatDetail.accessToken = accessToken
	t.WeChatDetail.fList = getflist(accessToken)
	fmt.Println("更新Flist", "time",time.Now().String(),"len",len(t.WeChatDetail.fList))
}
func (t *ToolManager) CheckStatus() {
	for host, client := range t.MapClient {
		height, err := client.GetRootBlockHeight()
		if err != nil {
			errMsg := fmt.Sprintf("time %s 节点出现故障：节点ip:%s, jsonRpc调用出错:%v", time.Now().String(), host, err.Error())
			t.SendMsg(errMsg)
		}
		//if height-t.MapLastHeight[host] < maxBlock {
		errMsg := fmt.Sprintf("time %s 节点出现故障：节点ip:%s %d秒root块增长小于%d, 当前root高度%d, 上次root高度%d", time.Now().String(), host, sleepInterval, maxBlock, height, t.MapLastHeight[host])
		t.SendMsg(errMsg)
		t.MapLastHeight[host] = height
		//}
	}
}

