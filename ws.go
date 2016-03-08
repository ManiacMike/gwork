package gwork

import (
	"golang.org/x/net/websocket"
)

func WsServer(ws *websocket.Conn) {
	var (
		err  error
		room *Room
	)
	uid := ws.Request().FormValue(conf.WsUidName)
	if uid == "" {
		Log(LogLevelInfo, "uid missing")
		if GenerateUid != nil {
			uid = GenerateUid()
		} else {
			uid = GenerateUnixNanoId()
		}
	}

	roomId := ws.Request().FormValue(conf.WsRidName)
	if roomId == "" {
		roomId = "default" //no room param
	}
	room, exist := roomList[roomId]
	if exist == false {
		room = NewRoom(roomId)
	}
	userExist, index := room.ExistUser(uid)
	if userExist == true {
		room.ChangeConn(index, ws)
	} else {
		Log(LogLevelInfo, "create new user")
		uid = room.NewUser(ws, uid)
	}

	for {
		var receiveMsg string
		if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
			room.RemoveUser(uid)
			break
		}
		receiveNodes := JsonDecode(receiveMsg)
		HandleRequest(receiveNodes.(map[string]interface{}), uid, room)
	}
}
