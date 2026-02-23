package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionServiceForTest(t *testing.T) *service.PermissionService {
	t.Helper()

	tmpFile := t.TempDir() + "/permission_middleware_test.db"
	dsn := tmpFile + "?_loc=auto&parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().Local().Truncate(time.Second)
		},
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			name TEXT,
			description TEXT,
			module TEXT,
			type TEXT,
			parent_id INTEGER,
			api_path TEXT,
			api_method TEXT,
			is_public BOOLEAN,
			status TEXT,
			sort INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create permissions table: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			name TEXT,
			description TEXT,
			is_system BOOLEAN,
			status TEXT,
			sort INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create roles table: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE role_permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_id INTEGER,
			permission_id INTEGER,
			created_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create role_permissions table: %v", err)
	}

	role := &model.Role{
		Code:     "staff",
		Name:     "Staff",
		Status:   "active",
		IsSystem: true,
		Sort:     1,
	}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	protected := &model.Permission{
		Code:      "task:list",
		Name:      "Task List",
		Module:    "task",
		Type:      "resource",
		ParentID:  0,
		APIPath:   "/api/v1/b/tasks",
		APIMethod: http.MethodGet,
		IsPublic:  false,
		Status:    "active",
		Sort:      1,
	}
	if err := db.Create(protected).Error; err != nil {
		t.Fatalf("failed to seed protected permission: %v", err)
	}

	public := &model.Permission{
		Code:      "public:ping",
		Name:      "Public Ping",
		Module:    "public",
		Type:      "resource",
		ParentID:  0,
		APIPath:   "/api/v1/public/ping",
		APIMethod: http.MethodGet,
		IsPublic:  true,
		Status:    "active",
		Sort:      2,
	}
	if err := db.Create(public).Error; err != nil {
		t.Fatalf("failed to seed public permission: %v", err)
	}

	rolePermission := &model.RolePermission{
		RoleID:       role.ID,
		PermissionID: protected.ID,
	}
	if err := db.Create(rolePermission).Error; err != nil {
		t.Fatalf("failed to seed role permission: %v", err)
	}

	return service.NewPermissionService(
		db,
		repository.NewPermissionRepository(db),
		repository.NewRoleRepository(db),
		repository.NewRolePermissionRepository(db),
		nil,
	)
}

func runPermissionRequest(t *testing.T, permService *service.PermissionService, method, path string, values map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	r := gin.New()
	r.Handle(method, "/*any",
		func(c *gin.Context) {
			for k, v := range values {
				c.Set(k, v)
			}
			c.Next()
		},
		PermissionMiddleware(permService),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		},
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	permService := setupPermissionServiceForTest(t)

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/tasks", map[string]any{
		"user_type": "c_end",
	}); w.Code != http.StatusOK {
		t.Fatalf("expected c_end to bypass permission check, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/tasks", map[string]any{
		"user_type": "b_end",
	}); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing identities, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/tasks", map[string]any{
		"user_type":       "b_end",
		"user_identities": "staff",
	}); w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for invalid identities format, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/tasks", map[string]any{
		"user_type":       "b_end",
		"user_identities": []string{"admin"},
	}); w.Code != http.StatusOK {
		t.Fatalf("expected admin to bypass permission check, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/public/ping", map[string]any{
		"user_type":       "b_end",
		"user_identities": []string{"staff"},
	}); w.Code != http.StatusOK {
		t.Fatalf("expected public API to bypass permission check, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/forbidden", map[string]any{
		"user_type":       "b_end",
		"user_identities": []string{"staff"},
	}); w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for forbidden path, got %d", w.Code)
	}

	if w := runPermissionRequest(t, permService, http.MethodGet, "/api/v1/b/tasks", map[string]any{
		"user_type":       "b_end",
		"user_identities": []string{"staff"},
	}); w.Code != http.StatusOK {
		t.Fatalf("expected allowed path to pass permission check, got %d", w.Code)
	}
}
