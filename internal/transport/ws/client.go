package ws

import (
	"context"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	UserID uuid.UUID
	conn   *websocket.Conn
	hub    *Hub
	send   chan []byte
	logger *slog.Logger
}

func NewClient(userID uuid.UUID, conn *websocket.Conn, hub *Hub, logger *slog.Logger) *Client {
	return &Client{
		UserID: userID,
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, 256),
		logger: logger,
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		c.logger.Warn("client send buffer full, dropping message",
			slog.String("user_id", c.UserID.String()))
	}
}

func (c *Client) ReadPump(ctx context.Context, handler func(userID uuid.UUID, data []byte)) {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.CloseNow()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) != -1 {
				c.logger.Info("websocket closed",
					slog.String("user_id", c.UserID.String()))
			} else {
				c.logger.Error("websocket read error",
					slog.String("user_id", c.UserID.String()),
					slog.String("error", err.Error()))
			}
			return
		}

		handler(c.UserID, data)
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.CloseNow()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.Close(websocket.StatusNormalClosure, "")
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()

			if err != nil {
				c.logger.Error("websocket write error",
					slog.String("user_id", c.UserID.String()),
					slog.String("error", err.Error()))
				return
			}

		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(pingCtx)
			cancel()

			if err != nil {
				return
			}

			c.hub.RefreshOnline(c.UserID)

		case <-ctx.Done():
			return
		}
	}
}
