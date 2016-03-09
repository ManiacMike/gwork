package gwork

import (
	"golang.org/x/net/websocket"
	"strings"
)

var roomList map[string]*Room //在线room列表

type User struct {
	Uid string
	Con *websocket.Conn
}

type UserList []User

type Room struct {
	RoomId   string
	Userlist []User
}

func NewRoom(roomId string) *Room {
	userlist := []User{}
	room := Room{RoomId: roomId, Userlist: userlist}
	go SendStats(StatsCmdNewRoom)
	roomList[roomId] = &room
	return &room
}

func (room *Room) NewUser(ws *websocket.Conn, uid string) string {
	room.Userlist = append(room.Userlist, User{uid, ws})
	Log(LogLevelInfo, "New user connect current user num", len(room.Userlist))
	if GetConnCallback != nil {
		go GetConnCallback(uid, room)
	}
	go SendStats(StatsCmdNewUser)
	return uid
}

func (room *Room) RemoveUser(uid string) {
	flag, find := room.ExistUser(uid)
	Log(LogLevelInfo, "user disconnect uid: ", uid)
	if flag == true {
		room.Userlist = append(room.Userlist[:find], room.Userlist[find+1:]...)
		if LoseConnCallback != nil {
			go LoseConnCallback(uid, room)
		}
		go SendStats(StatsCmdLostUser)
		if len(room.Userlist) == 0 {
			delete(roomList, room.RoomId)
			go SendStats(StatsCmdCloseRoom)
		}
	}
}

func (room *Room) ChangeConn(index int, con *websocket.Conn) {
	curUser := (room.Userlist)[index]
	curUser.Con.Close()
	(room.Userlist)[index].Con = con
}

func (room *Room) ExistUser(uid string) (bool, int) {
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
			room.RemoveUser(user.Uid)
			break
		}
	}
	return nil
}

func (room *Room) Push(user User, replyBody map[string]interface{}) error {
	replyBodyStr := JsonEncode(replyBody)
	if err := websocket.Message.Send(user.Con, replyBodyStr); err != nil {
		room.RemoveUser(user.Uid)
	}
	return nil
}

func (room *Room) PushByUid(uid string, replyBody map[string]interface{}) error {
	for _, user := range room.Userlist {
		if uid == user.Uid {
			return room.Push(user, replyBody)
		}
	}
	return nil
}
