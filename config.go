package gwork

type ConfigType struct {
	ServerPort    string
	WsUidName     string
	WsRidName     string
	LogQueueSize  uint
	LogBufferSize uint16
	LogLevel      LogLevel
	AdminPort     uint
}
