//go:build integration

// Package testutil 提供集成测试的基础设施
package testutil

import (
	"testing"
	"time"

	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/router"
	"community-elderly-care-platform/pkg/database"
	pkgredis "community-elderly-care-platform/pkg/redis"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// TestEnv 测试环境，包含所有测试所需的基础设施
type TestEnv struct {
	Engine *gin.Engine
	DB     *database.DB
	Redis  *pkgredis.Client
	Config *config.Config
}

// Setup 初始化完整测试环境：DB + Redis + 种子数据 + 路由
func Setup(t *testing.T) *TestEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db := setupDB(t)
	rdb := setupRedis(t)
	cfg := testConfig(t)

	// 先插入种子数据（NewDeps 中 PermissionService 和
	// GeofenceService 初始化时会查询 DB）
	SeedTestData(t, db.DB)

	engine := gin.New()
	engine.Use(gin.Recovery())
	if err := router.Register(engine, db, rdb, cfg); err != nil {
		t.Fatalf("Register 路由失败: %v", err)
	}

	return &TestEnv{
		Engine: engine,
		DB:     db,
		Redis:  rdb,
		Config: cfg,
	}
}

// setupDB 创建 SQLite 内存数据库并自动迁移
func setupDB(t *testing.T) *database.DB {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		t.Fatalf("打开 SQLite 失败: %v", err)
	}

	err = gdb.AutoMigrate(
		&model.User{},
		&model.UserIdentity{},
		&model.CustomerProfile{},
		&model.ServiceStation{},
		&model.ServiceZone{},
		&model.ServiceRequest{},
		&model.TaskAssignment{},
		&model.TaskHistory{},
		&model.Notification{},
		&model.News{},
		&model.Banner{},
		&model.Menu{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.Report{},
	)
	if err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}

	return &database.DB{DB: gdb}
}

// setupRedis 使用 miniredis 创建测试 Redis
func setupRedis(t *testing.T) *pkgredis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	return &pkgredis.Client{Client: rdb}
}

// testConfig 返回测试用配置
func testConfig(t *testing.T) *config.Config {
	t.Helper()
	storageDir := t.TempDir()
	return &config.Config{
		Server: config.ServerConfig{Port: 8080, Mode: "test", Name: "scare-test"},
		JWT: config.JWTConfig{
			Secret:           "test-jwt-secret-key-for-integration-tests",
			ExpiresIn:        24,
			RefreshExpiresIn: 168,
		},
		Storage: config.StorageConfig{
			Driver: "local",
			Local:  config.LocalStorageConfig{BasePath: storageDir, BaseURL: "http://localhost:8080/static"},
		},
	}
}

