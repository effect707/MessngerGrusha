package ws

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
)

type Server struct {
	hub          *Hub
	tokenManager *jwtpkg.TokenManager
	handler      func(c *Client, data []byte)
	logger       *slog.Logger
}

func NewServer(hub *Hub, tokenManager *jwtpkg.TokenManager, logger *slog.Logger) *Server {
	return &Server{
		hub:          hub,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (s *Server) SetHandler(handler func(c *Client, data []byte)) {
	s.handler = handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := s.tokenManager.ParseToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		s.logger.Error("websocket accept error", slog.String("error", err.Error()))
		return
	}

	client := NewClient(claims.UserID, conn, s.hub, s.logger)
	s.hub.Register(client)

	ctx, cancel := context.WithCancel(context.Background())
	go client.WritePump(ctx)
	client.ReadPump(ctx, func(_ uuid.UUID, data []byte) {
		if s.handler != nil {
			s.handler(client, data)
		}
	})
	cancel()
}
