package main

import (
	"encoding/json"

	"github.com/levigross/grequests"
)

const (
	WeChatAccessTokenAPI = "https://api.weixin.qq.com/cgi-bin/token"
	QQAccessTokenAPI = "https://api.q.qq.com/api/getToken"
)

type Token struct {
	AccessToken string `json:"access_token"`
	Expire      int    `json:"expires_in"`
}

// 获取AppID的access_token
func GetToken(atype AccountType, appid string, secret string) *Token {
	var params = map[string]string{
		"appid":      appid,
		"secret":     secret,
		"grant_type": "client_credential",
	}

	ro := &grequests.RequestOptions{
		Params: params,
	}

	var api = WeChatAccessTokenAPI
	if atype == QQ {
		api = QQAccessTokenAPI
	}

	res, _ := grequests.Get(api, ro)

	var token *Token
	if err := json.Unmarshal(res.Bytes(), &token); err != nil {
		return nil
	}

	return token
}
