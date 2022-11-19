package cfg

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/rs/zerolog/log"
)

var c models.Instance
var once sync.Once

func GetCurrentInstance() models.Instance {
	once.Do(func() {
		file, err := os.Open("config/cfg.json")
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening menu.json")
		}
		defer file.Close()

		byteValue, _ := io.ReadAll(file)

		var current models.Instance
		json.Unmarshal(byteValue, &current)

		c = current
	})

	return c
}
