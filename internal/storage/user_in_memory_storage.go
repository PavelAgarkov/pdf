package storage

import (
	"context"
	"time"
)

const (
	salt  = "1af1dfa857bf1d8814fe1af898 3c18080019922e557f15a8a"
	Timer = 5 * time.Minute
)

// OperationStorage т.к. теперь у операций есть пользователь, а не наоборот
type UserStorage struct {
	storage Storage
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
				s.Delete(key.(Hash2lvl))
				// отсюда нужно будет удалять файлы и архивы, но их пока нет
			}

			return true
		})
}

func (s *UserStorage) ClearAll() {
	s.storage.Clear()
}

func (s *UserStorage) Insert(key Hash2lvl, ud *UserData) {
	s.storage.Insert(key, ud)
}

func (s *UserStorage) Delete(key Hash2lvl) {
	s.storage.Delete(key)
}

func (s *UserStorage) Get(key Hash2lvl) (value *UserData, ok bool) {
	ud, ok := s.storage.Get(key)

	ud, assert := ud.(*UserData)
	if !assert {
		return nil, ok
	}

	return ud.(*UserData), ok
}

func (s *UserStorage) Put(key Hash2lvl, value *UserData) (previous *UserData, loaded bool) {
	previousUD, ok := s.storage.Put(key, value)

	previousUD, assert := previousUD.(*UserData)
	if !assert {
		return nil, ok
	}

	return previousUD.(*UserData), ok
}
