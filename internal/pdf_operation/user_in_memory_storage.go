package pdf_operation

import (
	"context"
	storage2 "pdf/internal/storage"
	"time"
)

const (
	Timer5  = 5 * time.Minute
	Timer10 = 10 * time.Minute
	Timer15 = 15 * time.Minute
)

// OperationStorage т.к. теперь у операций есть пользователь, а не наоборот
type OperationStorage struct {
	storage storage2.Storage
}

func NewInMemoryOperationStorage() *OperationStorage {
	return &OperationStorage{
		storage: storage2.NewMemory(),
	}
}

func (s *OperationStorage) Run(ctx context.Context, tickerTimer time.Duration) {
	tickerSetExpired := time.NewTicker(tickerTimer)
	tickerCleaner := time.NewTicker(tickerTimer * 2)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ClearAll()
				// удалить все файлы по операциям
			case <-tickerSetExpired.C:
				now := time.Now()
				s.setExpired(now)
			case <-tickerCleaner.C:
				s.clear()
			default:
			}
		}
	}()
}

func (s *OperationStorage) setExpired(now time.Time) {
	s.storage.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(Operation)
			uDur := operation.GetBaseOperation().GetUserData().GetExpiredAt().Unix()
			nDur := now.Unix()

			if nDur > uDur {
				operation.GetBaseOperation().SetStatus(StatusExpired)
				s.Put(key.(storage2.Hash2lvl), operation)
				return true
			}

			return false
		})
}

func (s *OperationStorage) clear() {
	s.storage.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(Operation)
			bo := operation.GetBaseOperation()

			if bo.CanDeleted() {
				s.Delete(key.(storage2.Hash2lvl))
				// тут нужно удалять все файлы для этой операции
				return true
			}

			return false
		})
}

func (s *OperationStorage) ClearAll() {
	s.storage.Clear()
}

func (s *OperationStorage) Insert(key storage2.Hash2lvl, operation Operation) {
	s.storage.Insert(key, operation)
}

func (s *OperationStorage) Delete(key storage2.Hash2lvl) {
	s.storage.Delete(key)
}

func (s *OperationStorage) Get(key storage2.Hash2lvl) (Operation, bool) {
	operation, ok := s.storage.Get(key)

	ud, assert := operation.(Operation)
	if !assert {
		return nil, ok
	}

	return ud.(Operation), ok
}

func (s *OperationStorage) Put(key storage2.Hash2lvl, value Operation) (Operation, bool) {
	previousOperation, ok := s.storage.Put(key, value)

	previousOperation, assert := previousOperation.(Operation)
	if !assert {
		return nil, ok
	}

	return previousOperation.(Operation), ok
}
