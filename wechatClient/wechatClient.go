package wechatClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/516108736/QKCMainnetMonitor/common"
	"github.com/tidwall/gjson"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

const (
	accessTokenFetchUrl = "https://api.weixin.qq.com/cgi-bin/token"
)

type WechatClient struct {
	weChatClt *core.Client
	config    common.WeChatDetail
	fList     FList
}

type FList struct {
	FList          []gjson.Result
	lastUpdateTime int64
}

func NewWeChatClient(config common.WeChatDetail) *WechatClient {
	c := &WechatClient{
		config: config,
	}
	c.UpdateFList()
	ats := core.NewDefaultAccessTokenServer(config.AppID, config.AppSecret, nil)
	c.weChatClt = core.NewClient(ats, nil)
	return c
}

func (t *WechatClient) UpdateFList() {
	t.fList.lastUpdateTime = time.Now().Unix()
	accessToken, _, err := t.fetchAccessToken()
	if err != nil {
		panic(err)
	}
	t.fList.FList = getflist(accessToken)
	fmt.Println("update Flist", "time", time.Now().String(), "len", len(t.fList.FList))
}

func (t *WechatClient) checkFListUpdate() {
	if time.Now().Unix()-t.fList.lastUpdateTime > 30*60 {
		t.UpdateFList()
	}
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

func (t *WechatClient) fetchAccessToken() (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		t.config.AppID,
		"&secret=",
		t.config.AppSecret}, "")

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
		fmt.Println("body", string(body))
		panic("获取access_token 失败")
	}
}

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

func (t *WechatClient) SendMsg(msg string) {
	t.checkFListUpdate()
	for _, v := range t.fList.FList {
		type AlertMsg struct {
			MsgErr template.DataItem `json:"msg_error"`
		}
		bn := AlertMsg{
			MsgErr: template.DataItem{Value: msg, Color: "#173177"},
		}
		msg := template.TemplateMessage2{
			ToUser:     v.Str,
			TemplateId: t.config.TemplateId,
			Data:       bn,
		}

		_, err := template.Send(t.weChatClt, msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("发送成功", v)
		break
	}

}
