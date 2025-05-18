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
	r := chi.NewRouter()

	r.Use(middleware.Logging)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	pgconn, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer pgconn.Close()

	userRepo := user.NewRepo(pgconn.DB())
	sessionRepo := session.NewRepo(pgconn.DB())
	categoriesRepo := categories.NewRepo(pgconn.DB())

	userService := user.NewService(user.ServiceConfig{
		UserRepo:     userRepo,
		CategoryRepo: categoriesRepo,
	})
	sessionService := session.NewService(session.ServiceConfig{
		SessionRepo: sessionRepo,
		UserService: userService,
	})
	authService := auth.NewService(auth.ServiceConfig{
		UserRepo:       userRepo,
		SessionRepo:    sessionRepo,
		SessionService: sessionService,
		AccessKey:      cfg.JWTAccessKey,
		RefreshKey:     cfg.JWTRefreshKey,
	})
	categoriesService := categories.NewService(categories.ServiceConfig{
		CategoriesRepo: categoriesRepo,
	})

	user.NewHandler(userService).RegisterRoutes(r)
	auth.NewHandler(authService, cfg.JWTAccessKey).RegisterRoutes(r)
	categories.NewHandler(categoriesService).RegisterRoutes(r)

	srv := server.New(server.Config{
		Port:         cfg.Port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Router:       r,
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
