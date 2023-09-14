package storage

import (
	"sync"
)

type Storage interface {
	Get(key any) (value any, ok bool)
	Insert(key, value any)
	Put(key, value any) (previous any, loaded bool)
	Delete(key any)
	Range(fn func(key interface{}, value interface{}) bool)
	Clear()
}

type Memory struct {
	storage sync.Map
}

func NewMemory() Storage {
	return &Memory{}
}

func (s *Memory) Get(key any) (value any, ok bool) {
	return s.storage.Load(key)
}

func (s *Memory) Insert(key, value any) {
	s.storage.Store(key, value)
}

func (s *Memory) Put(key, value any) (previous any, loaded bool) {
	return s.storage.Swap(key, value)
}

func (s *Memory) Delete(key any) {
	s.storage.Delete(key)
}

func (s *Memory) Clear() {
	s.storage.Range(
		func(key interface{}, value interface{}) bool {
			s.storage.Delete(key)
			return true
		})
}

func (s *Memory) Range(fn func(key interface{}, value interface{}) bool) {
	s.storage.Range(fn)
}
