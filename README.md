# go-rock

## 说明
  单通道的tcp复用, 超时重发, keepalive, 只支持单接口, 可用户微信服务器调试

## 设置nginx

```
server{
  listen 80;
  server_name wx.victorup.com;

  location ~ / {
     proxy_pass http://127.0.0.1:18080;
	 proxy_set_header Host $host;
  }
}

```

## 设置服务器

```
{
	"Mode": "server",
	"RemoteAddr": "106.185.26.101:18081",
	"RemoteServerAddr": "127.0.0.1:18080"
}
```
说明
- Mode: 服务器模式 server
- RemoteAddr: 服务器监听的客户端请求地址
- RemoteServerAddr: nginx转发http请求的地址


## 设置客户端

```
{
	"Mode": "server",
	"LocalAddr": "127.0.0.1:80",
	"RemoteAddr": "106.185.26.101:18081",
}
```

说明
- Mode: 服务器模式  client
- LocalAddr: 本地需要代理的地址
- RemoteServerAddr: 远程监听客户端请求的地址

## Todo
目前还没有实现http服务器，只实现了tcp接口映射到远端的服务器
- 实现一个http服务器
- 支持多域名，多用户使用
- 增加密码
