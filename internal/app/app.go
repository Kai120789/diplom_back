package app

import (
	"context"
	"fmt"
	"materials/internal/config"
	"materials/internal/service"
	"materials/internal/storage"
	"materials/internal/transport/http/handler"
	"materials/internal/transport/http/router"
	"materials/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	dbConn, err := storage.Connect(cfg.DbURL)
	if err != nil {
		panic(err)
	}

	PGStore := storage.NewPosgtresStorage(dbConn, log)

	s := service.NewService(service.Storage{
		UserStorage: PGStore.UserStore,
	}, cfg, log)

	h := handler.NewHandler(handler.Service{
		UserService: s.UserService,
	}, cfg)

	r := router.NewRouter(router.Handler{
		User: h.UserHandler,
	})

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	go func() {
		log.Info("starting server", zap.String("address", cfg.RunAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("server forced to shutdown:", zap.Error(err))
	}

	log.Info("server exiting")

}
