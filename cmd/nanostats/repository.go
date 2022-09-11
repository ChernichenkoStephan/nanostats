package main

import "github.com/ChernichenkoStephan/nanostats/internal/stats"

func initRepository(cfg Config, r *stats.IMRepository) {
	for _, c := range cfg.Chats {
		if c.Username != `` {
			r.Set(stats.Chat{
				Username:   c.Username,
				Shots:      make([]stats.Shot, 0),
				LastPostID: c.LastMsgID,
			})
		}
	}
}
