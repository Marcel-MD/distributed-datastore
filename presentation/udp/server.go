package udp

import (
	"encoding/json"
	"net"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/rs/zerolog/log"
)

const UdpPort = ":1053"
const UdpHost = "255.255.255.255"

func ListenAndServe() {

	addr, err := net.ResolveUDPAddr("udp4", UdpHost+UdpPort)
	if err != nil {
		log.Err(err).Msg("Error resolving UDP address")
		return
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Err(err).Msg("Error listening on UDP")
		return
	}
	defer conn.Close()

	client := GetClient()

	for {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Err(err).Msg("Error reading from UDP")
			continue
		}

		var instance models.Instance
		err = json.Unmarshal(buf[:n], &instance)
		if err != nil {
			log.Err(err).Msg("Error unmarshaling instance")
			continue
		}

		if client.HasInstance(instance.Host) {
			continue
		}

		client.AddInstance(instance)
	}
}
