package main

import "github.com/ChernichenkoStephan/nanostats/internal/tg"

func makeEmptyMessages(lastID int) []tg.Message {
	res := make([]tg.Message, 0)

	id := lastID - 9
	if lastID < 10 {
		id = 0
	}

	for ; id <= lastID; id++ {
		res = append(res, tg.Message{ID: id})
	}
	return res
}

func initRepository(cfg Config, r *tg.IMRepository) {
	for _, c := range cfg.Chats {
		if c.Username != `` {
			mm := makeEmptyMessages(c.LastMsgID)
			r.Set(tg.Chat{
				Username: c.Username,
				Shots:    make([]tg.Shot, 0),
				Messages: mm,
			})
		}
	}
}
