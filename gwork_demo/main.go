package main

import (
	"fmt"
	"github.com/ManiacMike/gwork"
	"net/http"
	"strconv"
	"time"
)

func main() {

	fmt.Println("访问 ip:port/demo/")
	http.Handle("/demo/", http.StripPrefix("/demo/", http.FileServer(http.Dir("Web"))))

	gwork.SetGenerateUid(func() string {
		id := int(time.Now().Unix())
		return strconv.Itoa(id)
	})

	gwork.SetGetConnCallback(func(uid string, room *gwork.Room) {
		welcome := map[string]interface{}{
			"type": "welcome",
			"id":   uid,
		}
		room.Broadcast(gwork.JsonEncode(welcome))
	})

	gwork.SetLoseConnCallback(func(uid string, room *gwork.Room) {
		close := map[string]interface{}{
			"type": "close",
			"id":   uid,
		}
		room.Broadcast(gwork.JsonEncode(close))
	})

	gwork.Init(func(receiveNodes map[string]interface{}, uid string, room *gwork.Room) {
		receiveType := receiveNodes["type"]
		if receiveType == "login" {

		} else if receiveType == "update" {
			var name interface{}
			var ok bool
			if name, ok = receiveNodes["name"]; ok == false {
				name = "Guest." + uid
			}
			x, _ := strconv.ParseFloat(receiveNodes["x"].(string), 64)
			y, _ := strconv.ParseFloat(receiveNodes["y"].(string), 64)
			angle, _ := strconv.ParseFloat(receiveNodes["angle"].(string), 64)
			momentum, _ := strconv.ParseFloat(receiveNodes["momentum"].(string), 64)
			reply := map[string]interface{}{
				"type":       "update",
				"id":         uid,
				"angle":      angle,
				"momentum":   momentum,
				"x":          x,
				"y":          y,
				"life":       1,
				"name":       name,
				"authorized": false,
			}
			room.Broadcast(gwork.JsonEncode(reply))
		} else if receiveType == "message" {
			reply := map[string]interface{}{
				"type":    "message",
				"id":      uid,
				"message": receiveNodes["message"].(string),
			}
			room.Broadcast(gwork.JsonEncode(reply))
		}
	})
}