// SeedTestData 插入基础测试数据
func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()
	now := time.Now()

	// 角色
	db.Create(&[]model.Role{
		{ID: 1, Code: "admin", Name: "系统管理员", IsSystem: true, Status: "active", Sort: 1},
		{ID: 2, Code: "station_manager", Name: "站点管理员", IsSystem: true, Status: "active", Sort: 2},
		{ID: 3, Code: "staff", Name: "工作人员", IsSystem: true, Status: "active", Sort: 3},
		{ID: 4, Code: "elderly", Name: "老年人", IsSystem: true, Status: "active", Sort: 4},
		{ID: 5, Code: "family", Name: "家属", IsSystem: true, Status: "active", Sort: 5},
	})

	// 用户（密码 Test@123 的 bcrypt hash）
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	db.Create(&[]model.User{
		{ID: 1, Phone: "13800000001", PasswordHash: hash, Name: "系统管理员", Status: "active", StationID: 0, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Phone: "13800000002", PasswordHash: hash, Name: "站点管理员", Status: "active", StationID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Phone: "13800000003", PasswordHash: hash, Name: "工作人员A", Status: "active", StationID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 4, Phone: "13800000004", PasswordHash: hash, Name: "工作人员B", Status: "active", StationID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 10, Phone: "13900000001", PasswordHash: hash, Name: "张大爷", Status: "active", StationID: 0, CreatedAt: now, UpdatedAt: now},
		{ID: 11, Phone: "13900000002", PasswordHash: hash, Name: "李家属", Status: "active", StationID: 0, CreatedAt: now, UpdatedAt: now},
	})

	// 用户身份
	db.Create(&[]model.UserIdentity{
		{ID: 1, UserID: 1, IdentityType: "admin", IsPrimary: true, StationID: 0, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 2, IdentityType: "station_manager", IsPrimary: true, StationID: 1, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
		{ID: 3, UserID: 3, IdentityType: "staff", IsPrimary: true, StationID: 1, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
		{ID: 4, UserID: 4, IdentityType: "staff", IsPrimary: true, StationID: 1, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
		{ID: 10, UserID: 10, IdentityType: "elderly", IsPrimary: true, StationID: 0, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
		{ID: 11, UserID: 11, IdentityType: "family", IsPrimary: true, StationID: 0, Status: "active", GrantedAt: now, CreatedAt: now, UpdatedAt: now},
	})

	// 服务站点
	db.Create(&model.ServiceStation{
		ID: 1, Name: "霍营街道养老服务站", Code: "HY001",
		Address: "北京市昌平区霍营街道", Longitude: 116.4, Latitude: 39.9,
		Capacity: 100, Status: "active", CreatedAt: now, UpdatedAt: now,
	})

	// 服务围栏
	db.Create(&model.ServiceZone{
		ID: 1, StationID: 1, Name: "霍营服务区",
		Points:   `[{"lng":116.3,"lat":39.8},{"lng":116.5,"lat":39.8},{"lng":116.5,"lat":40.0},{"lng":116.3,"lat":40.0}]`,
		Priority: 1, Status: "active", CreatedAt: now, UpdatedAt: now,
	})

	// 权限（公共权限 + 非公共权限）
	db.Create(&[]model.Permission{
		{ID: 1, Code: "system:user:list", Name: "用户列表", Module: "system", Type: "resource", APIPath: "/api/v1/b/users", APIMethod: "GET", Status: "active", Sort: 1},
		{ID: 2, Code: "system:user:get", Name: "用户详情", Module: "system", Type: "resource", APIPath: "/api/v1/b/users/*", APIMethod: "GET", Status: "active", Sort: 2},
		{ID: 3, Code: "service:request:list", Name: "请求列表", Module: "service", Type: "resource", APIPath: "/api/v1/b/requests", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 3},
		{ID: 4, Code: "service:request:get", Name: "请求详情", Module: "service", Type: "resource", APIPath: "/api/v1/b/requests/*", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 4},
		{ID: 5, Code: "service:task:list", Name: "任务列表", Module: "service", Type: "resource", APIPath: "/api/v1/b/tasks", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 5},
		{ID: 6, Code: "service:task:pool", Name: "任务池", Module: "service", Type: "resource", APIPath: "/api/v1/b/tasks/pool", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 6},
		{ID: 7, Code: "service:task:my", Name: "我的任务", Module: "service", Type: "resource", APIPath: "/api/v1/b/tasks/my", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 7},
		{ID: 8, Code: "service:task:get", Name: "任务详情", Module: "service", Type: "resource", APIPath: "/api/v1/b/tasks/*", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 8},
		{ID: 9, Code: "service:station:list", Name: "站点列表", Module: "service", Type: "resource", APIPath: "/api/v1/b/stations", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 9},
		{ID: 10, Code: "service:station:get", Name: "站点详情", Module: "service", Type: "resource", APIPath: "/api/v1/b/stations/*", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 10},
		{ID: 11, Code: "service:zone:list", Name: "围栏列表", Module: "service", Type: "resource", APIPath: "/api/v1/b/zones", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 11},
		{ID: 12, Code: "content:news:list", Name: "新闻列表", Module: "content", Type: "resource", APIPath: "/api/v1/b/news", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 12},
		{ID: 13, Code: "content:news:get", Name: "新闻详情", Module: "content", Type: "resource", APIPath: "/api/v1/b/news/*", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 13},
		{ID: 14, Code: "content:banner:list", Name: "Banner列表", Module: "content", Type: "resource", APIPath: "/api/v1/b/banners", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 14},
		{ID: 15, Code: "system:menu:list", Name: "菜单列表", Module: "system", Type: "resource", APIPath: "/api/v1/b/menus", APIMethod: "GET", Status: "active", Sort: 15},
		{ID: 16, Code: "system:menu:get", Name: "菜单详情", Module: "system", Type: "resource", APIPath: "/api/v1/b/menus/*", APIMethod: "GET", Status: "active", Sort: 16},
		{ID: 17, Code: "system:permission:tree", Name: "权限树", Module: "system", Type: "resource", APIPath: "/api/v1/b/permissions/tree", APIMethod: "GET", Status: "active", Sort: 17},
		{ID: 18, Code: "system:role:permissions", Name: "角色权限", Module: "system", Type: "resource", APIPath: "/api/v1/b/roles/*/permissions", APIMethod: "GET", Status: "active", Sort: 18},
		{ID: 19, Code: "system:notification:list", Name: "通知列表", Module: "system", Type: "resource", APIPath: "/api/v1/b/notifications", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 19},
		{ID: 20, Code: "statistics:dashboard", Name: "仪表盘统计", Module: "statistics", Type: "resource", APIPath: "/api/v1/b/statistics/dashboard", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 20},
		{ID: 21, Code: "statistics:tasks", Name: "任务统计", Module: "statistics", Type: "resource", APIPath: "/api/v1/b/statistics/tasks", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 21},
		{ID: 22, Code: "statistics:requests", Name: "请求统计", Module: "statistics", Type: "resource", APIPath: "/api/v1/b/statistics/requests", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 22},
		{ID: 23, Code: "statistics:today", Name: "今日统计", Module: "statistics", Type: "resource", APIPath: "/api/v1/b/statistics/today", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 23},
	})

	// 角色权限关联
	db.Create(&[]model.RolePermission{
		{ID: 1, RoleID: 2, PermissionID: 1},  // station_manager -> 用户列表
		{ID: 2, RoleID: 2, PermissionID: 9},  // station_manager -> 站点列表
		{ID: 3, RoleID: 2, PermissionID: 10}, // station_manager -> 站点详情
		{ID: 4, RoleID: 3, PermissionID: 6},  // staff -> 任务池
		{ID: 5, RoleID: 3, PermissionID: 7},  // staff -> 我的任务
		{ID: 6, RoleID: 3, PermissionID: 8},  // staff -> 任务详情
	})

	// 新闻
	db.Create(&model.News{ID: 1, Title: "测试新闻", Summary: "测试摘要", Content: "测试内容", Type: "notice", Status: "published", StationID: 1, AuthorID: 1, PublishAt: now, CreatedAt: now, UpdatedAt: now})

	// Banner
	db.Create(&model.Banner{ID: 1, StationID: 0, Title: "测试Banner", ImageURL: "/test.jpg", LinkType: "none", Sort: 1, Status: "active", CreatedAt: now, UpdatedAt: now})

	// 菜单
	db.Create(&[]model.Menu{
		{ID: 1, ParentID: 0, Name: "系统管理", Path: "/system", Icon: "setting", Sort: 1, Hidden: false, Status: "active", CreatedAt: now, UpdatedAt: now},
		{ID: 2, ParentID: 1, Name: "用户管理", Path: "/system/users", PermissionCode: "system:user:list", Sort: 1, Hidden: false, Status: "active", CreatedAt: now, UpdatedAt: now},
	})

	// 客户档案
	db.Create(&model.CustomerProfile{ID: 1, UserID: 10, Gender: "male", CustomerType: "elderly", CreatedAt: now, UpdatedAt: now})
}
