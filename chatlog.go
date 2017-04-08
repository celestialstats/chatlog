// Package chatlog provides a structure way of storing chatlogs for consumption by CelestialStats.
package chatlog

import (
	"encoding/json"
	"fmt"
	"log"
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

// NewChatLog returns a new ChatLog ready to recieve log entries and write them to disk.
func NewChatLog(LogDir, Protocol, Server string, MaxQueue int) *ChatLog {
	cl := new(ChatLog)
	cl.logDir = LogDir
	cl.Protocol = Protocol
	cl.Server = Server
	cl.logChannel = make(chan map[string]string, MaxQueue)
	go cl.Write()
	return cl
}

// Open opens a specific structured file for later writing.
func (chatLog *ChatLog) Open() {
	var logFilename = chatLog.GenerateFilename()
	var parentDir = path.Dir(logFilename)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		os.MkdirAll(parentDir, 0755)
	}
	f, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening file:", logFilename)
	}
	chatLog.logFile = f
}

// AddEntry adds a map representing the chat message to the log channel. It
// also appends the current Unix Timestamp in milliseconds to the map.
func (chatLog *ChatLog) AddEntry(testEntry map[string]string) {
	testEntry["Timestamp"] = strconv.FormatInt(time.Now().UTC().UnixNano()/int64(time.Millisecond), 36)
	chatLog.logChannel <- testEntry
}

// Write outputs any additions to the log channel to the current log file.
// If the log does not exist or is old this triggers the log to be opened
// or rotated. All entries are converted to JSON and stored on object per.
// line.
func (chatLog *ChatLog) Write() {
	for i := range chatLog.logChannel {
		chatLog.RotateIfNeeded()
		jsonVal, _ := json.Marshal(i)
		_, err := chatLog.logFile.WriteString(string(jsonVal) + "\n")
		if err != nil {
			log.Println("Error writing to file:", err)
		}
		chatLog.logFile.Sync()
	}
}

// RotateIfNeeded checks if the current ChatLog struct is referencing the
// proper log file. If the reference is incorrect or doesn't exist then this
// function opens the proper log, and if necessary closes the old one.
func (chatLog *ChatLog) RotateIfNeeded() {
	if chatLog.logFile == nil {
		// Open file if not opened
		chatLog.Open()
	} else if chatLog.logFile.Name() != chatLog.GenerateFilename() {
		// If the filename doesn't match where we should be writing
		// close the old file and reopen with a new name.
		chatLog.logFile.Close()
		chatLog.Open()
	}
}

// GenerateFilename returns a filename in the following format
// using the current timestamp:
// $LOGDIR/$PROTOCOL/$SERVER/YYYY/MM/DD/HH.csl
func (chatLog *ChatLog) GenerateFilename() string {
	return path.Join(
		chatLog.logDir,
		chatLog.Protocol,
		chatLog.Server,
		time.Now().UTC().Format("2006/01/02/15.csl"))
}
