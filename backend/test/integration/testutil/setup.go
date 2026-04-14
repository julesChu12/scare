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

	"github.com/gin-gonic/gin"
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
	if sqlDB, err := db.DB.DB(); err == nil {
		t.Cleanup(func() {
			_ = sqlDB.Close()
		})
	}
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

// setupDB 创建独立的 SQLite 测试数据库并初始化表结构
func setupDB(t *testing.T) *database.DB {
	t.Helper()

	tmpFile := t.TempDir() + "/integration.db"
	dsn := tmpFile + "?_loc=auto&parseTime=true"
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		t.Fatalf("打开 SQLite 失败: %v", err)
	}

	createTables(t, gdb)

	return &database.DB{DB: gdb}
}

func createTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone TEXT UNIQUE,
			password_hash TEXT,
			name TEXT,
			email TEXT,
			avatar TEXT,
			gender TEXT,
			birth_date DATE,
			id_card TEXT,
			station_id INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			id_card_hmac TEXT,
			id_card_masked TEXT
		);
		`,
		`
		CREATE TABLE user_identities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			identity_type TEXT NOT NULL,
			is_primary INTEGER NOT NULL DEFAULT 0,
			station_id INTEGER DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			granted_at DATETIME,
			granted_by INTEGER,
			revoked_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE customer_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			id_card TEXT,
			address TEXT,
			latitude REAL,
			longitude REAL,
			customer_type TEXT,
			emergency_contact TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			gender TEXT,
			birth_date DATE,
			health_status TEXT,
			disability_level TEXT,
			medical_history TEXT,
			special_needs TEXT
		);
		`,
		`
		CREATE TABLE service_stations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT,
			address TEXT,
			phone TEXT,
			latitude REAL,
			longitude REAL,
			service_area TEXT,
			capacity INTEGER DEFAULT 10,
			work_hours TEXT,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE service_zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			points TEXT NOT NULL,
			priority INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE service_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_no TEXT,
			user_id INTEGER NOT NULL,
			service_type TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			description TEXT,
			submit_location_lat REAL,
			submit_location_lng REAL,
			service_location_lat REAL,
			service_location_lng REAL,
			contact_name TEXT,
			contact_phone TEXT,
			address TEXT,
			appointment_time DATETIME,
			urgency TEXT,
			source_station_id INTEGER DEFAULT 0,
			station_id INTEGER DEFAULT 0,
			dispatch_basis TEXT,
			needs_manual_verify INTEGER DEFAULT 0,
			reject_reason TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			rating INTEGER DEFAULT 0,
			feedback TEXT
		);
		`,
		`
		CREATE TABLE task_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id INTEGER NOT NULL,
			station_id INTEGER NOT NULL,
			staff_id INTEGER DEFAULT 0,
			status TEXT DEFAULT 'dispatched',
			claimed_at DATETIME,
			completed_at DATETIME,
			rating INTEGER DEFAULT 0,
			feedback TEXT,
			staff_notes TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE task_histories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			request_id INTEGER NOT NULL,
			action TEXT NOT NULL,
			operator_id INTEGER NOT NULL,
			from_staff_id INTEGER DEFAULT 0,
			to_staff_id INTEGER DEFAULT 0,
			from_station_id INTEGER DEFAULT 0,
			to_station_id INTEGER DEFAULT 0,
			status_before TEXT,
			status_after TEXT,
			remark TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT,
			type TEXT,
			related_id INTEGER DEFAULT 0,
			related_type TEXT,
			channel TEXT NOT NULL,
			send_status TEXT DEFAULT 'pending',
			sent_at DATETIME,
			is_read INTEGER DEFAULT 0,
			read_at DATETIME,
			retry_count INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER DEFAULT 0,
			title TEXT NOT NULL,
			summary TEXT,
			content TEXT,
			cover_url TEXT,
			type TEXT DEFAULT 'news',
			status TEXT DEFAULT 'draft',
			author_id INTEGER DEFAULT 0,
			publish_at DATETIME,
			view_count INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE banners (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER NOT NULL DEFAULT 0,
			title TEXT,
			image_url TEXT NOT NULL,
			link_type TEXT NOT NULL DEFAULT 'none',
			link_value TEXT,
			sort INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE menus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			parent_id INTEGER NOT NULL DEFAULT 0,
			name TEXT NOT NULL,
			path TEXT,
			component TEXT,
			icon TEXT,
			permission_code TEXT,
			sort INTEGER NOT NULL DEFAULT 0,
			hidden INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			is_system INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			sort INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			module TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'resource',
			parent_id INTEGER NOT NULL DEFAULT 0,
			api_path TEXT,
			api_method TEXT,
			is_public INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			sort INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE role_permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_id INTEGER NOT NULL,
			permission_id INTEGER NOT NULL,
			created_at DATETIME
		);
		`,
		`
		CREATE TABLE reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			format TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL DEFAULT 0,
			station_id INTEGER DEFAULT 0,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			created_by INTEGER NOT NULL
		);
		`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("创建测试表失败: %v", err)
		}
	}
}

// setupRedis 返回 nil，服务层会退化到进程内存实现
func setupRedis(t *testing.T) *pkgredis.Client {
	t.Helper()
	return nil
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

	// 权限：按当前 RBAC 口径保留集成测试会用到的最小集合
	db.Create(&[]model.Permission{
		{ID: 3, Code: "public:auth:me", Name: "获取当前用户", Module: "public", Type: "resource", APIPath: "/api/v1/*/auth/me", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 1},
		{ID: 4, Code: "public:auth:logout", Name: "登出", Module: "public", Type: "resource", APIPath: "/api/v1/*/auth/logout", APIMethod: "POST", IsPublic: true, Status: "active", Sort: 2},
		{ID: 6, Code: "public:upload:file", Name: "上传文件", Module: "public", Type: "resource", APIPath: "/api/v1/*/upload", APIMethod: "POST", IsPublic: true, Status: "active", Sort: 3},
		{ID: 8, Code: "public:notification:list", Name: "查看通知", Module: "public", Type: "resource", APIPath: "/api/v1/*/notifications", APIMethod: "GET", IsPublic: true, Status: "active", Sort: 4},
		{ID: 9, Code: "public:notification:read", Name: "标记通知已读", Module: "public", Type: "resource", APIPath: "/api/v1/*/notifications/*/read", APIMethod: "POST", IsPublic: true, Status: "active", Sort: 5},
		{ID: 13, Code: "service:request:list", Name: "请求列表", Module: "service", Type: "button", APIPath: "/api/v1/b/requests", APIMethod: "GET", Status: "active", Sort: 13},
		{ID: 14, Code: "service:request:detail", Name: "请求详情", Module: "service", Type: "resource", APIPath: "/api/v1/b/requests/*", APIMethod: "GET", Status: "active", Sort: 14},
		{ID: 16, Code: "service:task:pool", Name: "任务池", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks/pool", APIMethod: "GET", Status: "active", Sort: 16},
		{ID: 17, Code: "service:task:my", Name: "我的任务", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks/my", APIMethod: "GET", Status: "active", Sort: 17},
		{ID: 18, Code: "service:task:claim", Name: "认领任务", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks/*/claim", APIMethod: "POST", Status: "active", Sort: 18},
		{ID: 19, Code: "service:task:complete", Name: "完成任务", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks/*/complete", APIMethod: "POST", Status: "active", Sort: 19},
		{ID: 22, Code: "station:list:view", Name: "查看站点", Module: "station", Type: "button", APIPath: "/api/v1/b/stations", APIMethod: "GET", Status: "active", Sort: 22},
		{ID: 23, Code: "station:list:detail", Name: "站点详情", Module: "station", Type: "resource", APIPath: "/api/v1/b/stations/*", APIMethod: "GET", Status: "active", Sort: 23},
		{ID: 24, Code: "station:list:create", Name: "创建站点", Module: "station", Type: "button", APIPath: "/api/v1/b/stations", APIMethod: "POST", Status: "active", Sort: 24},
		{ID: 25, Code: "station:list:update", Name: "编辑站点", Module: "station", Type: "button", APIPath: "/api/v1/b/stations/*", APIMethod: "PUT", Status: "active", Sort: 25},
		{ID: 26, Code: "station:list:delete", Name: "删除站点", Module: "station", Type: "button", APIPath: "/api/v1/b/stations/*", APIMethod: "DELETE", Status: "active", Sort: 26},
		{ID: 28, Code: "station:zone:list", Name: "围栏列表", Module: "station", Type: "button", APIPath: "/api/v1/b/zones", APIMethod: "GET", Status: "active", Sort: 28},
		{ID: 29, Code: "station:zone:create", Name: "创建围栏", Module: "station", Type: "button", APIPath: "/api/v1/b/zones", APIMethod: "POST", Status: "active", Sort: 29},
		{ID: 30, Code: "station:zone:update", Name: "编辑围栏", Module: "station", Type: "button", APIPath: "/api/v1/b/zones/*", APIMethod: "PUT", Status: "active", Sort: 30},
		{ID: 31, Code: "station:zone:delete", Name: "删除围栏", Module: "station", Type: "button", APIPath: "/api/v1/b/zones/*", APIMethod: "DELETE", Status: "active", Sort: 31},
		{ID: 34, Code: "system:user:list", Name: "用户列表", Module: "system", Type: "button", APIPath: "/api/v1/b/users", APIMethod: "GET", Status: "active", Sort: 34},
		{ID: 35, Code: "system:user:detail", Name: "用户详情", Module: "system", Type: "resource", APIPath: "/api/v1/b/users/*", APIMethod: "GET", Status: "active", Sort: 35},
		{ID: 36, Code: "system:user:create", Name: "创建用户", Module: "system", Type: "button", APIPath: "/api/v1/b/users", APIMethod: "POST", Status: "active", Sort: 36},
		{ID: 37, Code: "system:user:update", Name: "编辑用户", Module: "system", Type: "button", APIPath: "/api/v1/b/users/*", APIMethod: "PUT", Status: "active", Sort: 37},
		{ID: 41, Code: "system:role:permissions", Name: "角色权限", Module: "system", Type: "button", APIPath: "/api/v1/b/roles/*/permissions", APIMethod: "GET", Status: "active", Sort: 41},
		{ID: 42, Code: "system:role:update", Name: "更新权限", Module: "system", Type: "button", APIPath: "/api/v1/b/roles/*/permissions", APIMethod: "PUT", Status: "active", Sort: 42},
		{ID: 43, Code: "system:permission:tree", Name: "权限树", Module: "system", Type: "button", APIPath: "/api/v1/b/permissions/tree", APIMethod: "GET", Status: "active", Sort: 43},
		{ID: 45, Code: "system:menu:list", Name: "菜单列表", Module: "system", Type: "button", APIPath: "/api/v1/b/menus", APIMethod: "GET", Status: "active", Sort: 45},
		{ID: 46, Code: "system:menu:detail", Name: "菜单详情", Module: "system", Type: "resource", APIPath: "/api/v1/b/menus/*", APIMethod: "GET", Status: "active", Sort: 46},
		{ID: 47, Code: "system:menu:create", Name: "创建菜单", Module: "system", Type: "button", APIPath: "/api/v1/b/menus", APIMethod: "POST", Status: "active", Sort: 47},
		{ID: 48, Code: "system:menu:update", Name: "编辑菜单", Module: "system", Type: "button", APIPath: "/api/v1/b/menus/*", APIMethod: "PUT", Status: "active", Sort: 48},
		{ID: 49, Code: "system:menu:delete", Name: "删除菜单", Module: "system", Type: "button", APIPath: "/api/v1/b/menus/*", APIMethod: "DELETE", Status: "active", Sort: 49},
		{ID: 50, Code: "system:menu:sort", Name: "菜单排序", Module: "system", Type: "button", APIPath: "/api/v1/b/menus/sort", APIMethod: "PUT", Status: "active", Sort: 50},
		{ID: 53, Code: "content:banner:list", Name: "轮播图列表", Module: "content", Type: "button", APIPath: "/api/v1/b/banners", APIMethod: "GET", Status: "active", Sort: 53},
		{ID: 54, Code: "content:banner:create", Name: "创建轮播图", Module: "content", Type: "button", APIPath: "/api/v1/b/banners", APIMethod: "POST", Status: "active", Sort: 54},
		{ID: 55, Code: "content:banner:update", Name: "编辑轮播图", Module: "content", Type: "button", APIPath: "/api/v1/b/banners/*", APIMethod: "PUT", Status: "active", Sort: 55},
		{ID: 56, Code: "content:banner:delete", Name: "删除轮播图", Module: "content", Type: "button", APIPath: "/api/v1/b/banners/*", APIMethod: "DELETE", Status: "active", Sort: 56},
		{ID: 58, Code: "content:news:list", Name: "新闻列表", Module: "content", Type: "button", APIPath: "/api/v1/b/news", APIMethod: "GET", Status: "active", Sort: 58},
		{ID: 59, Code: "content:news:detail", Name: "新闻详情", Module: "content", Type: "resource", APIPath: "/api/v1/b/news/*", APIMethod: "GET", Status: "active", Sort: 59},
		{ID: 60, Code: "content:news:create", Name: "创建新闻", Module: "content", Type: "button", APIPath: "/api/v1/b/news", APIMethod: "POST", Status: "active", Sort: 60},
		{ID: 61, Code: "content:news:update", Name: "编辑新闻", Module: "content", Type: "button", APIPath: "/api/v1/b/news/*", APIMethod: "PUT", Status: "active", Sort: 61},
		{ID: 62, Code: "content:news:delete", Name: "删除新闻", Module: "content", Type: "button", APIPath: "/api/v1/b/news/*", APIMethod: "DELETE", Status: "active", Sort: 62},
		{ID: 65, Code: "data:statistics:view", Name: "查看统计", Module: "data", Type: "button", APIPath: "/api/v1/b/statistics/*", APIMethod: "GET", Status: "active", Sort: 65},
		{ID: 67, Code: "data:report:list", Name: "报表列表", Module: "data", Type: "button", APIPath: "/api/v1/b/reports", APIMethod: "GET", Status: "active", Sort: 67},
		{ID: 68, Code: "data:report:generate", Name: "生成报表", Module: "data", Type: "button", APIPath: "/api/v1/b/reports/generate", APIMethod: "POST", Status: "active", Sort: 68},
		{ID: 69, Code: "data:report:download", Name: "下载报表", Module: "data", Type: "button", APIPath: "/api/v1/b/reports/*/download", APIMethod: "GET", Status: "active", Sort: 69},
		{ID: 70, Code: "data:report:delete", Name: "删除报表", Module: "data", Type: "button", APIPath: "/api/v1/b/reports/*", APIMethod: "DELETE", Status: "active", Sort: 70},
		{ID: 71, Code: "service:task:transfer", Name: "转派任务", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks/*/transfer", APIMethod: "POST", Status: "active", Sort: 71},
		{ID: 72, Code: "service:task:list", Name: "任务列表", Module: "service", Type: "button", APIPath: "/api/v1/b/tasks", APIMethod: "GET", Status: "active", Sort: 72},
		{ID: 73, Code: "service:task:detail", Name: "任务详情", Module: "service", Type: "resource", APIPath: "/api/v1/b/tasks/*", APIMethod: "GET", Status: "active", Sort: 73},
		{ID: 74, Code: "service:request:update", Name: "编辑请求", Module: "service", Type: "button", APIPath: "/api/v1/b/requests/*", APIMethod: "PUT", Status: "active", Sort: 74},
		{ID: 75, Code: "service:request:cancel", Name: "取消请求", Module: "service", Type: "button", APIPath: "/api/v1/b/requests/*/cancel", APIMethod: "POST", Status: "active", Sort: 75},
	})

	// 角色权限关联：与当前种子保持一致的最小测试集合
	db.Create(&[]model.RolePermission{
		{ID: 1, RoleID: 2, PermissionID: 13},
		{ID: 2, RoleID: 2, PermissionID: 14},
		{ID: 3, RoleID: 2, PermissionID: 16},
		{ID: 4, RoleID: 2, PermissionID: 17},
		{ID: 5, RoleID: 2, PermissionID: 18},
		{ID: 6, RoleID: 2, PermissionID: 19},
		{ID: 7, RoleID: 2, PermissionID: 22},
		{ID: 8, RoleID: 2, PermissionID: 23},
		{ID: 9, RoleID: 2, PermissionID: 28},
		{ID: 10, RoleID: 2, PermissionID: 29},
		{ID: 11, RoleID: 2, PermissionID: 30},
		{ID: 12, RoleID: 2, PermissionID: 31},
		{ID: 13, RoleID: 2, PermissionID: 34},
		{ID: 14, RoleID: 2, PermissionID: 35},
		{ID: 15, RoleID: 2, PermissionID: 36},
		{ID: 16, RoleID: 2, PermissionID: 37},
		{ID: 17, RoleID: 2, PermissionID: 65},
		{ID: 18, RoleID: 2, PermissionID: 67},
		{ID: 19, RoleID: 2, PermissionID: 68},
		{ID: 20, RoleID: 2, PermissionID: 69},
		{ID: 21, RoleID: 2, PermissionID: 71},
		{ID: 22, RoleID: 2, PermissionID: 72},
		{ID: 23, RoleID: 2, PermissionID: 73},
		{ID: 24, RoleID: 2, PermissionID: 74},
		{ID: 25, RoleID: 2, PermissionID: 75},
		{ID: 26, RoleID: 3, PermissionID: 13},
		{ID: 27, RoleID: 3, PermissionID: 14},
		{ID: 28, RoleID: 3, PermissionID: 17},
		{ID: 29, RoleID: 3, PermissionID: 18},
		{ID: 30, RoleID: 3, PermissionID: 19},
		{ID: 31, RoleID: 3, PermissionID: 65},
		{ID: 32, RoleID: 3, PermissionID: 73},
	})

	// 新闻
	db.Create(&model.News{ID: 1, Title: "测试新闻", Summary: "测试摘要", Content: "测试内容", Type: "notice", Status: "published", StationID: 0, AuthorID: 1, PublishAt: now, CreatedAt: now, UpdatedAt: now})

	// Banner
	db.Create(&model.Banner{ID: 1, StationID: 0, Title: "测试Banner", ImageURL: "/test.jpg", LinkType: "none", Sort: 1, Status: "active", CreatedAt: now, UpdatedAt: now})

	// 菜单
	db.Create(&[]model.Menu{
		{ID: 1, ParentID: 0, Name: "系统管理", Path: "/system", Icon: "setting", Sort: 1, Hidden: false, Status: "active", CreatedAt: now, UpdatedAt: now},
		{ID: 2, ParentID: 1, Name: "用户管理", Path: "/system/users", PermissionCode: "system:user:list", Sort: 1, Hidden: false, Status: "active", CreatedAt: now, UpdatedAt: now},
	})

	// 客户档案
	db.Create(&[]model.CustomerProfile{
		{ID: 1, UserID: 10, Gender: "male", CustomerType: "elderly", Address: "北京市昌平区霍营街道", Latitude: 39.9, Longitude: 116.4, CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 11, Gender: "female", CustomerType: "family", Address: "北京市昌平区霍营街道", Latitude: 39.9, Longitude: 116.4, CreatedAt: now, UpdatedAt: now},
	})
}
