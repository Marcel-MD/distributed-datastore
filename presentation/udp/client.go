package udp

import (
	"encoding/json"
	"net"

	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

func BroadcastConfig() {
	addr, err := net.ResolveUDPAddr("udp4", UdpHost+UdpPort)

	if err != nil {
		log.Err(err).Msg("Error resolving UDP address")
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Err(err).Msg("Error dialing UDP")
		return
	}
	defer conn.Close()

	config := cfg.GetConfig()

	body, err := json.Marshal(config.Current)
	if err != nil {
		log.Err(err).Msg("Error marshaling config")
		return
	}

	_, err = conn.Write(body)
	if err != nil {
		log.Err(err).Msg("Error writing to UDP")
		return
	}
}
