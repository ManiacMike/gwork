package gwork

import (
	"golang.org/x/net/websocket"
	// "github.com/gorilla/websocket"
)

const (
	WsBroadTypeAll  = 1
	WsBroadTypeRoom = 2
	WsBroadTypeMap  = 3
)

func WsServer(ws *websocket.Conn) {
	var (
		err    error
		roomId string
		room   *Room
	)
	uid := getParam(ws, "uid")
	if uid == "" {
		Log(LogLevelInfo, "uid missing")
		if GenerateUid != nil {
			uid = GenerateUid()
		} else {
			uid = GenerateUnixNanoId()
		}
	}
	if conf.WsBroadType == WsBroadTypeRoom {
		roomId = getParam(ws, "rid")
	} else {
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

func getParam(ws *websocket.Conn, name string) string {
	var key string
	if name == "uid" {
		key = conf.WsUidName
	} else {
		key = conf.WsRidName
	}
	if conf.WsParamType == WsParamTypeGet {
		return ws.Request().FormValue(key)
	} else {
		if uidCookie, err := ws.Request().Cookie(key); err != nil {
			return ""
		} else {
			return uidCookie.Value
		}
	}
}
