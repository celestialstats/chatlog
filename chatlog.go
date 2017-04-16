// Package chatlog provides a structure way of storing chatlogs for consumption by CelestialStats.
package chatlog

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ChatLog struct {
	Protocol      string
	Server        string
	rmqHostname   string
	rmqPort       string
	rmqUsername   string
	rmqPassword   string
	rmqLogQueue   string
	rmqConnection *amqp.Connection
	rmqChannel    *amqp.Channel
	logChannel    chan map[string]string
}

// NewChatLog returns a new ChatLog ready to recieve log entries and write them to disk.
func NewChatLog(rmqHostname, rmqPort, rmqUsername, rmqPassword, rmqLogQueue, Protocol string, MaxQueue int) *ChatLog {
	cl := &ChatLog{
		Protocol:    Protocol,
		rmqHostname: rmqHostname,
		rmqPort:     rmqPort,
		rmqUsername: rmqUsername,
		rmqPassword: rmqPassword,
		rmqLogQueue: rmqLogQueue,
		logChannel:  make(chan map[string]string, MaxQueue),
	}
	go cl.queue()
	return cl
}

// Open opens a specific structured file for later writing.
func (chatLog *ChatLog) open() {
	if chatLog.rmqConnection == nil {
		conn, err := amqp.Dial(fmt.Sprintf(
			"amqp://%v:%v@%v:%v/",
			chatLog.rmqUsername,
			chatLog.rmqPassword,
			chatLog.rmqHostname,
			chatLog.rmqPort,
		))
		failOnError(err, "Failed to connect to RabbitMQ")
		chatLog.rmqConnection = conn
	}
	if chatLog.rmqChannel == nil {
		ch, err := chatLog.rmqConnection.Channel()
		failOnError(err, "Failed to open a channel")

		_, err = ch.QueueDeclare(
			chatLog.rmqLogQueue, // name
			true,                // durable
			false,               // delete when unused
			false,               // exclusive
			false,               // no-wait
			nil,                 // arguments
		)
		failOnError(err, "Failed to declare a queue")
		chatLog.rmqChannel = ch
	}
}

// AddEntry adds a map representing the chat message to the log channel. It
// also appends the current Unix Timestamp in milliseconds to the map.
func (chatLog *ChatLog) AddEntry(newEntry map[string]string) {
	newEntry["ServerType"] = "DISCORD"
	chatLog.logChannel <- newEntry
}

// Write outputs any additions to the log channel to the current log file.
// If the log does not exist or is old this triggers the log to be opened
// or rotated. All entries are converted to JSON and stored on object per.
// line.
func (chatLog *ChatLog) queue() {
	for i := range chatLog.logChannel {
		jsonVal, _ := json.Marshal(i)
		chatLog.open()
		err := chatLog.rmqChannel.Publish(
			"",                  // exchange
			chatLog.rmqLogQueue, // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(jsonVal),
			})
		failOnError(err, "Unable to submit to queue.")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
