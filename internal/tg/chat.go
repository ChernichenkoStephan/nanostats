package tg

type Chat struct {
	ID       int64
	TgID     int64
	Username string
	Title    string
	Type     string
	Shots    []Shot
	Messages []Message
}

func (c Chat) Participants() int {
	if len(c.Shots) == 0 {
		return 0
	}
	return c.Shots[len(c.Shots)-1].Amount
}
