package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/effect707/MessngerGrusha/api/gen/channel"
	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	channeluc "github.com/effect707/MessngerGrusha/internal/usecase/channel"
)

type ChannelHandler struct {
	pb.UnimplementedChannelServiceServer
	useCase *channeluc.UseCase
}

func NewChannelHandler(useCase *channeluc.UseCase) *ChannelHandler {
	return &ChannelHandler{useCase: useCase}
}

func (h *ChannelHandler) CreateChannel(ctx context.Context, req *pb.CreateChannelRequest) (*pb.CreateChannelResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := h.useCase.Create(ctx, channeluc.CreateInput{
		Slug:        req.GetSlug(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		IsPrivate:   req.GetIsPrivate(),
		OwnerID:     userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create channel")
	}

	return &pb.CreateChannelResponse{
		Channel: domainChannelToProto(ch),
	}, nil
}

func (h *ChannelHandler) GetChannel(ctx context.Context, req *pb.GetChannelRequest) (*pb.GetChannelResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelID, err := uuid.Parse(req.GetChannelId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel_id")
	}

	ch, err := h.useCase.GetByID(ctx, channelID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "private channel")
		}
		return nil, status.Error(codes.Internal, "failed to get channel")
	}

	return &pb.GetChannelResponse{
		Channel: domainChannelToProto(ch),
	}, nil
}

func (h *ChannelHandler) UpdateChannel(ctx context.Context, req *pb.UpdateChannelRequest) (*pb.UpdateChannelResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelID, err := uuid.Parse(req.GetChannelId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel_id")
	}

	ch, err := h.useCase.Update(ctx, channelID, userID, channeluc.UpdateInput{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		IsPrivate:   req.GetIsPrivate(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not the owner")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}
		return nil, status.Error(codes.Internal, "failed to update channel")
	}

	return &pb.UpdateChannelResponse{
		Channel: domainChannelToProto(ch),
	}, nil
}

func (h *ChannelHandler) DeleteChannel(ctx context.Context, req *pb.DeleteChannelRequest) (*pb.DeleteChannelResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelID, err := uuid.Parse(req.GetChannelId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel_id")
	}

	if err := h.useCase.Delete(ctx, channelID, userID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "not the owner")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete channel")
	}

	return &pb.DeleteChannelResponse{}, nil
}

func (h *ChannelHandler) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelID, err := uuid.Parse(req.GetChannelId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel_id")
	}

	if err := h.useCase.Subscribe(ctx, channelID, userID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "private channel")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}
		return nil, status.Error(codes.Internal, "failed to subscribe")
	}

	return &pb.SubscribeResponse{}, nil
}

func (h *ChannelHandler) Unsubscribe(ctx context.Context, req *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelID, err := uuid.Parse(req.GetChannelId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel_id")
	}

	if err := h.useCase.Unsubscribe(ctx, channelID, userID); err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to unsubscribe")
	}

	return &pb.UnsubscribeResponse{}, nil
}

func (h *ChannelHandler) GetPublicChannels(ctx context.Context, req *pb.GetPublicChannelsRequest) (*pb.GetPublicChannelsResponse, error) {
	channels, err := h.useCase.GetPublicChannels(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get channels")
	}

	resp := &pb.GetPublicChannelsResponse{}
	for _, ch := range channels {
		resp.Channels = append(resp.Channels, domainChannelToProto(&ch))
	}
	return resp, nil
}

func (h *ChannelHandler) GetMyChannels(ctx context.Context, _ *pb.GetMyChannelsRequest) (*pb.GetMyChannelsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channels, err := h.useCase.GetUserChannels(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get channels")
	}

	resp := &pb.GetMyChannelsResponse{}
	for _, ch := range channels {
		resp.Channels = append(resp.Channels, domainChannelToProto(&ch))
	}
	return resp, nil
}

func domainChannelToProto(ch *domain.Channel) *pb.Channel {
	avatarURL := ""
	if ch.AvatarURL != nil {
		avatarURL = *ch.AvatarURL
	}
	return &pb.Channel{
		Id:          ch.ID.String(),
		Slug:        ch.Slug,
		Name:        ch.Name,
		Description: ch.Description,
		AvatarUrl:   avatarURL,
		OwnerId:     ch.OwnerID.String(),
		IsPrivate:   ch.IsPrivate,
		CreatedAt:   ch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
