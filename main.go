package main

import (
	"os"

	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/Marcel-MD/distributed-datastore/presentation/http"
	"github.com/Marcel-MD/distributed-datastore/presentation/tcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	config := cfg.GetConfig()
	if config.Current.IsLeader {
		http.ListenAndServe()
	} else {
		tcp.ListenAndServe()
	}
}
