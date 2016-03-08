package gwork

import (
	"golang.org/x/net/websocket"
	"strings"
)

var roomList map[string]Room //在线room列表

func GetRoomUser(rid string) UserList {
	return roomList[rid].Userlist
}

type User struct {
	Uid string
	Con *websocket.Conn
}

type UserList []User

type Room struct {
	RoomId   string
	Userlist []User
}

func (room *Room) New(ws *websocket.Conn, uid string) string {
	room.Userlist = append(room.Userlist, User{uid, ws})
	Log(LogLevelInfo, "New user connect current user num", len(room.Userlist))
	if GetConnCallback != nil {
		go GetConnCallback(uid, room)
	}
	roomList[room.RoomId] = *room
	return uid
}

func (room *Room) Remove(uid string) {
	flag, find := room.Exist(uid)
	Log(LogLevelInfo, "user disconnect uid: ", uid)
	if flag == true {
		room.Userlist = append(room.Userlist[:find], room.Userlist[find+1:]...)
		roomList[room.RoomId] = *room
		if LoseConnCallback != nil {
			go LoseConnCallback(uid, room)
		}
	}
}

func (room *Room) ChangeConn(index int, con *websocket.Conn) {
	curUser := (room.Userlist)[index]
	curUser.Con.Close()
	(room.Userlist)[index].Con = con
	roomList[room.RoomId] = *room
}

func (room *Room) Exist(uid string) (bool, int) {
	var find int
	flag := false
	for i, v := range room.Userlist {
		if uid == v.Uid {
			find = i
			flag = true
			break
		}
	}
	return flag, find
}

func (room *Room) PushUserCount(event string, uid string) {
	userlist := []string{}
	for _, user := range room.Userlist {
		userlist = append(userlist, user.Uid)
	}
	replyBody := map[string]interface{}{
		"type":       event,
		"uid":        uid,
		"user_count": len(room.Userlist),
		"user_list":  strings.Join(userlist, ","),
	}
	room.Broadcast(replyBody)
}

func (room *Room) Broadcast(replyBody map[string]interface{}) error {
	replyBodyStr := JsonEncode(replyBody)
	for _, user := range room.Userlist {
		if err := websocket.Message.Send(user.Con, replyBodyStr); err != nil {
			room.Remove(user.Uid)
			break
		}
	}
	return nil
}

func (room *Room) Push(user User, replyBody map[string]interface{}) error {
	replyBodyStr := JsonEncode(replyBody)
	if err := websocket.Message.Send(user.Con, replyBodyStr); err != nil {
		room.Remove(user.Uid)
	}
	return nil
}
