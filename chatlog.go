package chatlog

import (
	"fmt"
	"time"
	"path"
	"os"
)

type ChatLog struct {
	Protocol string
	Server string
	logDir string
	logFile *os.File
	logChannel chan LogLine
}

func NewChatLog(LogDir, Protocol, Server string, MaxQueue int) *ChatLog {
	cl := new(ChatLog)
	cl.logDir = LogDir
	cl.Protocol = Protocol
	cl.Server = Server
	cl.logChannel = make(chan LogLine, MaxQueue)
	go cl.WriteLog()
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

func (chatLog *ChatLog) AddEntry(Timestamp time.Time, Initiator, LineType, Content string) {
	chatLog.logChannel <- LogLine {
		Timestamp: Timestamp,
		Initiator: Initiator,
		LineType: LineType,
		Content: Content,
	}
}

func (chatLog *ChatLog) WriteLog() {
	for i := range chatLog.logChannel {
		time.Sleep(time.Duration(500)*time.Millisecond)
		fmt.Println(">", i.Timestamp.UnixNano(), i.Content, " (", len(chatLog.logChannel), ")")
		
	}
}

func (chatLog *ChatLog) ComputeFilename() string {
	return path.Join(
		chatLog.logDir,
		chatLog.Protocol,
		chatLog.Server,
		time.Now().UTC().Format("2006/01/02/15.csl"))
}
