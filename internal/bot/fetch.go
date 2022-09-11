package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/stats"
	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

const (
	DEFAULT_POSTS_AMOUNT int = 20
	appId                    = 15656067
	appHash                  = `22e40a55ffc270cd196e10578d1d92da`
)

func idsToInputMessageClass(ids []int) []tg.InputMessageClass {
	imcs := make([]tg.InputMessageClass, 0)
	for _, id := range ids {
		imcs = append(imcs, &tg.InputMessageID{ID: id})
	}
	return imcs
}

func (b Bot) shot(username string, id int) (stats.Shot, error) {
	chat, err := b.botClient.ChatByUsername(username)
	if err != nil {
		return stats.Shot{}, fmt.Errorf("error during chat fetch, got %w", err)
	}

	chatPartis, err := b.botClient.Len(chat)
	if err != nil {
		b.lg.Info(err)
	}

	ids := b.validIDs(id)

	tgmm, err := b.getFullMessages(username, ids)
	if err != nil {
		return stats.Shot{}, fmt.Errorf("error during messages fetch, got %w", err)
	}

	mm := make([]stats.Message, 0)
	for _, tgm := range tgmm {

		var frlen int
		for _, rlen := range tgm.Reactions.Results {
			frlen += rlen.Count
		}

		m := stats.Message{
			ID:              tgm.ID,
			Text:            tgm.Message,
			Views:           tgm.Views,
			ReactionsAmount: frlen,
			CommentsLen:     tgm.Replies.Replies,
			PostDate:        tgm.Date,
		}

		mm = append(mm, m)

	}

	return stats.Shot{
		Messages:    mm,
		Subscribers: chatPartis,
		Created:     time.Now(),
	}, nil

}

func (b Bot) validIDs(id int) (ids []int) {
	ids = make([]int, 0)
	start := 1
	limit := id
	if id > b.messageLimit {
		start = id - b.messageLimit
	}
	for i := start; i <= limit; i++ {
		ids = append(ids, i)
	}
	return
}

func (b Bot) getFullMessages(username string, ids []int) (tg.MessageArray, error) {
	if len(ids) == 0 {
		return tg.MessageArray{}, fmt.Errorf("no messages ids")
	}
	b.lg.Infof("fetching messages for ids: %v", ids)
	messages := tg.MessageArray{}

	// No graceful shutdown.
	ctx := context.TODO()

	err := b.mtpClient.Run(ctx, func(ctx context.Context) error {
		// Checking auth status.
		status, err := b.mtpClient.Auth().Status(ctx)
		if err != nil {
			return err
		}
		// Can be already authenticated if we have valid session in
		// session storage.
		if !status.Authorized {
			// Otherwise, perform bot authentication.
			if _, err := b.mtpClient.Auth().Bot(ctx, b.token); err != nil {
				return err
			}
		}

		peerManager := peers.Options{
			// Logger: b.lg,
		}.Build(b.mtpClient.API())

		p, err := peerManager.ResolveDomain(ctx, username)
		if err != nil {
			b.lg.Error(fmt.Sprintf("%v", err))
			return err
		}

		if inputChannel, ok := peer.ToInputChannel(p.InputPeer()); ok {
			IDs := idsToInputMessageClass(ids)

			req := &tg.ChannelsGetMessagesRequest{
				// Channel/supergroup
				Channel: inputChannel, //InputChannelClass
				// IDs of messages to get
				ID: IDs, // []InputMessageClass
			}

			resp, err := b.mtpClient.API().ChannelsGetMessages(ctx, req)
			if err != nil {
				b.lg.Error(fmt.Sprintf("%v", err))
				return err
			}

			var temp interface{} = resp
			messages = temp.(*tg.MessagesChannelMessages).MapMessages().AsMessage()

		} else {
			return fmt.Errorf("not channel")
		}

		// All good, manually authenticated.
		b.lg.Info("Done")

		return nil
	})
	if err != nil {
		return tg.MessageArray{}, err
	}
	return messages, nil
}
