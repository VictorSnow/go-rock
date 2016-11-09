# go-rock

## 说明
  单通道的tcp复用, 超时重发, keepalive, 速度上会受限, 目前只支持单用户使用

## 设置nginx

```
server{
  listen 80;
  server_name wx.victorup.com;

  location ~ / {
     proxy_pass http://127.0.0.1:18080;
  }
}

```

## 设置服务器

```
servConfig := &ngServerConfig{"0.0.0.0:18080", "0.0.0.0:18081"}
server := newNgServer(servConfig)
server.start()
```

## 链接服务器

todo


# 说明

目前只验证了相关的协议方面，不稳定
