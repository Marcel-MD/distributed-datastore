package cfg

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Current   Instance   `json:"current"`
	Instances []Instance `json:"instances"`
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
		json.Unmarshal(byteValue, &c)
	})

	return c
}
