package gwork

const (
	WsParamTypeGet    = 1
	WsParamTypeCookie = 2
)

type ConfigType struct {
	ServerPort    string
	WsUidName     string
	WsRidName     string
	WsParamType   uint
	LogQueueSize  uint
	LogBufferSize uint16
	LogLevel      LogLevel
	AdminPort     string
}
