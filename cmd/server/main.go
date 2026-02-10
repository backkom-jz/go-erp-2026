package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-erp/internal/bootstrap"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	zapLogger, err := bootstrap.InitLogger(cfg.Log)
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer func() { _ = zapLogger.Sync() }()

	gormLogLevel := logger.Warn
	if cfg.Server.Mode == "debug" {
		gormLogLevel = logger.Info
	}

	db, err := bootstrap.InitDB(cfg.DB, gormLogLevel)
	if err != nil {
		zapLogger.Fatal("init db failed", zap.Error(err))
	}
	if sqlDB, err := db.DB(); err == nil {
		defer func() { _ = sqlDB.Close() }()
	}

	redisClient, err := bootstrap.InitRedis(cfg.Redis)
	if err != nil {
		zapLogger.Fatal("init redis failed", zap.Error(err))
	}
	if redisClient != nil {
		defer func() { _ = redisClient.Close() }()
	}

	r := bootstrap.InitRouter(cfg.Server, zapLogger)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}

	go func() {
		zapLogger.Info("server starting", zap.String("addr", cfg.Server.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("server shutdown failed", zap.Error(err))
	}
}
