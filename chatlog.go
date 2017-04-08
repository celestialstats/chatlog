package chatlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type ChatLog struct {
	Protocol   string
	Server     string
	logDir     string
	logFile    *os.File
	logChannel chan map[string]string
}

func NewChatLog(LogDir, Protocol, Server string, MaxQueue int) *ChatLog {
	cl := new(ChatLog)
	cl.logDir = LogDir
	cl.Protocol = Protocol
	cl.Server = Server
	cl.logChannel = make(chan map[string]string, MaxQueue)
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
		fmt.Println("ERROR")
	}
	chatLog.logFile = f
	fmt.Println(chatLog.logFile.Name())
	//defer chatLog.logFile.Close()
}

func (chatLog *ChatLog) AddEntry(testEntry map[string]string) {
	testEntry["Timestamp"] = strconv.FormatInt(time.Now().UTC().UnixNano()/int64(time.Millisecond), 36)
	chatLog.logChannel <- testEntry
}

func (chatLog *ChatLog) WriteLog() {
	for i := range chatLog.logChannel {
		chatLog.RotateIfNeeded()
		jsonVal, _ := json.Marshal(i)
		_, err := chatLog.logFile.WriteString(string(jsonVal) + "\n")
		if err != nil {
			panic(err)
		}
		chatLog.logFile.Sync()
		fmt.Println(">", string(jsonVal), " (", len(chatLog.logChannel), ")")
	}
}

func (chatLog *ChatLog) RotateIfNeeded() {
	if chatLog.logFile == nil {
		// Open file if not opened
		chatLog.OpenLog()
	} else if chatLog.logFile.Name() != chatLog.ComputeFilename() {
		// If the filename doesn't match where we should be writing
		// close the old file and reopen with a new name.
		chatLog.logFile.Close()
		chatLog.OpenLog()
	}
}

func (chatLog *ChatLog) ComputeFilename() string {
	return path.Join(
		chatLog.logDir,
		chatLog.Protocol,
		chatLog.Server,
		time.Now().UTC().Format("2006/01/02/15.csl"))
}
