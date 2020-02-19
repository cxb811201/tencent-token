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
	Accounts map[string]string
	DB       *buntdb.DB
	Web      *dotweb.DotWeb
	WxToken  *Token
}

type Account struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret"`
}

func NewApp() *App {
	a := &App{}
	a.Accounts = make(map[string]string)
	a.Web = dotweb.New()
	a.WxToken = new(Token)

	return a
}

// 读取配置文件中的appid和secret值到一个map中
func (a *App) SetAccounts(config *string) {
	accounts := make([]Account, 1)

	if _, err := os.Stat(*config); err != nil {
		fmt.Println("配置文件无法打开！")
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Println("配置文件读取失败！")
		os.Exit(1)
	}

	if err := json.Unmarshal(raw, &accounts); err != nil {
		fmt.Println("配置文件内容错误！")
		os.Exit(1)
	}

	for _, acc := range accounts {
		a.Accounts[acc.AppID] = acc.Secret
	}
}

func (a *App) Query(appid string, key string) string {
	var value string

	err := a.DB.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(appid + "_" + key)
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
func (a *App) UpdateToken(appid string) {
	timestamp := time.Now().Unix()

	a.DB.Update(func(tx *buntdb.Tx) error {
		tx.Delete(appid + "_timestamp")
		tx.Delete(appid + "_access_token")
		tx.Delete(appid + "_expires_in")

		tx.Set(appid+"_timestamp", strconv.FormatInt(timestamp, 10), nil)
		tx.Set(appid+"_access_token", a.WxToken.AccessToken, nil)
		tx.Set(appid+"_expires_in", strconv.Itoa(a.WxToken.Expire), nil)
		return nil
	})
}

// 启动AccessToken是否有效自动检测任务
func (a *App) StartTokenCheckTask() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for appid, secret := range a.Accounts {
				token := app.Query(appid, "access_token")
				if token != "" {
					fmt.Printf("access_toke: %s\n", token)
					ro := &grequests.RequestOptions{
						Params: map[string]string{"access_token": token},
					}

					res, _ := grequests.Get("https://api.weixin.qq.com/cgi-bin/getcallbackip", ro)

					var m map[string]interface{}
					if err := json.Unmarshal(res.Bytes(), &m); err == nil {
						if _, ok := m["errcode"]; ok {
							_ = app.WxToken.Get(appid, secret)
							app.UpdateToken(appid)
						}
					}
				}
			}
		}
	}()
}
