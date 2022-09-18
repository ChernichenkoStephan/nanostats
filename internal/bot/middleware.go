package bot

import (
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/tg"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func Logging(logger *zap.SugaredLogger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Infof("Message from '%v' with text '%v' from name %s in [%s|%d]", c.Sender().ID, c.Text(), c.Chat().FirstName, c.Chat().Type, c.Chat().ID)
			return next(c)
		}
	}
}

func (b *Bot) Auth() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			// check id
			return next(c)
		}
	}
}

func (b *Bot) LastIDUpdating() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if chat, ok := b.repo.GetWithTgID(c.Chat().ID); ok {
				m := tg.Message{
					ID:      c.Message().ID,
					Text:    c.Message().Text,
					Created: time.Now(),
				}
				chat.Messages = append(chat.Messages, m)
				b.repo.Set(chat)
			}
			return next(c)
		}
	}
}
