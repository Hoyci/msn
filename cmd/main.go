package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"msn/internal/config"
	"msn/internal/infra/database/pg"
	categoryRepository "msn/internal/infra/database/pg/repositories/category"
	roleRepository "msn/internal/infra/database/pg/repositories/role"
	sessionRepository "msn/internal/infra/database/pg/repositories/session"
	userRepository "msn/internal/infra/database/pg/repositories/user"
	authHandler "msn/internal/infra/http/handlers/auth"
	categoryhandler "msn/internal/infra/http/handlers/category"
	userHandler "msn/internal/infra/http/handlers/user"
	"msn/internal/infra/http/middlewares"
	"msn/internal/infra/http/server"
	"msn/internal/infra/jwt"
	"msn/internal/infra/logging"
	"msn/internal/infra/storage"
	"msn/internal/modules/auth"
	"msn/internal/modules/category"
	"msn/internal/modules/session"
	"msn/internal/modules/user"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.NewLogger(os.Stdout, cfg.Environment)
	slog.SetDefault(logger)

	slog.Info(fmt.Sprintf("Launching %s with the following settings:", cfg.AppName),
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	ctx := context.Background()
	router := chi.NewRouter()

	router.Use(middlewares.Logging)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	pgConn, err := pg.NewPostgresConnection(cfg.PostgresDSN)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer pgConn.Close()

	storageClient := storage.NewStorageClient(cfg.StorageURL, cfg.StorageAccessKey, cfg.StorageSecretKey)

	userRepo := userRepository.NewRepo(pgConn.DB())
	categoryRepo := categoryRepository.NewRepo(pgConn.DB())
	sessionRepo := sessionRepository.NewRepo(pgConn.DB())
	roleRepo := roleRepository.NewRepo(pgConn.DB())

	tokenProvider := jwt.NewProvider(cfg.JWTAccessKey, cfg.JWTRefreshKey)

	userService := user.NewService(user.ServiceConfig{
		UserRepo:      userRepo,
		CategoryRepo:  categoryRepo,
		RoleRepo:      roleRepo,
		StorageClient: storageClient,
	})
	sessionService := session.NewService(session.ServiceConfig{
		SessionRepo: sessionRepo,
		UserService: userService,
	})
	authService := auth.NewService(auth.ServiceConfig{
		UserRepo:       userRepo,
		SessionService: sessionService,
		TokenProvider:  *tokenProvider,
	})
	categoryService := category.NewService(category.ServiceConfig{
		CategoryRepo: categoryRepo,
	})

	authHandler.NewHandler(authService, cfg.JWTAccessKey).RegisterRoutes(router)
	userHandler.NewHandler(userService, storageClient).RegisterRoutes(router)
	categoryhandler.NewHandler(categoryService).RegisterRoutes(router)

	srv := server.New(server.Config{
		Port:         cfg.Port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Router:       router,
	})

	shutdownErr := srv.GracefulShutdown(ctx, time.Second*30)

	err = srv.Start()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	err = <-shutdownErr
	if err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}

	slog.Info("server shoutdown gracefully")
}
