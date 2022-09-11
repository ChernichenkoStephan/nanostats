package stats

import "time"

type Chat struct {
	ID         int64
	TgID       int64
	Username   string
	Title      string
	Type       string
	Shots      []Shot
	LastPostID int
}

type Message struct {
	ID              int
	Text            string
	Views           int
	ReactionsAmount int
	CommentsLen     int
	PostDate        int
}

type Shot struct {
	Messages    []Message
	Subscribers int
	Created     time.Time
}

func New() Chat {
	return Chat{}
}
