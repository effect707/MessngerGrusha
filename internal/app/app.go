package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	pb_auth "github.com/effect707/MessngerGrusha/api/gen/auth"
	pb_channel "github.com/effect707/MessngerGrusha/api/gen/channel"
	pb_chat "github.com/effect707/MessngerGrusha/api/gen/chat"
	pb_msg "github.com/effect707/MessngerGrusha/api/gen/message"
	pb_user "github.com/effect707/MessngerGrusha/api/gen/user"
	"github.com/effect707/MessngerGrusha/internal/config"
	"github.com/effect707/MessngerGrusha/internal/pkg/hasher"
	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
	miniorepo "github.com/effect707/MessngerGrusha/internal/repository/minio"
	"github.com/effect707/MessngerGrusha/internal/repository/postgres"
	redisrepo "github.com/effect707/MessngerGrusha/internal/repository/redis"
	grpctransport "github.com/effect707/MessngerGrusha/internal/transport/grpc"
	grpcinterceptor "github.com/effect707/MessngerGrusha/internal/transport/grpc/interceptor"
	httptransport "github.com/effect707/MessngerGrusha/internal/transport/http"
	"github.com/effect707/MessngerGrusha/internal/transport/ws"
	authuc "github.com/effect707/MessngerGrusha/internal/usecase/auth"
	channeluc "github.com/effect707/MessngerGrusha/internal/usecase/channel"
	chatuc "github.com/effect707/MessngerGrusha/internal/usecase/chat"
	fileuc "github.com/effect707/MessngerGrusha/internal/usecase/file"
	msguc "github.com/effect707/MessngerGrusha/internal/usecase/message"
	notifuc "github.com/effect707/MessngerGrusha/internal/usecase/notification"
	reactionuc "github.com/effect707/MessngerGrusha/internal/usecase/reaction"
)

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	grpcServer *grpc.Server
	httpServer *http.Server
	pool       *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	logger.Info("connected to PostgreSQL")

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err = redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	logger.Info("connected to Redis")

	fileStorage, err := miniorepo.NewFileStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.Bucket,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	if err := fileStorage.EnsureBucket(ctx); err != nil {
		return nil, fmt.Errorf("ensure minio bucket: %w", err)
	}
	logger.Info("connected to MinIO")

	userRepo := postgres.NewUserRepository(pool)
	tokenRepo := postgres.NewTokenRepository(pool)
	sessionRepo := redisrepo.NewSessionRepository(redisClient)
	msgRepo := postgres.NewMessageRepository(pool)
	chatRepo := postgres.NewChatRepository(pool)
	typingRepo := redisrepo.NewTypingRepository(redisClient)
	onlineRepo := redisrepo.NewOnlineStatusRepository(redisClient)
	attachmentRepo := postgres.NewAttachmentRepository(pool)
	reactionRepo := postgres.NewReactionRepository(pool)
	channelRepo := postgres.NewChannelRepository(pool)
	notifRepo := postgres.NewNotificationRepository(pool)

	tokenManager := jwtpkg.NewTokenManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	passwordHasher := hasher.NewBcryptHasher()

	authUseCase := authuc.NewUseCase(userRepo, tokenRepo, sessionRepo, passwordHasher, tokenManager)
	chatUseCase := chatuc.NewUseCase(chatRepo, userRepo)
	fileUseCase := fileuc.NewUseCase(fileStorage, attachmentRepo, msgRepo, chatRepo)
	reactionUseCase := reactionuc.NewUseCase(reactionRepo, msgRepo, chatRepo)
	channelUseCase := channeluc.NewUseCase(channelRepo)
	notifUseCase := notifuc.NewUseCase(notifRepo)

	hub := ws.NewHub(logger)
	hub.SetOnlineTracker(onlineRepo)
	broker := ws.NewLocalBroker(hub)
	msgUseCase := msguc.NewUseCase(msgRepo, chatRepo, broker)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcinterceptor.LoggingUnaryInterceptor(logger),
			grpcinterceptor.AuthUnaryInterceptor(tokenManager),
		),
	)

	authHandler := grpctransport.NewAuthHandler(authUseCase)
	pb_auth.RegisterAuthServiceServer(grpcServer, authHandler)

	msgHandler := grpctransport.NewMessageHandler(msgUseCase, reactionUseCase, fileUseCase)
	pb_msg.RegisterMessageServiceServer(grpcServer, msgHandler)

	chatHandler := grpctransport.NewChatHandler(chatUseCase)
	pb_chat.RegisterChatServiceServer(grpcServer, chatHandler)

	userHandler := grpctransport.NewUserHandler(userRepo, onlineRepo, notifUseCase)
	pb_user.RegisterUserServiceServer(grpcServer, userHandler)

	channelHandler := grpctransport.NewChannelHandler(channelUseCase)
	pb_channel.RegisterChannelServiceServer(grpcServer, channelHandler)

	wsServer := ws.NewServer(hub, tokenManager, logger)
	wsHandler := ws.NewWSHandler(hub, msgUseCase, typingRepo, chatRepo, logger)
	wsServer.SetHandler(wsHandler.HandleMessage)

	fileHandler := httptransport.NewFileHandler(fileUseCase, tokenManager, logger)

	mux := http.NewServeMux()
	mux.Handle("/ws", wsServer)
	mux.HandleFunc("/api/files/upload", fileHandler.Upload)
	mux.HandleFunc("/api/files/download", fileHandler.Download)

	grpcAddr := fmt.Sprintf("localhost:%d", cfg.GRPC.Port)
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
	)
	gwOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gwCtx := context.Background()
	if err := pb_auth.RegisterAuthServiceHandlerFromEndpoint(gwCtx, gwMux, grpcAddr, gwOpts); err != nil {
		return nil, fmt.Errorf("register auth gateway: %w", err)
	}
	if err := pb_chat.RegisterChatServiceHandlerFromEndpoint(gwCtx, gwMux, grpcAddr, gwOpts); err != nil {
		return nil, fmt.Errorf("register chat gateway: %w", err)
	}
	if err := pb_msg.RegisterMessageServiceHandlerFromEndpoint(gwCtx, gwMux, grpcAddr, gwOpts); err != nil {
		return nil, fmt.Errorf("register message gateway: %w", err)
	}
	if err := pb_channel.RegisterChannelServiceHandlerFromEndpoint(gwCtx, gwMux, grpcAddr, gwOpts); err != nil {
		return nil, fmt.Errorf("register channel gateway: %w", err)
	}
	if err := pb_user.RegisterUserServiceHandlerFromEndpoint(gwCtx, gwMux, grpcAddr, gwOpts); err != nil {
		return nil, fmt.Errorf("register user gateway: %w", err)
	}

	mux.Handle("/api/", gwMux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "ok")
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: mux,
	}

	return &App{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
		httpServer: httpServer,
		pool:       pool,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	errCh := make(chan error, 2)

	go func() {
		a.logger.Info("gRPC server starting", slog.Int("port", a.cfg.GRPC.Port))
		if err := a.grpcServer.Serve(grpcLis); err != nil {
			errCh <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	go func() {
		a.logger.Info("HTTP/WS server starting", slog.Int("port", a.cfg.HTTP.Port))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("shutting down gracefully...")
	case err := <-errCh:
		a.logger.Error("server error", slog.String("error", err.Error()))
		cancel()
	}

	a.grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}

	a.pool.Close()

	a.logger.Info("servers stopped")
	return nil
}
