package router

import (
	"community-elderly-care-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterBEndRoutes 注册 B 端所有路由
func RegisterBEndRoutes(api *gin.RouterGroup, secured *gin.RouterGroup, deps *Deps) {
	// =====================================================
	// 公开路由（无需认证）
	// =====================================================
	auth := api.Group("/b/auth")
	{
		auth.POST("/login", deps.BAuthHandler.Login)
		auth.POST("/refresh", deps.BAuthHandler.Refresh)
	}

	// =====================================================
	// 需认证路由
	// =====================================================
	bSecured := secured.Group("/b")
	bSecured.Use(middleware.RequireEndType("b_end"))

	// 基础认证（无需权限检查）
	bSecured.GET("/auth/me", deps.BAuthHandler.Me)
	bSecured.POST("/auth/logout", deps.LogoutHandler.Logout)
	bSecured.GET("/menus/user", deps.MenuHandler.GetUserMenus)

	// =====================================================
	// 需权限检查的路由（使用新的 PermissionMiddleware）
	// =====================================================
	protected := bSecured.Group("")
	protected.Use(middleware.PermissionMiddleware(deps.PermissionService))
	{
		// 服务请求管理
		protected.GET("/requests", deps.RequestHandler.List)
		protected.GET("/requests/:id", deps.RequestHandler.Get)
		protected.PUT("/requests/:id", deps.RequestHandler.Update)
		protected.PUT("/requests/:id/status", deps.RequestHandler.UpdateStatus)
		protected.POST("/requests/:id/cancel", deps.RequestHandler.CancelByAdmin)

		// 任务管理
		protected.GET("/tasks", deps.TaskHandler.List)
		protected.GET("/tasks/pool", deps.TaskHandler.ListPool)
		protected.GET("/tasks/my", deps.TaskHandler.ListMy)
		protected.GET("/tasks/:id", deps.TaskHandler.GetByID)
		protected.POST("/tasks/:id/claim", deps.TaskHandler.Claim)
		protected.POST("/tasks/:id/complete", deps.TaskHandler.Complete)
		protected.POST("/tasks/:id/transfer", deps.TaskHandler.Transfer)

		// 服务站点管理
		protected.GET("/stations", deps.StationHandler.List)
		protected.POST("/stations", deps.StationHandler.Create)
		protected.GET("/stations/:id", deps.StationHandler.Get)
		protected.PUT("/stations/:id", deps.StationHandler.Update)
		protected.DELETE("/stations/:id", deps.StationHandler.Delete)

		// 服务围栏管理
		protected.GET("/zones", deps.ZoneHandler.List)
		protected.POST("/zones", deps.ZoneHandler.Create)
		protected.PUT("/zones/:id", deps.ZoneHandler.Update)
		protected.DELETE("/zones/:id", deps.ZoneHandler.Delete)

		// 用户管理
		protected.GET("/users", deps.UserHandler.List)
		protected.POST("/users", deps.UserHandler.Create)
		protected.GET("/users/:id", deps.UserHandler.GetByID)
		protected.PUT("/users/:id", deps.UserHandler.Update)

		// 角色权限管理
		protected.GET("/permissions/tree", deps.PermissionHandler.GetPermissionTree)
		protected.GET("/roles/:role/permissions", deps.RoleHandler.GetRolePermissions)
		protected.PUT("/roles/:role/permissions", deps.RoleHandler.UpdateRolePermissions)
		protected.PUT("/users/:id/identities", deps.RoleHandler.UpdateUserIdentities)

		// 菜单管理
		protected.GET("/menus", deps.MenuHandler.List)
		protected.POST("/menus", deps.MenuHandler.Create)
		protected.GET("/menus/:id", deps.MenuHandler.GetByID)
		protected.PUT("/menus/:id", deps.MenuHandler.Update)
		protected.DELETE("/menus/:id", deps.MenuHandler.Delete)
		protected.PUT("/menus/sort", deps.MenuHandler.BatchUpdateSort)

		// 通知管理
		protected.GET("/notifications", deps.NotificationHandler.List)
		protected.POST("/notifications/:id/read", deps.NotificationHandler.MarkRead)

		// 文件上传
		protected.POST("/upload", deps.UploadHandler.Upload)

		// 轮播图管理
		protected.GET("/banners", deps.BannerHandler.List)
		protected.POST("/banners", deps.BannerHandler.Create)
		protected.PUT("/banners/:id", deps.BannerHandler.Update)
		protected.DELETE("/banners/:id", deps.BannerHandler.Delete)

		// 新闻管理
		protected.GET("/news", deps.NewsHandler.List)
		protected.POST("/news", deps.NewsHandler.Create)
		protected.GET("/news/:id", deps.NewsHandler.Get)
		protected.PUT("/news/:id", deps.NewsHandler.Update)
		protected.DELETE("/news/:id", deps.NewsHandler.Delete)

		// 老年人档案管理
		protected.GET("/customers", deps.CustomerHandler.List)
		protected.POST("/customers", deps.CustomerHandler.Create)
		protected.GET("/customers/:id", deps.CustomerHandler.Get)
		protected.PUT("/customers/:id", deps.CustomerHandler.Update)
		protected.GET("/customers/:id/service-records", deps.CustomerHandler.GetServiceRecords)

		// 统计数据
		protected.GET("/statistics/dashboard", deps.StatisticsHandler.GetDashboardStats)
		protected.GET("/statistics/tasks", deps.StatisticsHandler.GetTaskStats)
		protected.GET("/statistics/requests", deps.StatisticsHandler.GetRequestStats)
		protected.GET("/statistics/today", deps.StatisticsHandler.GetTodayStats)
		protected.GET("/statistics/overview", deps.StatisticsHandler.GetOverviewStats)
		protected.GET("/statistics/service-types", deps.StatisticsHandler.GetServiceTypeStats)
		protected.GET("/statistics/trend", deps.StatisticsHandler.GetRequestTrend)
		protected.GET("/statistics/efficiency", deps.StatisticsHandler.GetEfficiencyStats)
		protected.GET("/statistics/staff-ranking", deps.StatisticsHandler.GetStaffRanking)

		// 报表
		protected.POST("/reports/generate", deps.ReportHandler.GenerateReport)
		protected.GET("/reports", deps.ReportHandler.ListReports)
		protected.GET("/reports/:id/download", deps.ReportHandler.DownloadReport)
		protected.DELETE("/reports/:id", deps.ReportHandler.DeleteReport)
	}
}
