package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type MessageBroker interface {
	Publish(ctx context.Context, chatID string, msg []byte) error
	Subscribe(ctx context.Context, chatID string) (<-chan []byte, error)
	Unsubscribe(ctx context.Context, chatID string) error
}

type OnlineTracker interface {
	SetOnline(ctx context.Context, userID uuid.UUID) error
	SetOffline(ctx context.Context, userID uuid.UUID) error
}

type Hub struct {
	mu            sync.RWMutex
	clients       map[uuid.UUID]*Client
	onlineTracker OnlineTracker
	logger        *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]*Client),
		logger:  logger,
	}
}

func (h *Hub) SetOnlineTracker(tracker OnlineTracker) {
	h.onlineTracker = tracker
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.UserID] = client
	h.logger.Info("client connected", slog.String("user_id", client.UserID.String()))

	if h.onlineTracker != nil {
		if err := h.onlineTracker.SetOnline(context.Background(), client.UserID); err != nil {
			h.logger.Error("failed to set online status", slog.String("error", err.Error()))
		}
	}
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client.UserID)
	h.logger.Info("client disconnected", slog.String("user_id", client.UserID.String()))

	if h.onlineTracker != nil {
		if err := h.onlineTracker.SetOffline(context.Background(), client.UserID); err != nil {
			h.logger.Error("failed to set offline status", slog.String("error", err.Error()))
		}
	}
}

func (h *Hub) RefreshOnline(userID uuid.UUID) {
	if h.onlineTracker != nil {
		if err := h.onlineTracker.SetOnline(context.Background(), userID); err != nil {
			h.logger.Error("failed to refresh online status", slog.String("error", err.Error()))
		}
	}
}

func (h *Hub) GetClient(userID uuid.UUID) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[userID]
	return c, ok
}

func (h *Hub) SendToUser(userID uuid.UUID, msg any) error {
	client, ok := h.GetClient(userID)
	if !ok {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	client.Send(data)
	return nil
}

func (h *Hub) SendToUsers(userIDs []uuid.UUID, msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal message", slog.String("error", err.Error()))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, id := range userIDs {
		if client, ok := h.clients[id]; ok {
			client.Send(data)
		}
	}
}

type LocalBroker struct {
	hub *Hub
}

func NewLocalBroker(hub *Hub) *LocalBroker {
	return &LocalBroker{hub: hub}
}

func (b *LocalBroker) Publish(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (b *LocalBroker) Subscribe(_ context.Context, _ string) (<-chan []byte, error) {
	return nil, nil
}

func (b *LocalBroker) Unsubscribe(_ context.Context, _ string) error {
	return nil
}
