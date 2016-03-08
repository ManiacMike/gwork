package gwork

import (
	"fmt"
	"io"
	"net"
)

// Start admin server.
// Handle admin commands by TCP protocol.
func adminStart() {
	go func() {
		l, err := net.Listen("tcp", ":"+conf.AdminPort)
		if err != nil {
			Log(LogLevelError, err.Error())
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				Log(LogLevelError, err)
			}
			go handleCommand(conn)
		}
	}()
}

func handleCommand(c net.Conn) {
	defer c.Close()
	for {
		inBuf := make([]byte, 128)
		n, err := c.Read(inBuf)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(c, "read command error: %s\n", err)
			}
			continue
		}
		cmd := string(inBuf[:n-2])
		Logf(LogLevelInfo, "admin command: %s", cmd)
		switch cmd {
		case "stats":
			outBuf := StatsReport()
			fmt.Fprintln(c, outBuf)
		case "quit": // close connection
			return
		default:
			fmt.Fprintf(c, "unknown admin command: %s\n", cmd)
		}
	}
}
