package storage

import (
	"context"
	"pdf/internal/hash"
	"pdf/internal/pdf_operation"
	"sync"
	"time"
)

// OperationStorage т.к. теперь у операций есть пользователь, а не наоборот
type OperationStorage struct {
	sm sync.Map
}

func NewInMemoryOperationStorage() *OperationStorage {
	return &OperationStorage{
		sm: sync.Map{},
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
				s.clearExpired()
			default:
			}
		}
	}()
}

func (s *OperationStorage) setExpired(now time.Time) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.Operation)
			uDur := operation.GetBaseOperation().GetUserData().GetExpiredAt().Unix()
			nDur := now.Unix()

			if nDur > uDur {
				operation.GetBaseOperation().SetStatus(pdf_operation.StatusExpired)
				s.Put(key.(hash.Hash2lvl), operation)
				return true
			}

			return false
		})
}

func (s *OperationStorage) clearExpired() {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.Operation)
			bo := operation.GetBaseOperation()

			if bo.CanDeleted() {
				s.Delete(key.(hash.Hash2lvl))
				// тут нужно удалять все файлы для этой операции
				return true
			}

			return false
		})
}

func (s *OperationStorage) ClearAll() {
	s.Clear()
}

func (s *OperationStorage) Clear() {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			s.Delete(key.(hash.Hash2lvl))
			return true
		})
}

func (s *OperationStorage) Insert(key hash.Hash2lvl, operation pdf_operation.Operation) {
	s.sm.Store(key, operation)
}

func (s *OperationStorage) Delete(key hash.Hash2lvl) {
	s.sm.Delete(key)
}

func (s *OperationStorage) Get(key hash.Hash2lvl) (pdf_operation.Operation, bool) {
	operation, ok := s.sm.Load(key)

	ud, assert := operation.(pdf_operation.Operation)
	if !assert {
		return nil, ok
	}

	return ud.(pdf_operation.Operation), ok
}

func (s *OperationStorage) Put(key hash.Hash2lvl, value pdf_operation.Operation) (pdf_operation.Operation, bool) {
	previousOperation, ok := s.sm.Swap(key, value)

	previousOperation, assert := previousOperation.(pdf_operation.Operation)
	if !assert {
		return nil, ok
	}

	return previousOperation.(pdf_operation.Operation), ok
}

func (s *OperationStorage) Range(fn func(key interface{}, value interface{}) bool) {
	s.sm.Range(fn)
}
