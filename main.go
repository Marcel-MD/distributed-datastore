package main

import (
	"os"

	"github.com/Marcel-MD/distributed-datastore/presentation/http"
	"github.com/Marcel-MD/distributed-datastore/presentation/tcp"
	"github.com/Marcel-MD/distributed-datastore/presentation/udp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	go udp.ListenAndServe()

	client := udp.GetClient()
	go client.BroadcastConfig()

	go tcp.ListenAndServe()

	http.ListenAndServe()
}
