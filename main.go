package main

import (
	"net/http"
	"os"

	"github.com/Marcel-MD/distributed-datastore/presentation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	r := presentation.InitRouter()

	log.Fatal().Err(http.ListenAndServe(":8080", r)).Msg("Server failed to start")
}
