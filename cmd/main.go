package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"msn/internal/config"
	"msn/internal/infra/database/pg"
	"msn/internal/infra/http/middleware"
	"msn/internal/infra/http/server"
	"msn/internal/infra/logging"
	"msn/internal/modules/auth"
	"msn/internal/modules/categories"
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

	router.Use(middleware.Logging)

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

	userRepo := user.NewRepo(pgConn.DB())
	sessionRepo := session.NewRepo(pgConn.DB())
	categoriesRepo := categories.NewRepo(pgConn.DB())

	userService := user.NewUserService(user.ServiceConfig{
		UserRepo:     userRepo,
		CategoryRepo: categoriesRepo,
	})
	sessionService := session.NewSessionService(session.ServiceConfig{
		SessionRepo: sessionRepo,
		UserService: userService,
	})
	authService := auth.NewAuthService(auth.ServiceConfig{
		UserRepo:       userRepo,
		SessionRepo:    sessionRepo,
		SessionService: sessionService,
		AccessKey:      cfg.JWTAccessKey,
		RefreshKey:     cfg.JWTRefreshKey,
	})
	categoriesService := categories.NewService(categories.ServiceConfig{
		CategoriesRepo: categoriesRepo,
	})

	user.NewHandler(userService).RegisterRoutes(router)
	auth.NewHandler(authService, cfg.JWTAccessKey).RegisterRoutes(router)
	categories.NewHandler(categoriesService).RegisterRoutes(router)

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
