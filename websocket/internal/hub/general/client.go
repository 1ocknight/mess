package general

import (
	"fmt"
	"time"

	wsdto "github.com/1ocknight/mess/shared/dto/ws"
	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
)

type ClientConfig struct {
	MessageBuffer int           `yaml:"message_buffer"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	PingPeriod    time.Duration `yaml:"ping_timeout"`
}

type Client struct {
	SubjectID string
	send      chan *wsdto.WSMessage
	hub       *Hub
	conn      *websocket.Conn
}

func NewClient(subjectID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		SubjectID: subjectID,
		send:      make(chan *wsdto.WSMessage, hub.cfg.Client.MessageBuffer),
		hub:       hub,
		conn:      conn,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(c.hub.cfg.Client.ReadTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.hub.cfg.Client.ReadTimeout))
		return nil
	})

	for {
		mt, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.sendError(err)
			}
			break
		}

		if mt == websocket.CloseMessage {
			break
		}
	}
}

func (c *Client) sendError(err error) {
	c.hub.lg.Error(fmt.Errorf("subj: %v, err: %w", c.SubjectID, err))
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(c.hub.cfg.Client.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.cfg.Client.WriteTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.sendError(err)
				return
			}
			msg, err := message.GetBytes()
			if err != nil {
				c.sendError(err)
				return
			}

			w.Write(msg)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				msg, err := (<-c.send).GetBytes()
				if err != nil {
					c.sendError(err)
					return
				}
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				c.sendError(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.cfg.Client.WriteTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.sendError(err)
				return
			}
		}
	}
}
