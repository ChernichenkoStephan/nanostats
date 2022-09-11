package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ChernichenkoStephan/nanostats/internal/bot"
	"github.com/ChernichenkoStephan/nanostats/internal/stats"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.mux.Unlock()
	return nil
}

func initBot(cfg Config, lg *zap.SugaredLogger, repo *stats.IMRepository) (*bot.Bot, error) {

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botClient, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("error during bot init, got %w", err)
	}

	// Using custom session storage.
	// You can save session to database, e.g. Redis, MongoDB or postgres.
	// See memorySession for implementation details.
	sessionStorage := &memorySession{}

	mtpClient := telegram.NewClient(cfg.AppID, cfg.APIHash, telegram.Options{
		SessionStorage: sessionStorage,
		Logger:         lg.Desugar(),
	})

	opts := bot.Options{
		Token: cfg.Token,

		RequestsLimit: cfg.RequestLimit,
		RequestsDelay: cfg.RequestDelay,
		MessageLimit:  cfg.MessagesLimit,

		Repository: repo,
		BotClient:  botClient,
		MTPClient:  mtpClient,

		Lg:      lg,
		OutFile: cfg.OutputFileName,
	}

	b := bot.New(opts)

	var (
		add = `/add`
		del = `/delete`
		rep = `/report`
	)

	botClient.Use(bot.Logging(lg))

	botClient.Handle(`/start`, b.HandleStart)
	botClient.Handle(add, b.HandleAddChat)
	botClient.Handle(del, b.HandleDeleteChat)
	botClient.Handle(rep, b.HandleGetStats)

	commands := []tele.Command{
		{
			Text:        add,
			Description: `Adds chat to stats fetching`,
		},
		{
			Text:        del,
			Description: `Deletes chat from stats fetching`,
		},
		{
			Text:        rep,
			Description: `Returns chats stats`,
		},
	}

	err = botClient.SetCommands(commands)
	if err != nil {
		lg.Errorln("Command setup failed (on set)")
	}

	botClient.Use(b.LastIDUpdating())

	return b, nil
}
