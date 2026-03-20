package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/effect707/MessngerGrusha/api/gen/auth"
	"github.com/effect707/MessngerGrusha/internal/domain"
	"github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	authuc "github.com/effect707/MessngerGrusha/internal/usecase/auth"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	useCase *authuc.UseCase
}

func NewAuthHandler(useCase *authuc.UseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.useCase.Register(ctx, authuc.RegisterInput{
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		DisplayName: req.GetDisplayName(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &pb.RegisterResponse{
		User: domainUserToProto(user),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tokens, err := h.useCase.Login(ctx, authuc.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := h.useCase.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}
	return &pb.LogoutResponse{}, nil
}

func (h *AuthHandler) LogoutAll(ctx context.Context, _ *pb.LogoutAllRequest) (*pb.LogoutAllResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.useCase.LogoutAll(ctx, userID); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout all")
	}
	return &pb.LogoutAllResponse{}, nil
}

func (h *AuthHandler) RefreshTokens(ctx context.Context, req *pb.RefreshTokensRequest) (*pb.RefreshTokensResponse, error) {
	tokens, err := h.useCase.RefreshTokens(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		if errors.Is(err, domain.ErrTokenRevoked) {
			return nil, status.Error(codes.Unauthenticated, "token revoked")
		}
		return nil, status.Error(codes.Internal, "failed to refresh tokens")
	}

	return &pb.RefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func domainUserToProto(u *domain.User) *pb.User {
	avatarURL := ""
	if u.AvatarURL != nil {
		avatarURL = *u.AvatarURL
	}
	return &pb.User{
		Id:          u.ID.String(),
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarUrl:   avatarURL,
		Bio:         u.Bio,
		CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
