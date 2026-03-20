package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/effect707/MessngerGrusha/api/gen/chat"
	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	chatuc "github.com/effect707/MessngerGrusha/internal/usecase/chat"
)

type ChatHandler struct {
	pb.UnimplementedChatServiceServer
	useCase *chatuc.UseCase
}

func NewChatHandler(useCase *chatuc.UseCase) *ChatHandler {
	return &ChatHandler{useCase: useCase}
}

func (h *ChatHandler) CreateDirectChat(ctx context.Context, req *pb.CreateDirectChatRequest) (*pb.CreateDirectChatResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	recipientID, err := uuid.Parse(req.GetRecipientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid recipient_id")
	}

	chat, err := h.useCase.CreateDirectChat(ctx, userID, recipientID)
	if err != nil {
		slog.Error("CreateDirectChat failed", slog.String("error", err.Error()), slog.String("creator", userID.String()), slog.String("recipient", recipientID.String()))
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "recipient not found")
		}
		return nil, status.Error(codes.Internal, "failed to create chat")
	}

	return &pb.CreateDirectChatResponse{
		Chat: domainChatToProto(chat),
	}, nil
}

func (h *ChatHandler) CreateGroupChat(ctx context.Context, req *pb.CreateGroupChatRequest) (*pb.CreateGroupChatResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	memberIDs := make([]uuid.UUID, 0, len(req.GetMemberIds()))
	for _, idStr := range req.GetMemberIds() {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid member_id: "+idStr)
		}
		memberIDs = append(memberIDs, id)
	}

	chat, err := h.useCase.CreateGroupChat(ctx, req.GetName(), userID, memberIDs)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create group chat")
	}

	return &pb.CreateGroupChatResponse{
		Chat: domainChatToProto(chat),
	}, nil
}

func (h *ChatHandler) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	chat, err := h.useCase.GetChat(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}
		return nil, status.Error(codes.Internal, "failed to get chat")
	}

	return &pb.GetChatResponse{
		Chat: domainChatToProto(chat),
	}, nil
}

func (h *ChatHandler) GetUserChats(ctx context.Context, _ *pb.GetUserChatsRequest) (*pb.GetUserChatsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chats, err := h.useCase.GetUserChats(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get chats")
	}

	resp := &pb.GetUserChatsResponse{}
	for _, c := range chats {
		resp.Chats = append(resp.Chats, domainChatToProto(&c))
	}
	return resp, nil
}

func (h *ChatHandler) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.AddMemberResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	newMemberID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if err := h.useCase.AddMember(ctx, chatID, userID, newMemberID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "already a member")
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to add member")
	}

	return &pb.AddMemberResponse{}, nil
}

func (h *ChatHandler) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid chat_id")
	}

	targetID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if err := h.useCase.RemoveMember(ctx, chatID, userID, targetID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not a member")
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to remove member")
	}

	return &pb.RemoveMemberResponse{}, nil
}

func domainChatToProto(c *domain.Chat) *pb.Chat {
	name := ""
	if c.Name != nil {
		name = *c.Name
	}
	avatarURL := ""
	if c.AvatarURL != nil {
		avatarURL = *c.AvatarURL
	}
	return &pb.Chat{
		Id:        c.ID.String(),
		Type:      string(c.Type),
		Name:      name,
		AvatarUrl: avatarURL,
		CreatedBy: c.CreatedBy.String(),
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
