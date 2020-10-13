# Tencent Token
微信Access Token中控服务器，用来统一管理各个公众号的access_token，提供统一的接口进行获取和自动刷新Access Token。

使用dotweb重新写了这个项目，用dotweb分支命名，用Echo写的那个分支将不再维护。


## 项目特点

* REST风格的Web服务方式提供一致的获取接口
* 使用Dotweb框架编写的REST API，性能优秀
* 使用BuntDB作为access_token的内存缓存数据库，并且同时支持数据持久化
* 支持Basic Auth的认证方式，需要通过HTTP Basic认证才能访问，增强了安全性


## 安装

```bash
git clone https://github.com/cxb811201/tencent-token.git
cd tencent-token
go get -d -v .
go build
```


## 快速开始

1. 在安装之后，会生成一个 tencent-token 程序，修改account.json文件，把其中的appid和secret替换成你自己的值；如果你的配置文件放在了其他位置，需要使用 -config
参数指定配置文件的路径。
2. 在项目目录中运行 ./tencent-token 就启动了一个web服务。如果操作系统开启了防火墙，需要防火墙开放8000端口的访问；你也可以通过 -port 参数指定其他未使用的端口。
3. 如果有多个微信公众号的access_token需要管理，只需要在 account.json 文件中按格式把你的AppID和AppSecret添加到数组中就可以了。
4. 默认情况下，没有给框架配置日志输出，如果需要定制日志输出的格式，可以使用[dotlog](https://github.com/devfeel/dotlog)包修改这个项目的代码。
5. 本项目的HTTP路由只有一个 `/token`，需要加参数才能访问；格式为：`http://127.0.0.1:8000/token?appid=你的AppID&type=qq或wechat`，而且需要使用用户名和密码才能访问这个地址。


## 参与贡献

**所有的问题都可以发issue请求**

- 如果要做小的修改，请发出PR
- 对于大的改动，在发PR之前请先发issue请求进行讨论
- PR应该包含：
  * 测试用例
  * 文档
- 也可以通过下列方式贡献：
  * 报告issue
  * 推荐新的特性
  * 改善文档


## License

Apache License, Version 2.0

