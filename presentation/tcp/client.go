package tcp

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Set(key string, value []byte) error
	Delete(key string) error
}

type client struct {
	connections map[string]net.Conn
}

func (c *client) Set(key string, value []byte) error {
	action := models.Action{
		Command: models.SET,
		Key:     key,
		Value:   value,
	}

	return c.broadcast(action)
}

func (c *client) Delete(key string) error {
	action := models.Action{
		Command: models.DELETE,
		Key:     key,
	}

	return c.broadcast(action)
}

func (c *client) broadcast(action models.Action) error {
	data, err := json.Marshal(action)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal action")
		return err
	}

	for host, conn := range c.connections {
		_, err := conn.Write(data)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to write to connection %s", host)
			return err
		}
	}

	return nil
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
