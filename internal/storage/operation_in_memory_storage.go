package storage

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"strconv"
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

func (s *OperationStorage) Run(
	ctx context.Context,
	tickerTimer time.Duration,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) {
	tickerSetExpired := time.NewTicker(tickerTimer)
	tickerCleaner := time.NewTicker(tickerTimer)
	go func(tickerSetExpired, tickerCleaner *time.Ticker) {
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("storage RUN : Recovered. Panic: %s\n", r)
				loggerFactory.ErrorLog(errStr, "")
			}
		}()
		defer func() {
			tickerSetExpired.Stop()
			tickerCleaner.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerSetExpired.C:
				s.setExpired(time.Now())
			case <-tickerCleaner.C:
				s.clearExpired(adapterLocator, loggerFactory)
			}
		}
	}(tickerSetExpired, tickerCleaner)
}

func (s *OperationStorage) setExpired(now time.Time) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			uDur := operation.GetUserData().GetExpiredAt().Unix()
			nDur := now.Unix()

			if nDur > uDur {
				operation.SetStatus(internal.StatusExpired)
				s.Put(key.(internal.Hash2lvl), operation)
				return true
			}

			return false
		})
}

func (s *OperationStorage) clearExpired(adapterLocator *locator.Locator, loggerFactory *logger.Factory) {
	allInMemory := 0
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			allInMemory++
			if operation.CanDeleted() {
				pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
				rootDir := string(pathAdapter.GenerateRootDir(operation.GetUserData().GetHash2Lvl()))
				err := os.RemoveAll(rootDir)
				if err != nil {
					loggerFactory.PanicLog(fmt.Sprintf("storage: can't remove dir %s", rootDir), zap.Stack("").String)
					return false
				}
				s.Delete(key.(internal.Hash2lvl))
				return true
			}
			return true
		})
	if allInMemory > 0 {
		loggerFactory.InfoLog(fmt.Sprintf("key in storage %s", strconv.Itoa(allInMemory)))
	}
}

func (s *OperationStorage) clearAllFiles(adapterLocator *locator.Locator, loggerFactory *logger.Factory) {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			operation := value.(pdf_operation.OperationDataInterface)
			pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
			rootDir := string(pathAdapter.GenerateRootDir(operation.GetUserData().GetHash2Lvl()))
			err := os.RemoveAll(rootDir)
			if err != nil {
				loggerFactory.PanicLog(fmt.Sprintf("storage: can't clear all dir %s", rootDir), zap.Stack("").String)
				return false
			}
			return true
		})
}

func (s *OperationStorage) clear() {
	s.sm.Range(
		func(key interface{}, value interface{}) bool {
			s.Delete(key.(internal.Hash2lvl))
			return true
		})
}

func (s *OperationStorage) ClearStorageAndFilesystem(
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) {
	s.clearAllFiles(adapterLocator, loggerFactory)
	s.clear()
}

func (s *OperationStorage) Insert(key internal.Hash2lvl, operation pdf_operation.OperationDataInterface) {
	s.sm.Store(key, operation)
}

func (s *OperationStorage) Delete(key internal.Hash2lvl) {
	s.sm.Delete(key)
}

func (s *OperationStorage) Get(key internal.Hash2lvl) (pdf_operation.OperationDataInterface, bool) {
	operation, ok := s.sm.Load(key)

	operat, assert := operation.(pdf_operation.OperationDataInterface)
	if !assert {
		return nil, ok
	}

	return operat.(pdf_operation.OperationDataInterface), ok
}

func (s *OperationStorage) Put(key internal.Hash2lvl, value pdf_operation.OperationDataInterface) (pdf_operation.OperationDataInterface, bool) {
	previousOperation, ok := s.sm.Swap(key, value)

	previousOperation, assert := previousOperation.(pdf_operation.OperationDataInterface)
	if !assert {
		return nil, ok
	}

	return previousOperation.(pdf_operation.OperationDataInterface), ok
}
