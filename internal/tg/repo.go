package tg

import (
	"sort"
	"sync"
)

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

	res := make([]Chat, 0)
	for id, c := range repo.chats {
		if id == 0 {
			continue
		}
		res = append(res, c)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

func (repo *IMRepository) Get(ID int64) (Chat, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	c, ok := repo.chats[ID]
	return c, ok
}

func (repo *IMRepository) GetWithTgID(id int64) (Chat, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, c := range repo.chats {
		if c.TgID == id {
			return c, true
		}
	}
	return Chat{}, false
}

func (repo *IMRepository) GetWithUsername(username string) (Chat, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, c := range repo.chats {
		if c.Username == username {
			return c, true
		}
	}
	return Chat{}, false
}

func (repo *IMRepository) Set(c Chat) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if c.ID == 0 {
		c.ID = int64(len(repo.chats) + 1)
	}
	repo.chats[c.ID] = c
}

func (repo *IMRepository) Delete(c Chat) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.chats, c.ID)
}

func (repo *IMRepository) DeleteWithUsername(username string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, c := range repo.chats {
		if c.Username == username {
			repo.Delete(c)
		}
	}
}
