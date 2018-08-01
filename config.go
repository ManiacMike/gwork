package gwork

const (
	WsParamTypeGet    = 1
	WsParamTypeCookie = 2
)

type ConfigType struct {
	ServerPort    string
	WsUidName     string
	WsBroadType   uint
	WsRidName     string
	WsParamType   uint
	WsTlsEnable   uint
	WsTlsCrt      string
	WsTlsKey      string
	LogQueueSize  uint
	LogBufferSize uint16
	LogLevel      LogLevel
	AdminPort     string
}
