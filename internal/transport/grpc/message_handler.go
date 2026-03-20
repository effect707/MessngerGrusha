package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/effect707/MessngerGrusha/api/gen/message"
	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/pkg/pagination"
	"github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	fileuc "github.com/effect707/MessngerGrusha/internal/usecase/file"
	msguc "github.com/effect707/MessngerGrusha/internal/usecase/message"
	reactionuc "github.com/effect707/MessngerGrusha/internal/usecase/reaction"
)

type MessageHandler struct {
	pb.UnimplementedMessageServiceServer
	useCase    *msguc.UseCase
	reactionUC *reactionuc.UseCase
	fileUC     *fileuc.UseCase
}

func NewMessageHandler(useCase *msguc.UseCase, reactionUC *reactionuc.UseCase, fileUC *fileuc.UseCase) *MessageHandler {
	return &MessageHandler{
		useCase:    useCase,
		reactionUC: reactionUC,
		fileUC:     fileUC,
	}
}

func (h *MessageHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	input := msguc.SendInput{
		ChatID:   chatID,
		SenderID: userID,
		Type:     domain.MessageType(req.GetType()),
		Content:  req.GetContent(),
	}

	if req.ReplyToId != nil {
		replyID, err := uuid.Parse(*req.ReplyToId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid reply_to_id")
		}
		input.ReplyToID = &replyID
	}

	msg, err := h.useCase.Send(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member of this chat")
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to send message")
	}

	return &pb.SendMessageResponse{
		Message: domainMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) GetHistory(ctx context.Context, req *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	var cursor *pagination.Cursor
	if req.CursorId != nil && req.CursorCreatedAt != nil {
		t, err := time.Parse(time.RFC3339Nano, *req.CursorCreatedAt)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid cursor_created_at")
		}
		cursor = &pagination.Cursor{
			ID:        *req.CursorId,
			CreatedAt: t,
		}
	}

	page, err := h.useCase.GetHistory(ctx, chatID, userID, cursor, int(req.GetLimit()))
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member of this chat")
		}
		return nil, status.Error(codes.Internal, "failed to get history")
	}

	resp := &pb.GetHistoryResponse{
		HasMore: page.HasMore,
	}

	for _, msg := range page.Items {
		resp.Messages = append(resp.Messages, domainMessageToProto(&msg))
	}

	if page.NextCursor != nil {
		resp.NextCursorId = &page.NextCursor.ID
		t := page.NextCursor.CreatedAt.Format(time.RFC3339Nano)
		resp.NextCursorCreatedAt = &t
	}

	return resp, nil
}

func (h *MessageHandler) SearchMessages(ctx context.Context, req *pb.SearchMessagesRequest) (*pb.SearchMessagesResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	messages, err := h.useCase.SearchMessages(ctx, chatID, userID, req.GetQuery(), int(req.GetLimit()))
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member of this chat")
		}
		return nil, status.Error(codes.Internal, "failed to search messages")
	}

	resp := &pb.SearchMessagesResponse{}
	for _, msg := range messages {
		resp.Messages = append(resp.Messages, domainMessageToProto(&msg))
	}

	return resp, nil
}

func (h *MessageHandler) AddReaction(ctx context.Context, req *pb.AddReactionRequest) (*pb.AddReactionResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	messageID, err := uuid.Parse(req.GetMessageId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid message_id")
	}

	if err := h.reactionUC.AddReaction(ctx, messageID, userID, req.GetEmoji()); err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "message not found")
		}
		return nil, status.Error(codes.Internal, "failed to add reaction")
	}

	return &pb.AddReactionResponse{}, nil
}

func (h *MessageHandler) RemoveReaction(ctx context.Context, req *pb.RemoveReactionRequest) (*pb.RemoveReactionResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	messageID, err := uuid.Parse(req.GetMessageId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid message_id")
	}

	if err := h.reactionUC.RemoveReaction(ctx, messageID, userID, req.GetEmoji()); err != nil {
		return nil, status.Error(codes.Internal, "failed to remove reaction")
	}

	return &pb.RemoveReactionResponse{}, nil
}

func (h *MessageHandler) GetReactions(ctx context.Context, req *pb.GetReactionsRequest) (*pb.GetReactionsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	messageID, err := uuid.Parse(req.GetMessageId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid message_id")
	}

	reactions, err := h.reactionUC.GetReactions(ctx, messageID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		return nil, status.Error(codes.Internal, "failed to get reactions")
	}

	resp := &pb.GetReactionsResponse{}
	for _, r := range reactions {
		resp.Reactions = append(resp.Reactions, &pb.Reaction{
			MessageId: r.MessageID.String(),
			UserId:    r.UserID.String(),
			Emoji:     r.Emoji,
			CreatedAt: r.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	return resp, nil
}

func (h *MessageHandler) GetAttachments(ctx context.Context, req *pb.GetAttachmentsRequest) (*pb.GetAttachmentsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	messageID, err := uuid.Parse(req.GetMessageId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid message_id")
	}

	attachments, err := h.fileUC.GetAttachments(ctx, messageID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		return nil, status.Error(codes.Internal, "failed to get attachments")
	}

	resp := &pb.GetAttachmentsResponse{}
	for _, a := range attachments {
		att := &pb.Attachment{
			Id:        a.ID.String(),
			MessageId: a.MessageID.String(),
			FileName:  a.FileName,
			FileSize:  a.FileSize,
			MimeType:  a.MimeType,
			CreatedAt: a.CreatedAt.Format(time.RFC3339Nano),
		}
		if a.DurationMs != nil {
			d := int32(*a.DurationMs)
			att.DurationMs = &d
		}
		resp.Attachments = append(resp.Attachments, att)
	}

	return resp, nil
}

func domainMessageToProto(m *domain.Message) *pb.Message {
	msg := &pb.Message{
		Id:        m.ID.String(),
		SenderId:  m.SenderID.String(),
		Type:      string(m.Type),
		Content:   m.Content,
		IsEdited:  m.IsEdited,
		CreatedAt: m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339Nano),
	}
	if m.ChatID != nil {
		msg.ChatId = m.ChatID.String()
	}
	if m.ReplyToID != nil {
		s := m.ReplyToID.String()
		msg.ReplyToId = &s
	}
	return msg
}
