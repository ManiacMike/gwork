package gwork

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	// "fmt"
)

var HandleRequest func(map[string]interface{}, string, *Room)
var LoseConnCallback func(string, *Room)
var GetConnCallback func(string, *Room)
var GenerateUid func() string
var conf *ConfigType

const (
	Version = "0.1.0"
)

func Start() {
	roomList = make(map[string]*Room)
	http.Handle("/", websocket.Handler(WsServer))

	serverConfig := LoadConfig("server")
	wsConfig := LoadConfig("websocket")
	logConfig := LoadConfig("log")
	adminConfig := LoadConfig("admin")

	converted := convertInt(map[string]string{
		"log_queue_size":  logConfig["log_queue_size"],
		"log_buffer_size": logConfig["log_buffer_size"],
		"log_level":       logConfig["log_level"],
		"ws_param_type":   wsConfig["param_type"],
		"ws_broad_type":   wsConfig["broad_type"],
	})

	conf = &ConfigType{
		ServerPort:    serverConfig["port"],
		WsUidName:     wsConfig["uid_name"],
		WsBroadType:   uint(converted["ws_broad_type"]),
		WsRidName:     wsConfig["rid_name"],
		WsParamType:   uint(converted["ws_param_type"]),
		LogQueueSize:  uint(converted["log_queue_size"]),
		LogBufferSize: uint16(converted["log_buffer_size"]),
		LogLevel:      LogLevel(converted["log_level"]),
		AdminPort:     adminConfig["port"],
	}

	logStart()
	statsStart()
	adminStart()

	rejects := make(chan error, 1)
	go func(port string) {
		Logf(LogLevelNotice, "WebSocket Server listen on port: %s", conf.ServerPort)
		rejects <- http.ListenAndServe(":"+port, nil)
	}(conf.ServerPort)
	select {
	case err := <-rejects:
		log.Fatal("server", "Can't start server: %s", err)
		os.Exit(3)
	}
}

func convertInt(params map[string]string) map[string]int {
	converted := make(map[string]int)
	for k, v := range params {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(k, " convert to int error")
			os.Exit(3)
		}
		converted[k] = i
	}
	return converted
}

func LoadConfig(section string) map[string]string {
	configData, err := GetConfig("config.ini", section)
	if err != nil {
		log.Fatal(section+" config not found:", err)
		os.Exit(1)
	}
	//需要再config.ini中设置的参数
	var neededParams []string
	//默认参数
	var defaultParams map[string]string
	switch section {
	case "server":
		neededParams = []string{"port"}
	case "websocket":
		neededParams = []string{"broad_type"}
		defaultParams = map[string]string{
			"uid_name":   "uid",
			"rid_name":   "room_id",
			"param_type": "get",
		}
	case "log":
		defaultParams = map[string]string{
			"log_queue_size":  "1000",
			"log_buffer_size": "2",
			"log_level":       "1",
		}
	case "admin":
		neededParams = []string{"port"}
	default:
		neededParams = []string{}
		defaultParams = map[string]string{}
	}
	for _, param := range neededParams {
		if _, ok := configData[param]; ok == false {
			log.Fatal(section+" "+param+" config must be set:", err)
			os.Exit(1)
		}
	}
	for k, v := range defaultParams {
		if _, ok := configData[k]; ok == false {
			configData[k] = v
		}
	}
	return configData
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
