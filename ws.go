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
	uid := ws.Request().FormValue(wsConfig["uid"])
	if uid == "" {
		fmt.Println("uid missing")
		if GenerateUid != nil{
			uid = GenerateUid()
		}else{
			uid = GenerateId()
		}
	}
	var roomId string
	if _,ok := wsConfig["room_id"];ok == false{
		roomId = ""
	}else{
		roomId = ws.Request().FormValue(wsConfig["room_id"])
	}
	if roomId == "" {
		roomId = "default" //no room param
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
