package ws

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
	msguc "github.com/effect707/MessngerGrusha/internal/usecase/message"
)

type TypingNotifier interface {
	SetTyping(ctx context.Context, chatID, userID uuid.UUID) error
}

type MemberLister interface {
	GetMembers(ctx context.Context, chatID uuid.UUID) ([]domain.ChatMember, error)
}

type WSHandler struct {
	hub            *Hub
	msgUseCase     *msguc.UseCase
	typingNotifier TypingNotifier
	memberLister   MemberLister
	logger         *slog.Logger
}

func NewWSHandler(hub *Hub, msgUseCase *msguc.UseCase, typingNotifier TypingNotifier, memberLister MemberLister, logger *slog.Logger) *WSHandler {
	return &WSHandler{
		hub:            hub,
		msgUseCase:     msgUseCase,
		typingNotifier: typingNotifier,
		memberLister:   memberLister,
		logger:         logger,
	}
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
	MsgType string `json:"msg_type"`
}

type TypingPayload struct {
	ChatID string `json:"chat_id"`
}

type OutgoingMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type NewMessageEvent struct {
	ID        string `json:"id"`
	ChatID    string `json:"chat_id"`
	SenderID  string `json:"sender_id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type TypingEvent struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

func (h *WSHandler) HandleMessage(client *Client, data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		h.logger.Error("invalid ws message", slog.String("error", err.Error()))
		return
	}

	ctx := context.Background()

	switch msg.Type {
	case "send_message":
		h.handleSendMessage(ctx, client, msg.Payload)
	case "typing":
		h.handleTyping(ctx, client, msg.Payload)
	default:
		h.logger.Warn("unknown ws message type", slog.String("type", msg.Type))
	}
}

func (h *WSHandler) handleSendMessage(ctx context.Context, client *Client, payload json.RawMessage) {
	var p SendMessagePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		h.logger.Error("invalid send_message payload", slog.String("error", err.Error()))
		return
	}

	chatID, err := uuid.Parse(p.ChatID)
	if err != nil {
		h.logger.Error("invalid chat_id in ws message")
		return
	}

	msgType := domain.MessageTypeText
	if p.MsgType != "" {
		msgType = domain.MessageType(p.MsgType)
	}

	msg, err := h.msgUseCase.Send(ctx, msguc.SendInput{
		ChatID:   chatID,
		SenderID: client.UserID,
		Type:     msgType,
		Content:  p.Content,
	})
	if err != nil {
		h.logger.Error("failed to send message via ws",
			slog.String("error", err.Error()),
			slog.String("user_id", client.UserID.String()))
		return
	}

	event := OutgoingMessage{
		Type: "new_message",
		Payload: NewMessageEvent{
			ID:        msg.ID.String(),
			ChatID:    chatID.String(),
			SenderID:  client.UserID.String(),
			Type:      string(msg.Type),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	members, err := h.memberLister.GetMembers(ctx, chatID)
	if err != nil {
		h.logger.Error("failed to get chat members", slog.String("error", err.Error()))
		return
	}

	memberIDs := make([]uuid.UUID, 0, len(members))
	for _, m := range members {
		memberIDs = append(memberIDs, m.UserID)
	}

	h.hub.SendToUsers(memberIDs, event)
}

func (h *WSHandler) handleTyping(ctx context.Context, client *Client, payload json.RawMessage) {
	var p TypingPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return
	}

	chatID, err := uuid.Parse(p.ChatID)
	if err != nil {
		return
	}

	if err := h.typingNotifier.SetTyping(ctx, chatID, client.UserID); err != nil {
		h.logger.Error("failed to set typing", slog.String("error", err.Error()))
		return
	}

	members, err := h.memberLister.GetMembers(ctx, chatID)
	if err != nil {
		return
	}

	event := OutgoingMessage{
		Type: "typing",
		Payload: TypingEvent{
			ChatID: chatID.String(),
			UserID: client.UserID.String(),
		},
	}

	for _, m := range members {
		if m.UserID != client.UserID {
			h.hub.SendToUser(m.UserID, event)
		}
	}
}
