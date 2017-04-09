// Package chatlog provides a structure way of storing chatlogs for consumption by CelestialStats.
package chatlog

import (
	"encoding/json"
	"log"
	"os"
	"path"
	_ "strconv"
	"time"
)

type ChatLog struct {
	Protocol   string
	Server     string
	logDir     string
	logFiles   map[string]*os.File
	logChannel chan map[string]string
}

// NewChatLog returns a new ChatLog ready to recieve log entries and write them to disk.
func NewChatLog(LogDir, Protocol string, MaxQueue int) *ChatLog {
	cl := new(ChatLog)
	cl.logDir = LogDir
	cl.Protocol = Protocol
	cl.logFiles = make(map[string]*os.File)
	cl.logChannel = make(chan map[string]string, MaxQueue)
	go cl.Write()
	return cl
}

// Open opens a specific structured file for later writing.
func (chatLog *ChatLog) Open(Server string) *os.File {
	var logFilename = chatLog.GenerateFilename(Server)
	var parentDir = path.Dir(logFilename)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		os.MkdirAll(parentDir, 0755)
	}
	f, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening file:", logFilename)
	}
	return f
}

// AddEntry adds a map representing the chat message to the log channel. It
// also appends the current Unix Timestamp in milliseconds to the map.
func (chatLog *ChatLog) AddEntry(newEntry map[string]string) {
	chatLog.logChannel <- newEntry
}

// Write outputs any additions to the log channel to the current log file.
// If the log does not exist or is old this triggers the log to be opened
// or rotated. All entries are converted to JSON and stored on object per.
// line.
func (chatLog *ChatLog) Write() {
	for i := range chatLog.logChannel {
		curLogFile := chatLog.GetLogHandle(i["Server"])
		jsonVal, _ := json.Marshal(i)
		_, err := curLogFile.WriteString(string(jsonVal) + "\n")
		if err != nil {
			log.Println("Error writing to file:", err)
		}
		curLogFile.Sync()
	}
}

// GetLogHandle returns a pointer to the current log file we should be writing to.
func (chatLog *ChatLog) GetLogHandle(Server string) *os.File {
	if _, ok := chatLog.logFiles[Server]; ok {
		// A log file exists with this server name
		if chatLog.logFiles[Server].Name() != chatLog.GenerateFilename(Server) {
			// Filename doesn't match where we should be writing so close
			// and re-open with new name
			chatLog.logFiles[Server].Close()
			chatLog.logFiles[Server] = chatLog.Open(Server)
		}
	} else {
		// Chatlog for this server isn't open, so open it.
		chatLog.logFiles[Server] = chatLog.Open(Server)
	}
	return chatLog.logFiles[Server]
}

// GenerateFilename returns a filename in the following format
// using the current timestamp:
// $LOGDIR/$PROTOCOL/$SERVER/YYYY/MM/DD/HH.csl
func (chatLog *ChatLog) GenerateFilename(Server string) string {
	return path.Join(
		chatLog.logDir,
		chatLog.Protocol,
		Server,
		time.Now().UTC().Format("2006/01/02/15-04-05.csl"))
	//time.Now().UTC().Format("2006/01/02/15.csl"))
}
