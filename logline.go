package chatlog

import (
	"time"
)

type LogLine struct {
	Timestamp time.Time
	Initiator string
	LineType string
	Content string
}
