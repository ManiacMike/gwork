package gwork

import (
	"errors"
	"fmt"
	"time"
)

type ServerStats struct {
	UserCurrent uint
	UserPeak    uint
	RoomCurrent uint
	RoomPeak    uint
	Version     string
	StartTime   time.Time
}

const (
	StatsCmdNewUser = iota
	StatsCmdLostUser
	StatsCmdNewRoom
	StatsCmdCloseRoom
	StatsCmdReport
)

// A channel use to transfer status command
var statsChannel chan statsCmd

// the global status server
var stats *ServerStats

// A type use to store status command
type statsCmd struct {
	cmd       int
	replyChan chan string
}

func (stats *ServerStats) HandleCommand(cmd statsCmd) error {
	switch cmd.cmd {
	case StatsCmdNewUser:
		stats.UserCurrent++
		if stats.UserPeak < stats.UserCurrent {
			stats.UserPeak = stats.UserCurrent
		}
	case StatsCmdLostUser:
		stats.UserCurrent--
	case StatsCmdNewRoom:
		stats.RoomCurrent++
		if stats.RoomPeak < stats.RoomCurrent {
			stats.RoomPeak = stats.RoomCurrent
		}
	case StatsCmdCloseRoom:
		stats.RoomCurrent--
	case StatsCmdReport:
		cmd.replyChan <- stats.Report()
	default:
		Log(LogLevelWarning, "[stats server] command not found. command: %s", cmd)
		return errors.New("command not found")
	}
	return nil
}

func (stats *ServerStats) Report() string {
	uptime := UptimeFormat(uint32(time.Now().Sub(stats.StartTime)/time.Second), 2)

	return fmt.Sprintf(`===============================
Version: %s
Uptime: %s
Copyright (c) 2016 gwork
*******************************
config:
  ServerPort:          %s
  LogLevel:            %s
usage:
  Current User Num:    %d
  Current Room Num:    %d
  Peak User Num:       %d
  Peak Room Num:       %d
===============================`,
		stats.Version,
		uptime,
		conf.ServerPort,
		getLevelString(conf.LogLevel),
		stats.UserCurrent,
		stats.RoomCurrent,
		stats.UserPeak,
		stats.RoomPeak,
	)
}

// Start status server to collect status of the server
func statsStart() {
	// init server status
	stats = &ServerStats{
		Version:   Version,
		StartTime: time.Now(),
	}
	statsChannel = make(chan statsCmd, 5)
	go func(stats *ServerStats) {
		for {
			select {
			case cmd := <-statsChannel:
				stats.HandleCommand(cmd)
			}
		}
	}(stats)
}

// Send status command to status server
func SendStats(cmdCode int) (replyChan chan string) {
	var cmd statsCmd
	// only report command returns result
	// others commands have no result, so they don't need reply channel
	switch cmdCode {
	case StatsCmdReport:
		replyChan = make(chan string)
	default:
		replyChan = nil
	}
	cmd = statsCmd{cmdCode, replyChan}
	statsChannel <- cmd
	return
}

func StatsReport() string {
	replyChan := SendStats(StatsCmdReport)
	return <-replyChan
}
