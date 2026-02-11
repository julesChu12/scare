package router

import (
	"net/http"

	_ "community-elderly-care-platform/docs" // swagger docs
	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/handler"
	"community-elderly-care-platform/internal/middleware"
	"community-elderly-care-platform/pkg/database"
	"community-elderly-care-platform/pkg/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Register 注册所有路由
func Register(engine *gin.Engine, db *database.DB, rdb *redis.Client, cfg *config.Config) error {
	// 初始化依赖
	deps, err := NewDeps(db, rdb, cfg)
	if err != nil {
		return err
	}

	// 静态文件服务
	engine.Static("/static", "./storage")

	// 基础路由
	api := engine.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		handler.Respond(c, http.StatusOK, "ok", gin.H{"service": cfg.Server.Name})
	})

	// Swagger UI
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 需认证的路由组
	secured := api.Group("")
	secured.Use(middleware.AuthMiddleware(deps.JWTManager, deps.BlacklistService))

	// 注册 B 端路由
	RegisterBEndRoutes(api, secured, deps)

	// 注册 C 端路由
	RegisterCEndRoutes(api, secured, deps)

	return nil
}
