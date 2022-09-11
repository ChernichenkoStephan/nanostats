package stats

import "sync"

type IMRepository struct {
	mu    sync.Mutex
	chats map[int64]Chat
}

func NewRepository() *IMRepository {
	return &IMRepository{
		chats: make(map[int64]Chat),
	}
}

func (repo *IMRepository) GetAll() []Chat {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	res := make([]Chat, len(repo.chats)-1)
	for _, c := range repo.chats {
		res = append(res, c)
	}
	return res
}

func (repo *IMRepository) Get(ID int64) (Chat, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	c, ok := repo.chats[ID]
	return c, ok
}

func (repo *IMRepository) GetWithTgID(tgID int64) (Chat, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, c := range repo.chats {
		if c.TgID == tgID {
			return c, true
		}
	}
	return Chat{}, false
}

func (repo *IMRepository) Set(c Chat) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.chats[c.ID] = c
}

func (repo *IMRepository) Delete(c Chat) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.chats, c.ID)
}

func (repo *IMRepository) DeleteWithTgID(tgID int64) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, c := range repo.chats {
		if c.TgID == tgID {
			repo.Delete(c)
		}
	}
}
