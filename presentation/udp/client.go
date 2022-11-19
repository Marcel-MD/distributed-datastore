package udp

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/Marcel-MD/distributed-datastore/presentation/tcp"
	"github.com/rs/zerolog/log"
)

type Client interface {
	BroadcastConfig()
	HasInstance(host string) bool
	AddInstance(instance models.Instance)
}

type client struct {
	addr            *net.UDPAddr
	currentInstance models.Instance
	aliveInstances  map[string]int64
	tcpClient       tcp.Client
}

var (
	c    Client
	once sync.Once
)

func GetClient() Client {
	once.Do(func() {

		addr, err := net.ResolveUDPAddr("udp4", UdpHost+UdpPort)

		if err != nil {
			log.Err(err).Msg("Error resolving UDP address")
			return
		}

		c = &client{
			addr:            addr,
			currentInstance: models.GetCurrentInstance(),
			aliveInstances:  make(map[string]int64),
			tcpClient:       tcp.GetClient(),
		}
	})

	return c
}

func (c *client) BroadcastConfig() {

	for {
		time.Sleep(5 * time.Second)

		conn, err := net.DialUDP("udp", nil, c.addr)
		if err != nil {
			log.Err(err).Msg("Error dialing UDP")
			return
		}
		defer conn.Close()

		body, err := json.Marshal(c.currentInstance)
		if err != nil {
			log.Err(err).Msg("Error marshaling config")
			return
		}

		_, err = conn.Write(body)
		if err != nil {
			log.Err(err).Msg("Error writing to UDP")
			return
		}

		c.checkInstances()
	}
}

func (c *client) checkInstances() {
	for host, pingTime := range c.aliveInstances {
		if time.Now().Unix()-pingTime > 10 {
			delete(c.aliveInstances, host)
			c.tcpClient.RemoveConnection(host)

			log.Info().Msgf("Instance %s is dead", host)
		}
	}
}

func (c *client) HasInstance(host string) bool {
	if c.currentInstance.Host == host {
		return true
	}

	if _, ok := c.aliveInstances[host]; ok {
		c.aliveInstances[host] = time.Now().Unix()

		log.Debug().Msgf("Instance %s is alive", host)

		return true
	}

	return false
}

func (c *client) AddInstance(instance models.Instance) {
	if c.currentInstance.Host == instance.Host {
		return
	}

	c.aliveInstances[instance.Host] = time.Now().Unix()
	c.tcpClient.AddConnection(instance)

	log.Info().Msgf("Add new instance %s", instance.Host)
}
