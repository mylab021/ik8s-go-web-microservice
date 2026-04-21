package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mylab021/go-web-microservice/common/database"
	"github.com/mylab021/go-web-microservice/common/logger"
	"github.com/mylab021/go-web-microservice/user-service/internal/config"
	"github.com/mylab021/go-web-microservice/user-service/internal/infrastructure/http"
	"github.com/mylab021/go-web-microservice/user-service/internal/infrastructure/persistence/gorm"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. 初始化日志
	logger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatal("Failed to init logger:", err)
	}
	defer logger.Sync()

	// 3. 初始化数据库
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer db.Close()

	// 4. 运行数据库迁移
	if err := gorm.RunMigrations(db, cfg.Database); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// 5. 初始化仓储
	userRepo := gorm.NewUserRepository(db, logger)

	// 6. 初始化服务
	userService := service.NewUserService(userRepo, logger)

	// 7. 创建 HTTP 服务器
	server := http.NewServer(userService, logger, cfg)

	// 8. 启动服务器
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 9. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
