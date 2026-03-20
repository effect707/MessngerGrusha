package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		logger.Info("gRPC request",
			slog.String("method", info.FullMethod),
			slog.String("code", st.Code().String()),
			slog.Duration("duration", duration),
		)

		if err != nil {
			logger.Error("gRPC error",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
			)
		}

		return resp, err
	}
}
