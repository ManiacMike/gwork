package gwork

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type UserCountChangeReply struct {
	Type      string `json:"type"`
	Uid       string `json:"uid"`
	UserCount int    `json:"user_count"`
	UserList  string `json:"user_list"`
}

func WsServer(ws *websocket.Conn) {
	var err error
	uid := ws.Request().FormValue("uid")
	if uid == "" {
		fmt.Println("uid missing")
		return
	}
	roomId := ws.Request().FormValue("room_id")
	if roomId == "" {
		fmt.Println("roomId missing")
		return
	}
	room, exist := roomList[roomId]
	if exist == false {
		userlist := []User{}
		room = Room{RoomId: roomId, Userlist: userlist}
	}
	userExist, index := room.Exist(uid)
	if userExist == true {
		room.ChangeConn(index, ws)
	} else {
		fmt.Println("create new user")
		uid = room.New(ws, uid)
	}

	for {
		var receiveMsg string
		if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
			room = roomList[room.RoomId]
			fmt.Println("Can't receive,user ", uid, " lost connection")
			room.Remove(uid)
			break
		}
		room = roomList[room.RoomId]
		receiveNodes := JsonDecode(receiveMsg)
		HandleRequest(receiveNodes.(map[string]interface{}), uid, &room)
	}
}
