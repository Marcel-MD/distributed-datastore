package main

import (
	"os"
	"time"

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

	time.Sleep(1 * time.Second)

	udp.BroadcastConfig()

	time.Sleep(1 * time.Second)

	go tcp.ListenAndServe()

	time.Sleep(1 * time.Second)

	http.ListenAndServe()
}
