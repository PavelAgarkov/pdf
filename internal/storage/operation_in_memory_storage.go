package storage

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"pdf/internal/adapter"
	"pdf/internal/hash"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"sync"
	"time"
)

type OperationStorage struct {
	sm sync.Map
}

func NewInMemoryOperationStorage() *OperationStorage {
	return &OperationStorage{
		sm: sync.Map{},
	}
}

func (s *OperationStorage) Run(ctx context.Context, tickerTimer time.Duration, adapterLocator *adapter.Locator, loggerFactory *logger.Factory) {
	tickerSetExpired := time.NewTicker(tickerTimer)
	tickerCleaner := time.NewTicker(tickerTimer * 2)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.clearAllFiles(adapterLocator, loggerFactory)
				s.Clear()
			case <-tickerSetExpired.C:
				now := time.Now()
				s.setExpired(now)
			case <-tickerCleaner.C:
				s.clearExpired(adapterLocator, loggerFactory)
			default:
			}
		}
	}()
}

func (s *OperationStorage) setExpired(now time.Time) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			uDur := operation.GetUserData().GetExpiredAt().Unix()
			nDur := now.Unix()

			if nDur > uDur {
				operation.SetStatus(pdf_operation.StatusExpired)
				s.Put(key.(hash.Hash2lvl), operation)
				return true
			}

			return false
		})
}

func (s *OperationStorage) clearExpired(adapterLocator *adapter.Locator, loggerFactory *logger.Factory) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			if operation.CanDeleted() {
				pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
				rootDir := string(pathAdapter.GenerateDirPathToFiles(operation.GetUserData().GetHash2Lvl()))
				err := os.RemoveAll(rootDir)
				if err != nil {
					loggerFactory.PanicLog(fmt.Sprintf("storage: can't remove dir %s", rootDir), zap.Stack("").String)
					return false
				}
				s.Delete(key.(hash.Hash2lvl))
				return true
			}
			return true
		})
}

func (s *OperationStorage) clearAllFiles(adapterLocator *adapter.Locator, loggerFactory *logger.Factory) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
			rootDir := string(pathAdapter.GenerateDirPathToFiles(operation.GetUserData().GetHash2Lvl()))
			err := os.RemoveAll(rootDir)
			if err != nil {
				loggerFactory.PanicLog(fmt.Sprintf("storage: can't clear all dir %s", rootDir), zap.Stack("").String)
				return false
			}
			return true
		})
}

func (s *OperationStorage) Clear() {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			s.Delete(key.(hash.Hash2lvl))
			return true
		})
}

func (s *OperationStorage) Insert(key hash.Hash2lvl, operation pdf_operation.OperationDataInterface) {
	s.sm.Store(key, operation)
}

func (s *OperationStorage) Delete(key hash.Hash2lvl) {
	s.sm.Delete(key)
}

func (s *OperationStorage) Get(key hash.Hash2lvl) (pdf_operation.OperationDataInterface, bool) {
	operation, ok := s.sm.Load(key)

	operat, assert := operation.(pdf_operation.OperationDataInterface)
	if !assert {
		return nil, ok
	}

	return operat.(pdf_operation.OperationDataInterface), ok
}

func (s *OperationStorage) Put(key hash.Hash2lvl, value pdf_operation.OperationDataInterface) (pdf_operation.OperationDataInterface, bool) {
	previousOperation, ok := s.sm.Swap(key, value)

	previousOperation, assert := previousOperation.(pdf_operation.OperationDataInterface)
	if !assert {
		return nil, ok
	}

	return previousOperation.(pdf_operation.OperationDataInterface), ok
}

func (s *OperationStorage) Range(fn func(key interface{}, value interface{}) bool) {
	s.sm.Range(fn)
}
