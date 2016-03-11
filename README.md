## 线上demo
聊天室: http://120.24.55.22:8001/chat


小蝌蚪游戏: http://120.24.55.22:8006/demo

小蝌蚪项目地址 https://github.com/ManiacMike/go_websocket_chatroom

## 特性
* 高性能的golang websocket 服务器框架
* 支持全局推送，room推送和基于geohash的虚拟地图推送（可用于mmorpg游戏）三种推送方式
* 如果你使用golang，函数式编程，简单几行代码就可以实现一个websocket服务器

~~~ go
package main

import (
	"fmt"
	"github.com/ManiacMike/gwork"
	"time"
)

func main() {
	gwork.SetGetConnCallback(func(uid string, room *gwork.Room) {
		welcome := map[string]interface{}{
			"type":       "user_count",
			"user_count": len(room.Userlist),
		}
		room.Broadcast(welcome)
	})

	gwork.SetLoseConnCallback(func(uid string, room *gwork.Room) {
		close := map[string]interface{}{
			"type":       "user_count",
			"user_count": len(room.Userlist),
		}
		room.Broadcast(close)
	})

	gwork.SetRequestHandler(func(receiveNodes map[string]interface{}, uid string, room *gwork.Room) {
		reply := map[string]interface{}{
			"type":    "message",
			"content": receiveNodes["content"].(string),
			"uname":   receiveNodes["uname"].(string),
			"time":    time.Now().Unix(),
		}
		room.Broadcast(reply)
	})

	gwork.Start()
}

~~~

* 如果你使用其他的后端语言，请使用gateway的demo


## 配置及安装
在你的项目下新建config.ini文件
