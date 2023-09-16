package storage

import (
	"context"
	"sync"
	"time"
)

const salt = "1af1dfa857bf1d8814fe1af898 3c18080019922e557f15a8a"
const Timer = 5 * time.Minute

type UserData struct {
	hash1lvl  Hash1lvl // это и будет ключ для записи в куки
	hash2lvl  Hash2lvl // это и будет ключ для записи в основное хранилище и для ссылки
	files     sync.Map
	expiredAt time.Time
}

type UserStorage struct {
	storage Storage
}

func NewUserData(hash1lvl Hash1lvl, hash2lvl Hash2lvl, files []string, expiredAt time.Time) *UserData {
	ud := &UserData{
		hash1lvl:  hash1lvl,
		hash2lvl:  hash2lvl,
		files:     sync.Map{},
		expiredAt: expiredAt,
	}

	for key, filename := range files {
		ud.files.Store(key, filename)
	}

	return ud
}

func NewInMemoryUserStorage() *UserStorage {
	return &UserStorage{
		storage: NewMemory(),
	}
}

func (s *UserStorage) Run(ctx context.Context, tickerTimer time.Duration) {
	ticker := time.NewTicker(tickerTimer)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				s.clearExpired(now)
			default:
			}
		}
	}()
}

func (s *UserStorage) clearExpired(now time.Time) {
	s.storage.Range(
		func(key interface{}, value interface{}) bool {
			userData := value.(*UserData)
			uDur := userData.expiredAt.Unix()
			nDur := now.Unix()

			if nDur > uDur {
				s.Delete(key.(string))
				// отсюда нужно будет удалять файлы и архивы, но их пока нет
			}

			return true
		})
}

func (s *UserStorage) ClearAll() {
	s.storage.Clear()
}

func (s *UserStorage) Insert(key string, ud *UserData) {
	s.storage.Insert(key, ud)
}

func (s *UserStorage) Delete(key string) {
	s.storage.Delete(key)
}

func (s *UserStorage) Get(key string) (value *UserData, ok bool) {
	ud, ok := s.storage.Get(key)

	ud, assert := ud.(*UserData)
	if !assert {
		return nil, ok
	}

	return ud.(*UserData), ok
}

func (s *UserStorage) Put(key string, value *UserData) (previous *UserData, loaded bool) {
	previousUD, ok := s.storage.Put(key, value)

	previousUD, assert := previousUD.(*UserData)
	if !assert {
		return nil, ok
	}

	return previousUD.(*UserData), ok
}
