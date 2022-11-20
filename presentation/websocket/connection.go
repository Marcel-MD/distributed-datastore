package websocket

import (
	"time"

	"github.com/Marcel-MD/distributed-datastore/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type connection struct {
	ws   *websocket.Conn
	send chan models.Action
}

func (c *connection) readPump() {
	log.Debug().Msg("Starting websocket read pump")

	defer func() {
		log.Debug().Msg("Stopping websocket read pump")
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var action models.Action

		err := c.ws.ReadJSON(&action)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Debug().Msg("Normal websocket close")
				break
			}
			log.Err(err).Msg("Unexpected websocket close")
			break
		}
	}
}

func (c *connection) writePump() {
	log.Debug().Msg("Starting websocket write pump")

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Debug().Msg("Stopping websocket write pump")
		ticker.Stop()
		c.ws.Close()
		delete(channels, c.ws)
	}()

	for {
		select {
		case action, ok := <-c.send:
			if !ok {
				log.Info().Msg("Closing websocket connection")
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.writeJSON(action); err != nil {
				log.Err(err).Msg("Failed to write action to websocket")
				c.write(websocket.CloseMessage, []byte{})
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Err(err).Msg("Failed to write ping")
				return
			}
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *connection) writeJSON(v interface{}) error {
	s.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return s.ws.WriteJSON(v)
}
