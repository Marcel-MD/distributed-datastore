package tcp

import (
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/Marcel-MD/distributed-datastore/domain"
	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/Marcel-MD/distributed-datastore/presentation/websocket"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Update(key string, value []byte) error
	Delete(key string) error

	AddConnection(instance models.Instance)
	RemoveConnection(host string)
	SyncData()
}

type client struct {
	connections map[string]net.Conn
}

func (c *client) Get(key string) ([]byte, error) {
	action := models.Action{
		Command: models.GET,
		Key:     key,
	}

	data, err := c.broadcast(action)

	return data, err
}

func (c *client) Set(key string, value []byte) error {
	action := models.Action{
		Command: models.SET,
		Key:     key,
		Value:   value,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) Update(key string, value []byte) error {
	action := models.Action{
		Command: models.UPDATE,
		Key:     key,
		Value:   value,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) Delete(key string) error {
	action := models.Action{
		Command: models.DELETE,
		Key:     key,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) broadcast(action models.Action) ([]byte, error) {
	data, err := json.Marshal(action)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal action")
		return nil, err
	}

	websocket.Broadcast(action)

	for _, conn := range c.connections {
		_, err := conn.Write(data)
		if err != nil {
			continue
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			continue
		}

		if string(buffer[0:n]) == models.ERROR {
			continue
		}

		return buffer[0:n], nil
	}

	return nil, errors.New("failed to broadcast")
}

func (c *client) AddConnection(instance models.Instance) {

	if _, ok := c.connections[instance.Host]; ok {
		return
	}

	address := instance.Host + instance.TcpPort
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Err(err).Msgf("Error connecting to %s", address)
		return
	}

	c.connections[instance.Host] = conn
}

func (c *client) RemoveConnection(host string) {
	conn, ok := c.connections[host]
	if !ok {
		log.Error().Msgf("Connection to %s not found", host)
		return
	}

	conn.Close()
	delete(c.connections, host)
}

var c Client
var once sync.Once

func GetClient() Client {
	once.Do(func() {
		c = &client{
			connections: make(map[string]net.Conn),
		}
	})

	return c
}

func (c *client) SyncData() {
	s := domain.GetStore()
	keys := s.GetKeys()

	for _, key := range keys {
		value, err := s.Get(key)
		if err != nil {
			continue
		}

		c.Set(key, value)
	}
}
