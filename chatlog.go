package chatlog

import (
	"time"
	"path"
	"os"
)

type ChatLog struct {
	Protocol string
	Server string
	logDir string
	logFile *os.File
}

func NewChatLog(LogDir, Protocol, Server string) *ChatLog {
	cl := new(ChatLog)
	cl.logDir = LogDir
	cl.Protocol = Protocol
	cl.Server = Server
	return cl
}

func (chatLog *ChatLog) OpenLog() {
	var logFilename = chatLog.ComputeFilename()
	var parentDir = path.Dir(logFilename)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		os.MkdirAll(parentDir, 0755)
	}
	f, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	// Do real error checking here.
	if err != nil {
		chatLog.logFile = f
	}
}

func (chatLog *ChatLog) AddEntry(Initiator, LineType, Content string) {
	
}

func (chatLog *ChatLog) ComputeFilename() string {
	return path.Join(
		chatLog.logDir,
		chatLog.Protocol,
		chatLog.Server,
		time.Now().Format("2006/01/02/15.csl"))
}
