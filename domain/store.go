package domain

import (
	"errors"
	"sync"
)

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type store struct {
	data map[string][]byte
}

var once = sync.Once{}
var instance Store

func NewStore() Store {

	once.Do(func() {
		instance = &store{
			data: make(map[string][]byte),
		}
	})

	return instance
}

func (s *store) Get(key string) ([]byte, error) {

	if value, ok := s.data[key]; ok {
		return value, nil
	}

	return nil, errors.New("key not found")
}

func (s *store) Set(key string, value []byte) error {
	s.data[key] = value
	return nil
}

func (s *store) Delete(key string) error {

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return nil
	}

	return errors.New("key not found")
}
