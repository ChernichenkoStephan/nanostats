package stats

import (
	"errors"
	"fmt"
	"time"
)

type Stats struct {
	Name           string
	Username       string
	Participants   int
	Delta          int
	DeltaWeek      int
	LastPostViews  int
	AvgViews       float32
	H24PostsAmount int
	LastPostDate   int
}

func (s Stats) String() string {
	t := time.Unix(int64(s.LastPostDate), 0)
	datestring := fmt.Sprintf("%02d.%02d[%02d:%02d]",
		t.Day(), t.Month(), t.Hour(), t.Minute())

	var str string
	str += fmt.Sprintf("%s %s\n", s.Name, s.Username)
	str += fmt.Sprintf("am: %d, dt: %d %d\n", s.Participants, s.Delta, s.DeltaWeek)
	str += fmt.Sprintf("last: %d | avg: %d\n", s.LastPostViews, int(s.AvgViews))
	str += fmt.Sprintf("24/p: %d | last: %s\n", s.H24PostsAmount, datestring)

	return str
}

func (s *Stats) Equals(another *Stats) bool {
	return s.Username == another.Username
}

func last[T any](slice []T) T {
	return slice[len(slice)-1]
}

func GetStats(chats []Chat) (stats []Stats) {
	stats = make([]Stats, 0)
	for _, c := range chats {
		if s, err := getChatStats(c, 20); err == nil {
			stats = append(stats, s)
		}
	}
	return
}

func getChatStats(c Chat, postsAmount int) (Stats, error) {
	if len(c.Shots) == 0 {
		return Stats{}, errors.New(`empty shots list`)
	}

	lastShot := last(c.Shots)
	mm := getMessages(c.Shots, postsAmount)

	avgv := getAvgViews(mm)

	prewShot := lastShot
	if len(c.Shots) > 1 {
		prewShot = c.Shots[len(c.Shots)-2]
	}
	delta := lastShot.Subscribers - prewShot.Subscribers

	// settig just wirst if we have not many shots
	weekShot := c.Shots[0]
	weekAgo := time.Now().Add(time.Hour * -24 * 7)
	for i := len(c.Shots) - 1; i > 0; i-- {
		if c.Shots[i].Created.Unix() < weekAgo.Unix() {
			weekShot = c.Shots[i]
			break
		}
	}
	deltaWeel := lastShot.Subscribers - weekShot.Subscribers

	return Stats{
		Name:           c.Title,
		Username:       c.Username,
		Participants:   last(c.Shots).Subscribers,
		Delta:          delta,
		DeltaWeek:      deltaWeel,
		LastPostViews:  int(last(lastShot.Messages).Views),
		AvgViews:       avgv,
		H24PostsAmount: getAmountPostsFor24H(lastShot.Messages),
		LastPostDate:   last(lastShot.Messages).PostDate,
	}, nil
}

func getMessages(shots []Shot, amount int) []Message {
	var messages []Message
	var i int
	for _, s := range shots {
		for _, m := range s.Messages {
			messages = append(messages, m)
			i++
			if i == amount {
				return messages
			}
		}
	}
	return messages
}

func getAmountPostsFor24H(messages []Message) (amount int) {
	dayAgo := time.Now().Add(time.Hour * -24)
	for i := len(messages) - 1; i > 0; i-- {
		if messages[i].PostDate < int(dayAgo.Unix()) {
			break
		}
		amount++
	}
	return
}

func getAvgViews(messages []Message) float32 {
	var sum int
	for _, m := range messages {
		sum += m.Views
	}
	return float32(sum) / float32(len(messages)+1)
}
