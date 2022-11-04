package domain

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type store struct {
	data map[string][]byte
}

var once sync.Once
var instance Store

func GetStore() Store {

	once.Do(func() {
		log.Info().Msg("Initializing store")

		instance = &store{
			data: make(map[string][]byte),
		}
	})

	return instance
}

func (s *store) Get(key string) ([]byte, error) {
	log.Info().Msgf("Getting key %s", key)

	if value, ok := s.data[key]; ok {
		return value, nil
	}

	return nil, errors.New("key not found")
}

func (s *store) Set(key string, value []byte) error {
	log.Info().Msgf("Setting key %s", key)

	s.data[key] = value
	return nil
}

func (s *store) Delete(key string) error {
	log.Info().Msgf("Deleting key %s", key)

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return nil
	}

	return errors.New("key not found")
}
