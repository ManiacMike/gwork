package main

import (
	"fmt"
	"github.com/ManiacMike/gwork"
	"net/http"
	"strconv"
	"time"
)

var TheMap *gwork.MapType

func main() {

	fmt.Println("访问 ip:port/demo/")
	http.Handle("/demo/", http.StripPrefix("/demo/", http.FileServer(http.Dir("Web"))))

	TheMap = gwork.NewMap(64800, 64800)
	gwork.SetGenerateUid(func() string {
		id := int(time.Now().Unix())
		return strconv.Itoa(id)
	})

	gwork.SetGetConnCallback(func(uid string, room *gwork.Room) {
		TheMap.AddCoord(uid, 0, 0)
		welcome := map[string]interface{}{
			"type": "welcome",
			"id":   uid,
		}
		room.Broadcast(welcome)
	})

	gwork.SetLoseConnCallback(func(uid string, room *gwork.Room) {
		TheMap.DeleteCoordNode(uid)
		close := map[string]interface{}{
			"type": "closed",
			"id":   uid,
		}
		room.Broadcast(close)
	})

	gwork.SetRequestHandler(func(receiveNodes map[string]interface{}, uid string, room *gwork.Room) {
		receiveType := receiveNodes["type"]
		if receiveType == "update" {
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
			go TheMap.UpdateCoord(uid, int(x), int(y))
			go BroadcastNearby(x, y, reply, room)
			// room.Broadcast(reply)
		} else if receiveType == "message" {
			reply := map[string]interface{}{
				"type":    "message",
				"id":      uid,
				"message": receiveNodes["message"].(string),
			}
			room.Broadcast(reply)
		}
	})

	gwork.Start()
}

func BroadcastNearby(x float64, y float64, reply map[string]interface{}, room *gwork.Room) {
	uids := TheMap.QueryNearestSquare(int(x), int(y))
	for _, uid := range uids {
		room.PushByUid(uid, reply)
	}
}
