package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/effect707/MessngerGrusha/api/gen/user"
	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	notifuc "github.com/effect707/MessngerGrusha/internal/usecase/notification"
)

type UserGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type OnlineChecker interface {
	IsOnline(ctx context.Context, userID uuid.UUID) (bool, error)
}

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userRepo      UserGetter
	onlineChecker OnlineChecker
	notifUC       *notifuc.UseCase
}

func NewUserHandler(userRepo UserGetter, onlineChecker OnlineChecker, notifUC *notifuc.UseCase) *UserHandler {
	return &UserHandler{
		userRepo:      userRepo,
		onlineChecker: onlineChecker,
		notifUC:       notifUC,
	}
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var user *domain.User
	var err error

	userID, parseErr := uuid.Parse(req.GetUserId())
	if parseErr == nil {
		user, err = h.userRepo.GetByID(ctx, userID)
	} else {

		user, err = h.userRepo.GetByUsername(ctx, req.GetUserId())
	}

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return &pb.GetProfileResponse{
		User: &pb.User{
			Id:          user.ID.String(),
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			AvatarUrl:   avatarURL,
			Bio:         user.Bio,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *UserHandler) GetOnlineStatus(ctx context.Context, req *pb.GetOnlineStatusRequest) (*pb.GetOnlineStatusResponse, error) {
	statuses := make(map[string]bool, len(req.GetUserIds()))

	for _, idStr := range req.GetUserIds() {
		userID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		online, err := h.onlineChecker.IsOnline(ctx, userID)
		if err != nil {
			continue
		}

		statuses[idStr] = online
	}

	return &pb.GetOnlineStatusResponse{
		Statuses: statuses,
	}, nil
}

func (h *UserHandler) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var notifications []domain.Notification
	if req.GetUnreadOnly() {
		notifications, err = h.notifUC.GetUnread(ctx, userID, int(req.GetLimit()))
	} else {
		notifications, err = h.notifUC.GetNotifications(ctx, userID, int(req.GetLimit()))
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get notifications")
	}

	resp := &pb.GetNotificationsResponse{}
	for _, n := range notifications {
		resp.Notifications = append(resp.Notifications, &pb.Notification{
			Id:        n.ID.String(),
			Type:      string(n.Type),
			Payload:   string(n.Payload),
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	return resp, nil
}

func (h *UserHandler) MarkNotificationRead(ctx context.Context, req *pb.MarkNotificationReadRequest) (*pb.MarkNotificationReadResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	notifID, err := uuid.Parse(req.GetNotificationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid notification_id")
	}

	if err := h.notifUC.MarkRead(ctx, notifID, userID); err != nil {
		return nil, status.Error(codes.Internal, "failed to mark read")
	}

	return &pb.MarkNotificationReadResponse{}, nil
}

func (h *UserHandler) MarkAllNotificationsRead(ctx context.Context, _ *pb.MarkAllNotificationsReadRequest) (*pb.MarkAllNotificationsReadResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.notifUC.MarkAllRead(ctx, userID); err != nil {
		return nil, status.Error(codes.Internal, "failed to mark all read")
	}

	return &pb.MarkAllNotificationsReadResponse{}, nil
}

func (h *UserHandler) GetUnreadCount(ctx context.Context, _ *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	count, err := h.notifUC.CountUnread(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count unread")
	}

	return &pb.GetUnreadCountResponse{Count: count}, nil
}
