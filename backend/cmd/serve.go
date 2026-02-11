package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/middleware"
	"community-elderly-care-platform/internal/router"
	"community-elderly-care-platform/pkg/database"
	"community-elderly-care-platform/pkg/logger"
	"community-elderly-care-platform/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	// serve 命令标志
	configPath string
	port       int
	mode       string
	dbHost     string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 sCare 后端服务",
	Long: `启动 sCare 社区养老服务平台后端服务

示例:
  scare serve                    # 使用默认配置启动
  scare serve --port 8081        # 指定端口
  scare serve --mode release     # 生产模式
  scare serve --config /etc/scare/.env  # 指定配置文件`,
	Run: runServe,
}

func init() {
	// serve 命令标志
	serveCmd.Flags().StringVarP(&configPath, "config", "c", ".env", "配置文件路径")
	serveCmd.Flags().IntVarP(&port, "port", "p", 0, "服务端口（0表示使用配置文件）")
	serveCmd.Flags().StringVar(&mode, "mode", "", "运行模式: debug/release/test")
	serveCmd.Flags().StringVar(&dbHost, "db-host", "", "数据库主机")
}

func runServe(cmd *cobra.Command, args []string) {
	// 如果命令行指定了不同的配置文件，通过环境变量传递给配置加载器
	if configPath != ".env" {
		log.Printf("使用自定义配置文件: %s", configPath)
		os.Setenv("CONFIG_PATH", configPath)
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 显示加载的配置信息
	log.Printf("✅ 配置加载成功")
	log.Printf("   数据库: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("   Redis: %s:%d", cfg.Redis.Host, cfg.Redis.Port)

	// 覆盖命令行标志
	if port > 0 {
		log.Printf("   端口覆盖: %d -> %d", cfg.Server.Port, port)
		cfg.Server.Port = port
	}
	if mode != "" {
		log.Printf("   模式覆盖: %s -> %s", cfg.Server.Mode, mode)
		cfg.Server.Mode = mode
	}
	if dbHost != "" {
		log.Printf("   数据库主机覆盖: %s -> %s", cfg.Database.Host, dbHost)
		cfg.Database.Host = dbHost
	}

	// 初始化日志
	if err := logger.Init(logger.Config{
		Level:        cfg.Log.Level,
		Format:       cfg.Log.Format,
		OutputPath:   cfg.Log.OutputPath,
		ErrorPath:    cfg.Log.ErrorPath,
		MaxSize:      cfg.Log.MaxSize,
		MaxBackups:   cfg.Log.MaxBackups,
		MaxAge:       cfg.Log.MaxAge,
		Compress:     cfg.Log.Compress,
		EnableFile:   cfg.Log.EnableFile,
		EnableStdout: cfg.Log.EnableStdout,
	}); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// 初始化数据库
	logger.Info("正在连接数据库: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	db, err := database.InitMySQL(database.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		Charset:      cfg.Database.Charset,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
	})
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	logger.Info("数据库连接成功")

	// 初始化 Redis
	logger.Info("正在连接 Redis: %s:%d", cfg.Redis.Host, cfg.Redis.Port)
	rdb, err := redis.InitRedis(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	logger.Info("Redis 连接成功")

	// 设置 Gin 模式
	switch strings.ToLower(cfg.Server.Mode) {
	case "release", "production":
		gin.SetMode(gin.ReleaseMode)
		logger.Info("运行模式: Release")
	case "test":
		gin.SetMode(gin.TestMode)
		logger.Info("运行模式: Test")
	default:
		gin.SetMode(gin.DebugMode)
		logger.Info("运行模式: Debug")
	}

	// 创建路由引擎
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 注册中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// 注册路由
	if err := router.Register(r, db, rdb, cfg); err != nil {
		log.Fatalf("初始化路由失败: %v", err)
	}

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		separator := "=================================================="
		logger.Info(separator)
		logger.Info("sCare 服务启动成功")
		logger.Info("服务地址: http://localhost:%d", cfg.Server.Port)
		logger.Info("API 文档: http://localhost:%d/swagger/index.html", cfg.Server.Port)
		logger.Info(separator)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败: %v", err)
	}

	logger.Info("服务器已关闭")
}
