package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/tg"
	mtptg "github.com/gotd/td/tg"
)

// func last[T any](slice []T) T {
// 	return slice[len(slice)-1]
// }

const (
	DEFAULT_POSTS_AMOUNT int = 20
	appId                    = 15656067
	appHash                  = `22e40a55ffc270cd196e10578d1d92da`
)

// func (b Bot) stamp(ctx context.Context, chat tg.Chat, t tg.ShotType) (tg.Shot, error) {
// 	c, err := b.botClient.ChatByUsername(chat.Username)
// 	if err != nil {
// 		return tg.Shot{}, fmt.Errorf("error during chat fetch, got %w\n", err)
// 	}

// 	chatPartis, err := b.botClient.Len(c)
// 	if err != nil {
// 		b.lg.Info(err)
// 	}

// 	return tg.Shot{
// 		Type:    t,
// 		Amount:  chatPartis,
// 		Created: time.Now(),
// 	}, nil

// }

func (b Bot) update(ctx context.Context, c *tg.Chat, t tg.ShotType) error {
	b.lg.Infof("updating %s len: %d", c.Username, len(c.Messages))

	chat, err := b.botClient.ChatByUsername(c.Username)
	if err != nil {
		return fmt.Errorf("error during chat fetch, got %w", err)
	}

	chatPartis, err := b.botClient.Len(chat)
	if err != nil {
		b.lg.Info(err)
	}

	c.TgID = chat.ID
	c.Title = chat.FirstName
	if chat.Title != `` {
		c.Title = chat.Title
	}
	c.Type = string(chat.Type)

	c.Shots = append(c.Shots, tg.Shot{
		Type:    t,
		Amount:  chatPartis,
		Created: time.Now(),
	})

	if len(c.Messages) == 0 {
		// TODO: proper error handle
		return nil
	}

	var start int
	end := len(c.Messages) - 1
	for i := end; i >= 0; i-- {
		start = i

		if c.Messages[i].Views != 0 {
			break
		}
	}

	emptyMessageIDs := make([]int, 0)
	for i := start; i <= end; i++ {
		emptyMessageIDs = append(emptyMessageIDs, c.Messages[i].ID)
	}

	tail := len(emptyMessageIDs) % 6
	tgmm, err := b.mtpClient.getFullMessages(ctx, c.Username, emptyMessageIDs[0:tail])
	if err != nil {
		return fmt.Errorf("error during messages fetch, got %w", err)
	}

	for i := tail; i < len(emptyMessageIDs); i += 6 {

		ctgmm, err := b.mtpClient.getFullMessages(ctx, c.Username, emptyMessageIDs[i:i+6])
		if err != nil {
			return fmt.Errorf("error during messages fetch, got %w", err)
		}

		tgmm = append(tgmm, ctgmm...)
	}

	fillMessages(c.Messages[start:end+1], tgmm)

	return nil
}

func fillMessages(m2Fill []tg.Message, source mtptg.MessageArray) {

	for i := 0; i < len(m2Fill); i++ {

		// just for security
		if m2Fill[i].ID == source[i].ID {

			var frlen int
			for _, rlen := range source[i].Reactions.Results {
				frlen += rlen.Count
			}

			m2Fill[i].Text = source[i].Message
			m2Fill[i].Views = source[i].Views
			m2Fill[i].Reactions = frlen
			m2Fill[i].Comments = source[i].Replies.Replies
			m2Fill[i].Created = time.Unix(int64(source[i].Date), 0)

		}
	}
}
