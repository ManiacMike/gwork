## 线上demo
聊天室: http://120.24.55.22:8001/chat


小蝌蚪游戏: http://120.24.55.22:8006/demo

小蝌蚪项目地址 https://github.com/ManiacMike/go_websocket_chatroom

## 特性
* 高性能的golang websocket 服务器框架
* 支持全局推送，room推送和基于geohash的虚拟地图推送（可用于mmorpg游戏）三种推送方式
* 如果你使用golang，简单就可以实现一个websocket聊天室服务器

~~~ go
package main

import (
	"fmt"
	"github.com/ManiacMike/gwork"
	"time"
)

func main() {
  //设置新建用户连接的callback
	gwork.SetGetConnCallback(func(uid string, room *gwork.Room) {
		welcome := map[string]interface{}{
			"type":       "user_count",
			"user_count": len(room.Userlist),
		}
		room.Broadcast(welcome)
	})

  //设置丢失用户连接的callback
	gwork.SetLoseConnCallback(func(uid string, room *gwork.Room) {
		close := map[string]interface{}{
			"type":       "user_count",
			"user_count": len(room.Userlist),
		}
		room.Broadcast(close)
	})

  //设置处理客户端请求的方法
	gwork.SetRequestHandler(func(receiveNodes map[string]interface{}, uid string, room *gwork.Room) {
		reply := map[string]interface{}{
			"type":    "message",
			"content": receiveNodes["content"].(string),
			"uname":   receiveNodes["uname"].(string),
			"time":    time.Now().Unix(),
		}
		room.Broadcast(reply)
	})

  //读取配置，启动日志，stats，网络监听
	gwork.Start()
}

~~~

* 如果你使用其他的后端语言，请使用gateway的demo
* 简单的服务器状态信息

~~~
Mikes-iMac:~ Mike$ telnet 127.0.0.1 8011
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
stats
===============================
Version: 0.1.0
Uptime: 2 minutes, 19 seconds
Copyright (c) 2016 gwork
*******************************
config:
ServerPort:          8001
LogLevel:            INFO
usage:
Current User Num:    1
Current Room Num:    1
Peak User Num:       1
Peak Room Num:       1
===============================
quit
~~~ 

## 配置及安装
在你的项目下新建config.ini文件
