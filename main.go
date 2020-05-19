package main

import (

	//	"time"

	//	"github.com/patrickmn/go-cache"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

const (
	appID                = "wx900ca629d0906a34"
	appSecret            = "1d51ae95382ce4a4b81885b785f469fa"
	accessTokenFetchUrl  = "https://api.weixin.qq.com/cgi-bin/token"
	customServicePostUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send"
)

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

func fetchAccessToken() (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, err
	}

	//Json Decoding
	if bytes.Contains(body, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return "", 0.0, err
		}
		return atr.AccessToken, atr.ExpiresIn, nil
	} else {
		panic("获取access_token 失败")
	}
}

//获取关注者列表
// 一天只能获取100次
func getflist(access_token string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	flist := gjson.Get(string(body), "data.openid").Array()
	return flist
}

// BookUpdateMSG 消息结构
type BookUpdateMSG struct {
	IP     template.DataItem `json:"ip"`
	MsgErr template.DataItem `json:"msg_error"`
}
type TextMsgContent struct {
	Content string `json:"content"`
}

func (t *ToolManager) SendMsg(ip string, msgErr string) {
	t.checkFListUpdate()
	fmt.Println("FList", len(t.WeChatDetail.fList))
	for _, v := range t.WeChatDetail.fList {
		fmt.Println("openID", v)
		bn := BookUpdateMSG{
			IP:     template.DataItem{Value: ip, Color: "#173177"},
			MsgErr: template.DataItem{Value: msgErr, Color: "#173177"},
		}
		msg := template.TemplateMessage2{
			ToUser:     v.Str,
			TemplateId: "JFpYC0Q9dX-eUrY8PRWKjCFCMfjo5MR5fdhFUCLpOhM",
			//TemplateId: "R0omMNf2-Meb0U548lMP57oNbAYthj76JRQrvYyM-aE",//测试数据
			Data:       bn,
		}

		_, err := template.Send(t.WeChatClt, msg)
		fmt.Println(err, msg.Data)
	}

}

func main() {
	manager := NewToolManager([]string{"http://13.228.159.171:38391","http://52.194.81.124:38391"})
	for true {
		manager.CheckStatus()
		time.Sleep(sleepInterval * time.Second)
	}
}

const (
	sleepInterval = 5*60
	maxBlock      = 2
	minPeers      = 3
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
	WeChatClt     *core.Client
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
	ats := core.NewDefaultAccessTokenServer(appID, appSecret, nil)
	m.WeChatClt = core.NewClient(ats, nil)

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
	fmt.Println("更新Flist", "time", time.Now().String(), "len", len(t.WeChatDetail.fList))
}

func (t *ToolManager) CheckStatus() {
	for host, client := range t.MapClient {
		peerNumber, err := client.GetPeers()
		if err != nil {
			t.SendMsg(host, err.Error())
		}
		height, err := client.GetRootBlockHeight()
		if err != nil {
			t.SendMsg(host, err.Error())
		}

		if height-t.MapLastHeight[host] <= maxBlock {
			errMsg := fmt.Sprintf("time %s 节点出现故障：%d秒root的块数%d<=预期增长块数%d, 当前root高度%d, 上次root高度%d  peer个数%d", time.Now().String(), sleepInterval, height-t.MapLastHeight[host], maxBlock, height, t.MapLastHeight[host], peerNumber)
			t.SendMsg(host, errMsg)
		}

		if peerNumber <= minPeers {
			t.SendMsg(host, fmt.Sprintf("节点peer数目%d<=%d", peerNumber, minPeers))
		}
		t.MapLastHeight[host] = height
		fmt.Println("检查完毕", time.Now().String(), "host", host, "root高度", height, "peer数量", peerNumber)
	}
}
