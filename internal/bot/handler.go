package bot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"emperror.dev/errors"

	"github.com/ChernichenkoStephan/nanostats/internal/stats"
	tele "gopkg.in/telebot.v3"
)

var help string = "\\add @username0 @username1 | add chats to fetch list\n\\del @username0 @username1 delete chats from fetch list\n\\report print report"

func (b Bot) Start() {
	b.botClient.Start()
}

func (b Bot) HandleAddChats(c tele.Context) error {
	resp := "Added:\n"
	names := c.Args()

	if len(names) == 0 {
		err := c.Send(`Empty username list, usage: /add @name0 @name1`)
		if err != nil {
			b.lg.Errorln(err)
		}
	}

	err := c.Send(`Adding`)
	if err != nil {
		b.lg.Errorln(err)
	}

	for _, n := range names {
		err := b.addChat(n)
		if err != nil {
			b.lg.Errorln(err)
		}
		resp += n + "\n"
	}

	b.lg.Infoln(b.repo.GetAll())

	return c.Send(resp)
}

func (b Bot) HandleDeleteChats(c tele.Context) error {
	names := c.Args()
	for _, n := range names {
		b.repo.DeleteWithUsername(n)
	}
	return c.Send(`Done`)
}

func (b Bot) HandleGetStats(c tele.Context) error {
	fstRespErr := c.Send(`Fething...`)

	ss := b.getStats()

	var respond string
	for _, s := range ss {
		respond += fmt.Sprintf("%s\n", s)
	}

	err := b.makeReport(ss)
	if err != nil {
		b.lg.Errorf("got error during report making, got: %v", err)
	}

	secndRespErr := c.Send(respond)

	/*
		type Document struct {
			File

			// (Optional)
			Thumbnail            *Photo `json:"thumb,omitempty"`
			Caption              string `json:"caption,omitempty"`
			MIME                 string `json:"mime_type"`
			FileName             string `json:"file_name,omitempty"`
			DisableTypeDetection bool   `json:"disable_content_type_detection,omitempty"`
		}
	*/

	f := &tele.Document{
		File:     tele.FromDisk(b.outFile),
		FileName: `report.txt`,
	}

	_, fileSendErr := b.botClient.Send(c.Sender(), f)
	if fileSendErr != nil {
		b.lg.Errorln(fileSendErr)
	}

	// return errors.Combine(fstRespErr, secndRespErr)
	return errors.Combine(fstRespErr, secndRespErr, fileSendErr)
}

func (b Bot) HandleStart(c tele.Context) error {
	c.Send("Stats bot says: hi!\n" + help)
	return nil
}

func (b Bot) HandleHelp(c tele.Context) error {
	c.Send(help)
	return nil
}

func (b Bot) getStats() []stats.Stats {
	cc := b.repo.GetAll()

	b.lg.Infof("Fetching: %v\n", cc)

	err := b.withMTPSession(func(ctx context.Context) error {

		for i, c := range cc {

			if c.Username == `` {
				b.lg.Errorln(`empty username`)
				continue
			}

			b.lg.Infof("Shoting: %d, %s, %s\n", c.ID, c.Title, c.Username)

			s, err := b.shot(ctx, c.Username, c.LastPostID)
			if err != nil {
				b.lg.Errorln(err)
				continue
			}

			c.Shots = append(c.Shots, s)

			cc[i] = c

			if i%b.requestsLimit == 0 {
				time.Sleep(time.Duration(b.requestsDelay))
			}
		}

		return nil
	})

	if err != nil {
		b.lg.Error(err)
	}

	return stats.GetStats(cc)
}

func (b Bot) makeReport(ss []stats.Stats) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	exPath := filepath.Dir(ex)
	path := fmt.Sprintf("%s/../%s", exPath, b.outFile)

	err = os.Truncate(path, 100)
	if err != nil {
		return errors.Wrap(err, `got error durung truncation`)
	}

	return stats.OutputStats(ss, path)
}
