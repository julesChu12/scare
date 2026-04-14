package router

import (
	"community-elderly-care-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterCEndRoutes 注册 C 端所有路由
func RegisterCEndRoutes(api *gin.RouterGroup, secured *gin.RouterGroup, deps *Deps) {
	// =====================================================
	// 公开路由（无需认证）
	// =====================================================
	auth := api.Group("/c/auth")
	{
		auth.POST("/send-code", deps.CAuthHandler.SendCode)
		auth.POST("/quick-start", deps.CAuthHandler.QuickStart)
		auth.POST("/login", deps.CAuthHandler.Login)
		auth.POST("/refresh", deps.CAuthHandler.Refresh)
		auth.POST("/register", deps.CAuthHandler.Register)
	}

	public := api.Group("/c")
	{
		public.POST("/geocode", deps.CProfileHandler.Geocode)
		public.GET("/geocode/reverse", deps.CProfileHandler.ReverseGeocode)
		public.POST("/geocode/reverse", deps.CProfileHandler.ReverseGeocode)
		public.GET("/stations/:id", deps.StationHandler.GetPublic)
		public.GET("/stations/match", deps.StationHandler.MatchStation)
		public.POST("/stations/match", deps.StationHandler.MatchStation)
		public.POST("/upload", deps.UploadHandler.Upload)
		public.GET("/news", deps.NewsHandler.List)
		public.GET("/news/:id", deps.NewsHandler.Get)
		public.GET("/banners", deps.BannerHandler.ListForC)
	}

	// =====================================================
	// 需认证路由
	// =====================================================
	cSecured := secured.Group("/c")
	cSecured.Use(middleware.RequireEndType("c_end"))
	{
		// 认证相关
		cSecured.GET("/auth/me", deps.CAuthHandler.Me)
		cSecured.GET("/auth/check", deps.CAuthHandler.CheckToken)
		cSecured.POST("/auth/logout", deps.LogoutHandler.Logout)

		// 个人资料
		cSecured.PUT("/profile", deps.CProfileHandler.UpdateProfile)

		// 服务请求
		cSecured.POST("/requests", deps.RequestHandler.Create)
		cSecured.GET("/requests", deps.RequestHandler.List)
		cSecured.GET("/requests/:id", deps.RequestHandler.Get)
		cSecured.POST("/requests/:id/cancel", deps.RequestHandler.Cancel)
		cSecured.POST("/requests/:id/rate", deps.RequestHandler.Rate)

		// 通知
		cSecured.GET("/notifications", deps.NotificationHandler.List)
		cSecured.POST("/notifications/:id/read", deps.NotificationHandler.MarkRead)
	}
}
