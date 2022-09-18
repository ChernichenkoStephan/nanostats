package bot

import (
	"context"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/tg"
	"github.com/pkg/errors"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	token string

	requestsLimit int
	requestsDelay int
	messageLimit  int

	repo      *tg.IMRepository
	botClient *tele.Bot
	mtpClient *MTPClient

	lg      *zap.SugaredLogger
	outFile string
}

type Options struct {
	Token   string
	AppID   int
	APIHash string

	RequestsLimit int
	RequestsDelay int
	MessageLimit  int

	Repository *tg.IMRepository
	BotClient  *tele.Bot

	Lg      *zap.SugaredLogger
	OutFile string
}

// func NewMTP(token string, appID int, apiHash string, log *zap.Logger) *MTPClient {

func New(opt Options) *Bot {
	mtp := NewMTP(opt.Token, opt.AppID, opt.APIHash, opt.Lg.Desugar())

	b := &Bot{
		token: opt.Token,

		requestsLimit: opt.MessageLimit,
		requestsDelay: opt.RequestsDelay,
		messageLimit:  opt.RequestsLimit,

		repo:      opt.Repository,
		botClient: opt.BotClient,
		mtpClient: mtp,

		lg:      opt.Lg,
		outFile: opt.OutFile,
	}

	var i int
	ctx := context.Background()

	b.mtpClient.withSession(ctx, func(ctx context.Context) error {
		for _, c := range opt.Repository.GetAll() {

			err := b.update(ctx, &c, tg.MAMUAL)
			if err != nil {
				b.lg.Errorln(err)
				continue
			}

			b.repo.Set(c)

			if i%opt.RequestsLimit == 0 {
				time.Sleep(time.Duration(opt.RequestsDelay))
			}
			i++
		}
		return nil
	})

	return b
}

func (b Bot) addChat(username string) error {
	b.lg.Infof("Fetching %v\n", username)
	cf, err := b.botClient.ChatByUsername(username)
	if err != nil {
		return errors.Wrap(err, `error during fetching chat to add`)
	}
	b.lg.Infoln(cf)

	title := cf.FirstName
	if cf.Title != `` {
		title = cf.Title
	}

	c := tg.Chat{
		Title:    title,
		TgID:     cf.ID,
		Username: username,
		Type:     string(cf.Type),
	}

	b.repo.Set(c)

	b.lg.Infof("Added: %v", b.repo.GetAll())

	return nil
}
