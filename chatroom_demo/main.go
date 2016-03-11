package main

import (
	"fmt"
	"github.com/ManiacMike/gwork"
	"net/http"
	"time"
)

func StaticServer(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "chat.html")
	return
}

func main() {
	fmt.Println("浏览器访问 http://yourhost:port/chat")
	http.HandleFunc("/chat", StaticServer)

	gwork.SetGetConnCallback(func(uid string, room *gwork.Room) {
		cookie := map[string]interface{}{
			"type": "session",
			"uid":  uid,
		}
		room.PushByUid(uid, cookie)
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
