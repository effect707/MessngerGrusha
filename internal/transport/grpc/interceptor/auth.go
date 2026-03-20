package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var publicMethods = map[string]bool{
	"/auth.AuthService/Register":      true,
	"/auth.AuthService/Login":         true,
	"/auth.AuthService/RefreshTokens": true,
}

func AuthUnaryInterceptor(tokenManager *jwtpkg.TokenManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		userID, err := extractUserID(ctx, tokenManager)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		return handler(ctx, req)
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "user id not found in context")
	}
	return id, nil
}

func extractUserID(ctx context.Context, tokenManager *jwtpkg.TokenManager) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")

	claims, err := tokenManager.ParseToken(token)
	if err != nil {
		return uuid.UUID{}, status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims.UserID, nil
}
