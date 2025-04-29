package main

import (
	"context"
	"errors"
	"log/slog"
	"msn/services/user-service/internal/config"
	"msn/services/user-service/internal/infra/database/pg"
	"msn/services/user-service/internal/infra/http/server"
	"msn/services/user-service/internal/modules/user"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	cfg := config.GetConfig()
	ctx := context.Background()
	r := chi.NewRouter()

	pgconn, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer pgconn.Close()

	userRepo := user.NewRepo(pgconn.DB())

	userService := user.NewService(user.ServiceConfig{
		UserRepo: userRepo,
	})

	user.NewHandler(userService).RegisterRoutes(r)

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
