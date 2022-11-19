package websocket

import (
	"net/http"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Serve(w http.ResponseWriter, r *http.Request) error {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade websocket connection")
		return err
	}

	send := make(chan models.Action)
	c := &connection{ws: ws, send: send}

	go c.readPump()
	go c.writePump()

	channels[ws] = send

	return nil
}

var channels = make(map[*websocket.Conn]chan models.Action)

func Broadcast(action models.Action) {
	for _, ch := range channels {
		ch <- action
	}
}
