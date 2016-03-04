package gwork

import (
	"fmt"
	"golang.org/x/net/websocket"
	"strings"
)

var roomList map[string]Room //在线room列表
var HandleRequest func(map[string]interface{}, string, *Room)

func Init(f func(map[string]interface{}, string, *Room)) {
	roomList = make(map[string]Room)
	HandleRequest = f
}

func GetRoom(rid string) Room {
	return roomList[rid]
}

type User struct {
	uid string
	con *websocket.Conn
}

type UserList []User

type Room struct {
	RoomId   string
	Userlist []User
}

func (room *Room) New(ws *websocket.Conn, uid string) string {
	room.Userlist = append(room.Userlist, User{uid, ws})
	fmt.Println("New user connect current user num", len(room.Userlist))
	go room.PushUserCount("user_connect", uid)
	roomList[room.RoomId] = *room
	return uid
}

func (room *Room) Remove(uid string) {
	flag, find := room.Exist(uid)
	fmt.Println("user disconnect uid: ", uid)
	if flag == true {
		room.Userlist = append(room.Userlist[:find], room.Userlist[find+1:]...)
		go room.PushUserCount("user_disconnect", uid)
		roomList[room.RoomId] = *room
	}
}

func (room *Room) ChangeConn(index int, con *websocket.Conn) {
	fmt.Println("visitor exist change connection")
	curUser := (room.Userlist)[index]
	curUser.con.Close()
	(room.Userlist)[index].con = con
	roomList[room.RoomId] = *room
}

func (room *Room) Exist(uid string) (bool, int) {
	var find int
	flag := false
	for i, v := range room.Userlist {
		if uid == v.uid {
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
		userlist = append(userlist, user.uid)
	}
	userCount := UserCountChangeReply{event, uid, len(room.Userlist), strings.Join(userlist, ",")}
	replyBodyStr := JsonEncode(userCount)
	room.Broadcast(replyBodyStr)
}

func (room *Room) Broadcast(replyBodyStr string) error {
	fmt.Println("Broadcast ", room.RoomId, " room user", len(room.Userlist))
	for _, user := range room.Userlist {
		if err := websocket.Message.Send(user.con, replyBodyStr); err != nil {
			fmt.Println("Can't send user ", user.uid, " lost connection")
			room.Remove(user.uid)
			break
		}
	}
	return nil
}

func (room *Room) Push(user User, replyBodyStr string) error {
	fmt.Println("Push ", room.RoomId, user.uid)
	if err := websocket.Message.Send(user.con, replyBodyStr); err != nil {
		fmt.Println("Can't send user ", user.uid, " lost connection")
		room.Remove(user.uid)
	}
	return nil
}
