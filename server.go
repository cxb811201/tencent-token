package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/devfeel/dotweb"
	"github.com/levigross/grequests"
	"github.com/tidwall/buntdb"
)

type App struct {
	Accounts []*Account
	DB       *buntdb.DB
	Web      *dotweb.DotWeb
}

type AccountType string

const (
	QQ     AccountType = "qq"
	WeChat AccountType = "wechat"
)

type Account struct {
	Type   AccountType `json:"type"`
	AppID  string      `json:"appid"`
	Secret string      `json:"secret"`
}

func NewApp() *App {
	a := &App{}
	a.Web = dotweb.New()

	return a
}

// 读取配置文件中的appid和secret值到一个map中
func (a *App) SetAccounts(config *string) {
	if _, err := os.Stat(*config); err != nil {
		fmt.Println("配置文件无法打开！")
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Println("配置文件读取失败！")
		os.Exit(1)
	}

	if err := json.Unmarshal(raw, &a.Accounts); err != nil {
		fmt.Println("配置文件内容错误！")
		os.Exit(1)
	}
}

func (a *App) Query(atype AccountType, appid string, key string) string {
	var value string

	err := a.DB.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(fmt.Sprintf("%s_%s_%s", atype, appid, key))
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		value = ""
	}

	return value
}

// 更新AppID上下文环境中的Access Token和到期时间
func (a *App) UpdateToken(atype AccountType, appid string, token *Token) {
	timestamp := time.Now().Unix()

	a.DB.Update(func(tx *buntdb.Tx) error {
		tx.Delete(fmt.Sprintf("%s_%s_timestamp", atype, appid))
		tx.Delete(fmt.Sprintf("%s_%s_access_token", atype, appid))
		tx.Delete(fmt.Sprintf("%s_%s_expires_in", atype, appid))

		tx.Set(fmt.Sprintf("%s_%s_timestamp", atype, appid), strconv.FormatInt(timestamp, 10), nil)
		tx.Set(fmt.Sprintf("%s_%s_access_token", atype, appid), token.AccessToken, nil)
		tx.Set(fmt.Sprintf("%s_%s_expires_in", atype, appid), strconv.Itoa(token.Expire), nil)
		return nil
	})
}

// 启动AccessToken是否有效自动检测任务
func (a *App) StartTokenCheckTask() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for _, account := range a.Accounts {
				if account.Type != WeChat {
					continue
				}
				token := app.Query(account.Type, account.AppID, "access_token")
				if token != "" {
					fmt.Printf("access_toke: %s\n", token)
					ro := &grequests.RequestOptions{
						Params: map[string]string{"access_token": token},
					}

					res, _ := grequests.Get("https://api.weixin.qq.com/cgi-bin/getcallbackip", ro)

					var m map[string]interface{}
					if err := json.Unmarshal(res.Bytes(), &m); err == nil {
						if _, ok := m["errcode"]; ok {
							if token := GetToken(account.Type, account.AppID, account.Secret); token != nil {
								app.UpdateToken(account.Type, account.AppID, token)
							}
						}
					}
				}
			}
		}
	}()
}
