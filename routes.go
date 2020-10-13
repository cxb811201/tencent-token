package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/basicauth"
)

type ResBody struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}

var message = ResBody{
	Status:      "failed",
	AccessToken: "",
}

func tokenHandler(ctx dotweb.Context) error {
	appid := ctx.QueryString("appid")
	if appid == "" {
		log.Println("ERROR: 没有提供AppID参数")
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}

	atype := AccountType(ctx.QueryString("type"))
	if atype == "" {
		log.Println("ERROR: 没有提供Type参数")
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}

	var account *Account
	for _, _account := range app.Accounts {
		if _account.Type == atype && _account.AppID == appid {
			account = _account
		}
	}

	if account != nil {
		var access_token string
		var record_time string
		var expires_in string

		// 查询数据库中是否已经存在这个AppID的access_token
		record_time = app.Query(account.Type, account.AppID, "timestamp")
		access_token = app.Query(account.Type, account.AppID, "access_token")
		expires_in = app.Query(account.Type, account.AppID, "expires_in")
		expire_time, _ := strconv.ParseInt(record_time, 10, 64)
		timeout, _ := strconv.ParseInt(expires_in, 10, 64)

		if access_token != "" {
			// 如果数据库中已经存在了Token，就检查过期时间，如果过期了就去GetToken获取
			curTime := time.Now().Unix()
			if curTime >= expire_time+timeout-300 {
				token := GetToken(account.Type, account.AppID, account.Secret)
				// 没获得access_token就返回Failed消息
				if token == nil {
					log.Println("ERROR: 没有获得access_token.")
					return ctx.WriteJsonC(http.StatusNotFound, message)
				}

				//获取Token之后更新运行时环境，然后返回access_token
				app.UpdateToken(account.Type, account.AppID, token)
				message.AccessToken = token.AccessToken
			} else {
				message.AccessToken = access_token
			}
		} else {
			token := GetToken(account.Type, account.AppID, account.Secret)
			if token == nil {
				log.Println("ERROR: 没有获得access_token.")
				return ctx.WriteJsonC(http.StatusNotFound, message)
			}
			app.UpdateToken(account.Type, account.AppID, token)
			message.AccessToken = token.AccessToken
		}

		message.Status = "success"
		return ctx.WriteJson(message)
	}

	log.Println("ERROR: AppID或Type无效")
	// 如果提交的appid不在配置文件中，就返回Failed消息
	return ctx.WriteJsonC(http.StatusNotFound, message)
}

func InitRoute(server *dotweb.HttpServer) {
	// 定义Basic Auth的用户名和密码用来防止接口被恶意访问
	var form = map[string]string{
		"user": "api",
		"pass": "wechat",
	}

	option := basicauth.BasicAuthOption{}
	option.Auth = func(name, pwd string) bool {
		if name == form["user"] && pwd == form["pass"] {
			return true
		}
		return false
	}

	server.GET("/token", tokenHandler).Use(basicauth.Middleware(option))
}
