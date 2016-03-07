package gwork

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
)

var HandleRequest func(map[string]interface{}, string, *Room)
var LoseConnCallback func(string, *Room)
var GetConnCallback func(string, *Room)
var wsConfig map[string]string
var GenerateUid func() string

func Start() {
	var err error
	roomList = make(map[string]Room)
	http.Handle("/", websocket.Handler(WsServer))
	serverConfig, err := GetConfig("config.ini", "server")
	if err != nil {
		log.Fatal("server config error:", err)
		os.Exit(1)
	}
	wsConfig, err = GetConfig("config.ini", "websocket")
	if err != nil {
		log.Fatal("websocket config error:", err)
		os.Exit(1)
	}
	if _, ok := wsConfig["uid"]; ok == false {
		log.Fatal("websocket config uid error:", err)
		os.Exit(1)
	}

	fmt.Println("WebSocket Server listen on port:", serverConfig["port"])

	rejects := make(chan error, 1)
	go func(port string) {
		rejects <- http.ListenAndServe(":"+port, nil)
	}(serverConfig["port"])

	select {
	case err := <-rejects:
		log.Fatal("server", "Can't start server: %s", err)
		os.Exit(3)
	}
}

func SetRequestHandler(f func(map[string]interface{}, string, *Room)) {
	HandleRequest = f
}

func SetGetConnCallback(f func(string, *Room)) {
	GetConnCallback = f
}

func SetLoseConnCallback(f func(string, *Room)) {
	LoseConnCallback = f
}

func SetGenerateUid(f func() string) {
	GenerateUid = f
}
