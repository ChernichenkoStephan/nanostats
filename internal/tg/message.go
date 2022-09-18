package tg

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int
	Text      string
	Views     int
	Reactions int
	Comments  int
	Created   time.Time
}

func (m Message) CreatedString() string {
	t := m.Created
	return fmt.Sprintf("%02d.%02d[%02d:%02d]",
		t.Day(), t.Month(), t.Hour(), t.Minute())
}
