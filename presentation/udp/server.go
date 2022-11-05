package udp

import (
	"math/rand"
	"net"
	"time"

	"github.com/Marcel-MD/distributed-datastore/domain"
	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

func ListenAndServe() {
	s := domain.GetStore()
	config := cfg.GetConfig()
	port := config.Current.UdpPort

	addr, err := net.ResolveUDPAddr("udp4", ":"+port)
	if err != nil {
		log.Err(err).Msg("Error resolving udp address")
		return
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Err(err).Msgf("Error listening on port %s", port)
		return
	}

	log.Info().Msg("Listening on port " + port)

	defer conn.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Err(err).Msg("Error reading from udp")
			continue
		}

		key := string(buffer[:n])

		value, err := s.Get(key)
		if err != nil {
			log.Err(err).Msg("Error getting value from store")
			conn.WriteToUDP([]byte("Key not found"), addr)
			continue
		}

		_, err = conn.WriteToUDP(value, addr)
		if err != nil {
			log.Err(err).Msg("Error writing to udp")
			continue
		}
	}
}
