package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
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

var openID = "oMQEkwtXXkTdrWjszFhsnOvFvfC8"

type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

// {
//	"touser":"OPENID",
//	"msgtype":"text",
//	"text":
//	{
//		"content":"Hello World"
//	}
// }
type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	MsgType string         `json:"msgtype"`
	Text    TextMsgContent `json:"text"`
}

type TextMsgContent struct {
	Content string `json:"content"`
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
		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return "", 0.0, err
		}
		return "", 0.0, fmt.Errorf("%s", ater.Errmsg)
	}
}

func pushCustomMsg(accessToken, toUser, msg string) error {
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		fmt.Println("pushCustomMsg err", err)
		return err
	}

	postReq, err := http.NewRequest("POST",
		strings.Join([]string{customServicePostUrl, "?access_token=", accessToken}, ""),
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
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

func (t *ToolManager) SendMsg(msg string) {
	t.checkFListUpdate()
	fmt.Println("准备发送Msg", msg, "FList'Len", len(t.WeChatDetail.fList))
	for _, openID := range t.WeChatDetail.fList {
		fmt.Println("openID", openID, msg)
		err := pushCustomMsg(t.WeChatDetail.accessToken, openID.Str, msg)
		if err != nil {
			log.Println("Push custom service message err:", err)
			return
		}
	}

}
func main() {
	manager := NewToolManager([]string{"http://13.228.159.171:38391","http://52.194.81.124:38391"})
	for true {
		manager.CheckStatus()
		time.Sleep(sleepInterval * time.Second)
	}
}