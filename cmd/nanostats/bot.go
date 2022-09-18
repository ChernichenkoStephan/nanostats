package main

import (
	"fmt"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/bot"
	"github.com/ChernichenkoStephan/nanostats/internal/tg"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func initBot(cfg Config, lg *zap.SugaredLogger, repo *tg.IMRepository) (*bot.Bot, error) {

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botClient, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("error during bot init, got %w", err)
	}

	opts := bot.Options{
		Token:   cfg.Token,
		AppID:   cfg.AppID,
		APIHash: cfg.APIHash,

		RequestsLimit: cfg.RequestLimit,
		RequestsDelay: cfg.RequestDelay,
		MessageLimit:  cfg.MessagesLimit,

		Repository: repo,
		BotClient:  botClient,

		Lg:      lg,
		OutFile: cfg.OutputFileName,
	}

	b := bot.New(opts)

	var (
		start  = `/start`
		add    = `/add`
		del    = `/del`
		shot   = `/shot`
		report = `/report`
		help   = `/help`
	)

	botClient.Use(bot.Logging(lg))

	botClient.Handle(start, b.HandleStart)
	botClient.Handle(add, b.HandleAddChats)
	botClient.Handle(del, b.HandleDeleteChats)
	botClient.Handle(shot, b.HandleShot)
	botClient.Handle(report, b.HandleReport)
	botClient.Handle(help, b.HandleHelp)

	commands := []tele.Command{
		{
			Text:        help,
			Description: `Prints commands list with usages`,
		},
		{
			Text:        shot,
			Description: `Returns chats stats`,
		},
		{
			Text:        report,
			Description: `Returns chats stats in file`,
		},
		{
			Text:        add,
			Description: `Adds chat to list`,
		},
		{
			Text:        del,
			Description: `Deletes chat from list`,
		},
	}

	err = botClient.SetCommands(commands)
	if err != nil {
		lg.Errorln("Command setup failed (on set)")
	}

	botClient.Use(b.LastIDUpdating())

	return b, nil
}
