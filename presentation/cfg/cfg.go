package cfg

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Current   Instance
	Instances map[string]Instance
}

type Instance struct {
	IsLeader bool   `json:"is_leader"`
	Host     string `json:"host"`
	HttpPort string `json:"http_port"`
	TcpPort  string `json:"tcp_port"`
}

var c Config
var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		file, err := os.Open("config/cfg.json")
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening menu.json")
		}
		defer file.Close()

		byteValue, _ := io.ReadAll(file)

		var current Instance
		json.Unmarshal(byteValue, &current)

		c.Current = current
		c.Instances = make(map[string]Instance)
	})

	return c
}

func (c *Config) AddInstance(instance Instance) {
	if c.Current.Host == instance.Host {
		return
	}

	c.Instances[instance.Host] = instance
}

func (c *Config) RemoveInstance(instance Instance) {
	if c.Current.Host == instance.Host {
		return
	}

	delete(c.Instances, instance.Host)
}
