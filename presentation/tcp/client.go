package tcp

import (
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Update(key string, value []byte) error
	Delete(key string) error
}

type client struct {
	connections map[string]net.Conn
}

func (c *client) Get(key string) ([]byte, error) {
	action := Action{
		Command: GET,
		Key:     key,
	}

	data, err := c.broadcast(action)

	return data, err
}

func (c *client) Set(key string, value []byte) error {
	action := Action{
		Command: SET,
		Key:     key,
		Value:   value,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) Update(key string, value []byte) error {
	action := Action{
		Command: UPDATE,
		Key:     key,
		Value:   value,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) Delete(key string) error {
	action := Action{
		Command: DELETE,
		Key:     key,
	}

	_, err := c.broadcast(action)

	return err
}

func (c *client) broadcast(action Action) ([]byte, error) {
	data, err := json.Marshal(action)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal action")
		return nil, err
	}

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

		if string(buffer[0:n]) == ERROR {
			continue
		}

		return buffer[0:n], nil
	}

	return nil, errors.New("failed to broadcast")
}

var c Client
var once sync.Once

func GetClient() Client {
	once.Do(func() {
		connect()
	})

	return c
}

func connect() {
	config := cfg.GetConfig()

	client := client{
		connections: make(map[string]net.Conn),
	}

	for _, instance := range config.Instances {
		address := instance.Host + ":" + instance.TcpPort
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Err(err).Msgf("Error connecting to %s", address)
			continue
		}

		client.connections[instance.Host] = conn
	}

	c = &client
}
