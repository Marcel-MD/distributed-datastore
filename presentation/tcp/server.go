package tcp

import (
	"encoding/json"
	"net"

	"github.com/Marcel-MD/distributed-datastore/domain"
	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/rs/zerolog/log"
)

func ListenAndServe() {
	config := cfg.GetConfig()
	port := config.Current.TcpPort

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Err(err).Msg("Error listening on port " + port)
		return
	}

	log.Info().Msg("Listening on port " + port)

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Err(err).Msg("Error accepting connection")
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	s := domain.GetStore()

	for {
		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			log.Err(err).Msg("Error reading from connection")
			return
		}

		buffer = buffer[:n]

		var a models.Action
		err = json.Unmarshal(buffer, &a)
		if err != nil {
			log.Err(err).Msg("Error unmarshaling action")
			return
		}

		switch a.Command {
		case models.SET:
			err = s.Set(a.Key, a.Value)
			if err != nil {
				log.Err(err).Msg("Error setting value")
			}

		case models.DELETE:
			err = s.Delete(a.Key)
			if err != nil {
				log.Err(err).Msg("Error deleting value")
			}
		}
	}
}
