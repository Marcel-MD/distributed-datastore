package models

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

type Instance struct {
	IsLeader bool   `json:"is_leader"`
	Host     string `json:"host"`
	HttpPort string `json:"http_port"`
	TcpPort  string `json:"tcp_port"`
}

var c Instance
var once sync.Once

func GetCurrentInstance() Instance {
	once.Do(func() {
		file, err := os.Open("config/cfg.json")
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening menu.json")
		}
		defer file.Close()

		byteValue, _ := io.ReadAll(file)

		var current Instance
		json.Unmarshal(byteValue, &current)

		c = current
	})

	return c
}
