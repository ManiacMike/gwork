package gwork

import (
	"golang.org/x/net/websocket"
)

func WsServer(ws *websocket.Conn) {
	var err error
	uid := ws.Request().FormValue(conf.WsUidName)
	if uid == "" {
		Log(LogLevelInfo, "uid missing")
		if GenerateUid != nil {
			uid = GenerateUid()
		} else {
			uid = GenerateId()
		}
	}

	roomId := ws.Request().FormValue(conf.WsRidName)
	if roomId == "" {
		roomId = "default" //no room param
	}
	room, exist := roomList[roomId]
	if exist == false {
		userlist := []User{}
		room = Room{RoomId: roomId, Userlist: userlist}
		go SendStats(StatsCmdNewRoom)
	}
	userExist, index := room.Exist(uid)
	if userExist == true {
		room.ChangeConn(index, ws)
	} else {
		Log(LogLevelInfo, "create new user")
		uid = room.New(ws, uid)
	}

	for {
		var receiveMsg string
		if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
			room = roomList[room.RoomId]
			room.Remove(uid)
			break
		}
		room = roomList[room.RoomId]
		receiveNodes := JsonDecode(receiveMsg)
		HandleRequest(receiveNodes.(map[string]interface{}), uid, &room)
	}
}
