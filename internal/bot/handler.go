package bot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"emperror.dev/errors"

	"github.com/ChernichenkoStephan/nanostats/internal/stats"
	tele "gopkg.in/telebot.v3"
)

func (b Bot) Start() {
	b.botClient.Start()
}

func (b Bot) HandleAddChat(c tele.Context) error {
	resp := fmt.Sprintf("c.chat:%v\nc.msg:%v\nc.sndr:%v\n", c.Chat(), c.Message().Chat, c.Sender())
	return c.Send(resp)
}

func (b Bot) HandleDeleteChat(c tele.Context) error {
	return c.Send(`Deleted`)
}

/*
a := &tele.Audio{File: tele.FromDisk("file.ogg")}

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // false

// Will upload the file from disk and send it to the recipient
b.Send(recipient, a)

// Next time you'll be sending this very *Audio, Telebot won't
// re-upload the same file but rather utilize its Telegram FileID
b.Send(otherRecipient, a)

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // true
fmt.Println(a.FileID) // <Telegram file ID>
*/

func (b Bot) HandleGetStats(c tele.Context) error {
	fstRespErr := c.Send(`Fething...`)

	ss := b.getStats()

	var respond string
	for _, s := range ss {
		respond += fmt.Sprintf("%s\n", s)
	}

	b.makeReport(ss)

	secndRespErr := c.Send(respond)

	// f := tele.FromDisk(`out.txt`)
	// fileSendErr := c.Send(f)

	return errors.Combine(fstRespErr, secndRespErr)
	// return errors.Combine(fstRespErr, secndRespErr, fileSendErr)
}

func (b Bot) HandleStart(c tele.Context) error {
	c.Send(`Stats bot says: hi!`)
	return nil
}

func (b Bot) getStats() []stats.Stats {
	cc := b.repo.GetAll()
	for i, c := range cc {

		s, err := b.shot(c.Username, c.LastPostID)
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

	return stats.GetStats(cc)
}

func (b Bot) makeReport(ss []stats.Stats) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	exPath := filepath.Dir(ex)
	path := fmt.Sprintf("%s/%s", exPath, b.outFile)

	err = os.Truncate(path, 100)
	if err != nil {
		return errors.Wrap(err, `got error durung truncation`)
	}

	return stats.OutputStats(ss, path)
}
