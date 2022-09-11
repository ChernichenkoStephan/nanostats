package bot

import (
	"fmt"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/stats"

	mtp "github.com/gotd/td/telegram"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	token string

	requestsLimit int
	requestsDelay int
	messageLimit  int

	repo      *stats.IMRepository
	botClient *tele.Bot
	mtpClient *mtp.Client

	lg      *zap.SugaredLogger
	outFile string
}

type Options struct {
	Token string

	RequestsLimit int
	RequestsDelay int
	MessageLimit  int

	Repository *stats.IMRepository
	BotClient  *tele.Bot
	MTPClient  *mtp.Client

	Lg      *zap.SugaredLogger
	OutFile string
}

func New(opt Options) *Bot {
	var i int
	for _, c := range opt.Repository.GetAll() {
		opt.Lg.Infof("Fetching %v with %d", c.Username, c.LastPostID)
		cf, err := opt.BotClient.ChatByUsername(c.Username)
		if err != nil {
			opt.Lg.Errorln(err)
			continue
		}
		c.TgID = cf.ID
		c.Title = cf.FirstName
		if cf.Title != `` {
			c.Title = cf.Title
		}
		c.Type = fmt.Sprintf("%v", cf.Type)
		opt.Repository.Set(c)

		if i%opt.RequestsLimit == 0 {
			time.Sleep(time.Duration(opt.RequestsDelay))
		}
		i++
	}
	return &Bot{
		token: opt.Token,

		requestsLimit: opt.MessageLimit,
		requestsDelay: opt.RequestsDelay,
		messageLimit:  opt.RequestsLimit,

		repo:      opt.Repository,
		botClient: opt.BotClient,
		mtpClient: opt.MTPClient,

		lg:      opt.Lg,
		outFile: opt.OutFile,
	}
}
