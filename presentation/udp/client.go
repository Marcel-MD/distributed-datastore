package udp

import (
	"errors"
	"net"
	"sync"

	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Get(key string) ([]byte, error)
}

type client struct {
	connections map[string]*net.UDPConn
}

func (c *client) Get(key string) ([]byte, error) {
	for _, conn := range c.connections {
		conn.Write([]byte(key))
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return nil, err
		}
		return buffer[0:n], nil
	}

	return nil, errors.New("no connection available")
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
		connections: make(map[string]*net.UDPConn),
	}

	for _, instance := range config.Instances {
		address := instance.Host + ":" + instance.UdpPort

		s, err := net.ResolveUDPAddr("udp4", address)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to resolve address %s", address)
			continue
		}

		conn, err := net.DialUDP("udp4", nil, s)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to connect to %s", address)
			continue
		}

		client.connections[instance.Host] = conn
	}

	c = &client
}
